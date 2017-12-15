package proc

import (
    "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering SamplerProc")
    core.GetRegistryInstance()["SamplerProc"] = NewSamplerProc
}

type SamplerProc struct {
    *core.ComponentBase
    Every uint64
}


func NewSamplerProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
    log.Info("Creating SamplerProc")
    return &SamplerProc{core.NewComponentBase(inQ, outQ, cfg), uint64(cfg["every"].(float64))}
}


func (p *SamplerProc) Run() {
    log.Debug("SamplerProc Starting ... ")
    p.MustStop = false
    for !p.MustStop {
        log.Debug("SamplerProc Reading")
        e := <- p.InQ
        p.StatsAddMesg()

        if (p.Stats.MsgCount % p.Every) != 0 {
            log.Debug("SamplerProc Dropping")
            e = nil
            continue
        }

        log.Debug("SamplerProc Forwarding")
        p.OutQ<-e

        // Stats
        p.PrintStats("Log", 50000)

    }

    log.Info("LogProc Stopping!?")
}
