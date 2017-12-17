package proc

import (
    "os"
    "io"
    "fmt"
    "time"
    "sync"
    "bufio"
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering InListProc")
    GetRegistryInstance()["InListProc"] = NewInListProc
}

type InListProc struct {
    *ComponentBase
    List map[string]bool
    ListLock *sync.Mutex
    FilePath string
    InField string
    OutField string
    ReloadMinutes int
}

func NewInListProc(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating InListProc")

    // Set this modules log level

    list := map[string]bool{}
    tmp, ok1 := cfg["list"].([]interface{})

    if ok1 {
        // Convert list to map for faster lookup
        for _, item := range tmp {
            list[item.(string)] = true
        }
    } else {
        list = nil
    }

    fpath, ok2 := cfg["filepath"].(string)

    if !ok1 && !ok2 {
        panic("InList requires a list either in 'list' or in 'filepath'")
    }

    m := &InListProc{NewComponentBase(inQ, outQ, cfg),
        list, &sync.Mutex{}, fpath,
        cfg["in_field"].(string),
        cfg["out_field"].(string),
        int(cfg["reload_minutes"].(float64))}

    m.Tag = "INLIST-LOG"
    return m
}

func (p *InListProc) Run() {
    log.Debug("InListProc Starting ... ")

    // Spawn the loader
    if p.List == nil {
        go func(p *InListProc) {
            p.loadList()
            time.Sleep(time.Duration(p.ReloadMinutes)*time.Minute)
        }(p)
    }

    p.MustStop = false
    cfg_error := false

    for !p.MustStop {
        log.Debug("InListProc Reading")
        e := <- p.InQ

        what, ok := e.Data[p.InField]
        if !ok {
            // This is a user error, maybe error once?
            if !cfg_error {
                log.Error("Cannot find field ", p.InField)
                cfg_error = true
            }
        }

        whatstr := fmt.Sprintf("%v", what)
        p.ListLock.Lock()
        if _, ok := p.List[whatstr]; ok {
            e.Data[p.OutField] = true
        } else {
            e.Data[p.OutField] = false
        }
        p.ListLock.Unlock()



        p.OutQ<-e

        // Stats
        p.StatsAddMesg()
        p.PrintStats()

    }

    log.Info("InListProc Stopping!?")
}

func (p *InListProc) loadList() {
    f, err := os.Open(p.FilePath)
    if err != nil {
        log.Error("INLIST: Could not load list file")
        return
    }

    p.ListLock.Lock()

    log.Warn("INLIST: Reading file")
    reader := bufio.NewReader(f)

    count := 1

    p.List = map[string]bool{}

    line, _, err := reader.ReadLine()
    for err != io.EOF {
        p.List[string(line)] = true
        count += 1
        line, _, err = reader.ReadLine()
    }

    log.Info("INLIST: Done! Loaded ", count, " items!")

    f.Close()
    p.ListLock.Unlock()
}
