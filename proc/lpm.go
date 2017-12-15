package proc

import (
    "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering LPMProg")
    core.GetRegistryInstance()["LPMProg"] = NewLPMProg
}

type LPMProg struct {
    *core.ComponentBase
    FilePath string
}

func NewLPMProg(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
    log.Info("Creating LPMProg")

    fpath, ok := cfg["filepath"].(string)
    if !ok {
        panic("LPM 'filepath' missing")
    }

    return &LPMProg{core.NewComponentBase(inQ, outQ, cfg), fpath}
}


func (p *LPMProg) Run() {
    log.Debug("LPMProg Starting ... ")
    p.MustStop = false
    for !p.MustStop {
        log.Debug("LPMProg Reading")
        e := <- p.InQ

        // TODO: process

        p.OutQ<-e


        // Stats
        p.StatsAddMesg()
        p.PrintStats("LPM", 50000)

    }

    log.Info("LPMProg Stopping")
}
