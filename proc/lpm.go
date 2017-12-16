package proc

import (
    "os"
    "io"
    "sync"
    "time"
    "bufio"
    "bytes"
    "strings"
    . "gopipe/core"

    "encoding/json"
    "github.com/asergeyev/nradix"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering LPMProc")
    GetRegistryInstance()["LPMProc"] = NewLPMProc
}

type LPMOutField struct {
    NewKey string
    MetaKey string
}

type LPMProc struct {
    *ComponentBase
    Tree *nradix.Tree
    TreeLock *sync.Mutex
    FilePath string
    ReloadMinutes int
    InFields []string
    OutFields []LPMOutField
}

func NewLPMProc(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating LPMProc")

    fpath, ok := cfg["filepath"].(string)
    if !ok {
        panic("LPM 'filepath' missing")
    }

    in_fields := []string{}
    if tmp, ok := cfg["in_fields"].([]interface{}); ok {
        in_fields = InterfaceToStringArray(tmp)
    }


    out_fields := []LPMOutField{}
    tmpof, ok := cfg["out_fields"].([]interface{})
    for _, v := range tmpof {
        v2 := v.(map[string]interface{})
        out_fields = append(
            out_fields,
            LPMOutField{v2["newkey"].(string),  v2["metakey"].(string)})
    }

    m := &LPMProc{NewComponentBase(inQ, outQ, cfg),
         nradix.NewTree(100), &sync.Mutex{}, fpath,
         int(cfg["reload_minutes"].(float64)),
         in_fields, out_fields}

     m.Tag = "PROC-LPM"
     return m
}


func (p *LPMProc) Run() {
    log.Debug("LPMProc Starting ... ")

    // Spawn the loader
    go func(p *LPMProc) {
        p.loadTree()
        time.Sleep(time.Duration(p.ReloadMinutes)*time.Minute)
    }(p)

    p.MustStop = false
    cfg_error := false

    for !p.MustStop {
        // Do not read until we lock the tree!
        log.Debug("LPMProc Reading ", p.MustStop)
        e := <- p.InQ

        p.TreeLock.Lock()

        for _, ifield := range p.InFields {

            what, ok := e.Data[ifield].(string)
            if !ok {
                // This is a user error, maybe error once?
                if !cfg_error {
                    log.Error("Cannot find field ", ifield)
                    cfg_error = true
                }
                continue
            }
            // Get the node
            meta, err := p.Tree.FindCIDR(what)
            if err != nil {
                log.Error("LPM error in find: ", err.Error())
                continue
            }
            if meta == nil {
                log.Debug("Could not find prefix for '", ifield, "' -> ", e.Data[ifield].(string))
                continue
            }


            // Generate new fields
            for _, ofield := range p.OutFields {
                new_field := strings.Replace(ofield.NewKey, "{{in_field}}", ifield, 1)
                e.Data[new_field] = meta.(map[string]interface{})[ofield.MetaKey]
            }
        }

        // Now unlock and push
        p.TreeLock.Unlock()

        p.OutQ<-e

        // Stats
        p.StatsAddMesg()
        p.PrintStats()

    }

    log.Info("LPMProc Stopping")
}

func (p *LPMProc) loadTree() {
    p.TreeLock.Lock()

    f, err := os.Open(p.FilePath)
    if err != nil {
        log.Error("LPM: Could not load prefix file")
        p.TreeLock.Unlock()
        return
    }

    p.Tree = nradix.NewTree(100)

    log.Warn("LPM: Reading file")
    reader := bufio.NewReader(f)

    count := 1

    line, _, err := reader.ReadLine()
    for err != io.EOF {
        json_data := map[string]interface{}{}
        parts := bytes.Split(line, []byte(" "))
        meta := bytes.Join(parts[1:], []byte(""))
        if json.Unmarshal(meta, &json_data) != nil {
            log.Error("LPM: Unable to parse prefix meta-data: ", string(meta))
        }

        json_data["prefix"] = string(parts[0])
        p.Tree.AddCIDRb(parts[0], json_data)
        count += 1
        line, _, err = reader.ReadLine()
    }

    log.Info("LPM: Done! Loaded ", count, " prefixes!")
    f.Close()
    p.TreeLock.Unlock()

}
