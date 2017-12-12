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
    config core.Config
    inQ chan core.Event
    outQ chan core.Event
    mustStop bool
}

func NewLogProc(inQ chan core.Event, outQ chan core.Event, cfg core.Config) core.Component {
    log.Info("Creating LogProc")
    return &LogProc{cfg, inQ, outQ, false}
}

func (p *LogProc) Stop() {
    p.mustStop = true
}

func (p *LogProc) Run() {
    log.Debug("LogProc Starting ... ")
    p.mustStop = false
    for !p.mustStop {
        log.Debug("LogProc Reading")
        e := <- p.inQ
        log.Info("Log Proc " + e.ToString())

        if p.outQ != nil {
            log.Debug("LogProc Pushing")
            p.outQ<-e
        }
    }

    log.Info("LogProc Stopping!?")
}
