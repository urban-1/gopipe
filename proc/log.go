package input

import (
    "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering LogProc")
    core.GetRegistryInstance()["LogProc"] = NewLogProc
}

type LogProc struct {
    core.ComponentBase
}

func NewLogProc(inQ chan core.Event, outQ chan core.Event, cfg core.Config) core.Component {
    log.Info("Creating LogProc")
    return &LogProc{*core.NewComponentBase(inQ, outQ, cfg)}
}

func (p *LogProc) Stop() {
    p.MustStop = true
}

func (p *LogProc) Run() {
    log.Debug("LogProc Starting ... ")
    p.MustStop = false
    for !p.MustStop {
        log.Debug("LogProc Reading")
        e := <- p.InQ
        log.Info("Log Proc " + e.ToString())

        if p.OutQ != nil {
            log.Debug("LogProc Pushing")
            p.OutQ<-e
        }
    }

    log.Info("LogProc Stopping!?")
}
