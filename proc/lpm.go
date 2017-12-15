package proc

import (
    "os"
    "io"
    "sync"
    "time"
    "bufio"
    "bytes"
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
    Key string
    Value string
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
    if tmp, ok := cfg["in_fields"].(InterfaceArray); ok {
        in_fields = tmp.ToStringArray()
    }

    out_fields := []LPMOutField{}
    tmpof, ok := cfg["out_fields"].([]map[string]interface{})
    for _, v := range tmpof {
        out_fields = append(
            out_fields,
            LPMOutField{v["key"].(string),  v["value"].(string)})
    }

    return &LPMProc{NewComponentBase(inQ, outQ, cfg),
         nil, &sync.Mutex{}, fpath,
         int(cfg["reload_minutes"].(float64)),
         in_fields, out_fields}
}


func (p *LPMProc) Run() {
    log.Debug("LPMProc Starting ... ")

    // Spawn the loader
    go func(p *LPMProc) {
        p.loadTree()
        time.Sleep(time.Duration(p.ReloadMinutes)*time.Minute)
    }(p)

    p.MustStop = false
    for !p.MustStop {
        // Do not read until we lock the tree!
        log.Debug("LPMProc Reading")
        e := <- p.InQ

        p.TreeLock.Lock()
        // TODO: process

        // Now unlock and push
        p.TreeLock.Unlock()

        p.OutQ<-e


        // Stats
        p.StatsAddMesg()
        p.PrintStats("LPM", 50000)

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

    log.Error("LPM: Reading file")
    reader := bufio.NewReader(f)
    json_data := map[string]interface{}{}
    count := 1

    line, _, err := reader.ReadLine()
    for err != io.EOF {

        parts := bytes.Split(line, []byte(" "))
        meta := bytes.Join(parts[1:], []byte(""))
        if json.Unmarshal(meta, &json_data) != nil {
            log.Error("LPM: Unable to parse prefix meta-data: ", string(meta))
        }

        p.Tree.AddCIDRb(parts[0], json_data)
        count += 1
        if (count % 100000) == 0 {
            //log.Info(".")
            log.Info(string(parts[0]), string(meta))
        }

        line, _, err = reader.ReadLine()
    }

    log.Info("LPM: Done! Loaded ", count, " prefixes!")
    f.Close()
    p.TreeLock.Unlock()

}
