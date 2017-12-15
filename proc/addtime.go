package proc

import (
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering AddTimeProc")
    GetRegistryInstance()["AddTimeProc"] = NewAddTimeProc
}

type AddTimeProc struct {
    *ComponentBase
    FieldName string
}

func NewAddTimeProc(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating AddTimeProc")
    field_name, ok := cfg["field_name"].(string)
    if !ok {
        field_name = "timestamp"
    }
    return &AddTimeProc{NewComponentBase(inQ, outQ, cfg), field_name}
}


func (p *AddTimeProc) Run() {
    log.Debug("AddTimeProc Starting ... ")
    p.MustStop = false
    for !p.MustStop {
        log.Debug("AddTimeProc Reading")
        e := <- p.InQ

        e.Data[p.FieldName] = uint64(e.Timestamp.UnixNano()/1000000)
        p.OutQ<-e

        // Stats
        p.StatsAddMesg()
        p.PrintStats("ADD-TS", 50000)

    }

    log.Info("AddTimeProc Stopping!?")
}
