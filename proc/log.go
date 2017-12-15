package proc

import (
    "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering LogProc")
    core.GetRegistryInstance()["LogProc"] = NewLogProc
}

type LogProc struct {
    *core.ComponentBase
    logFunc func(args ...interface{})
}

func NewLogProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
    log.Info("Creating LogProc")

    // Set this modules log level
    logFunc := log.Debug
    if level, ok := cfg["level"].(string); ok {
        switch level {
        case "debug":
            logFunc = log.Debug
        case "info":
            logFunc = log.Info
        case "warn":
            logFunc = log.Warn
        }
    }
    return &LogProc{core.NewComponentBase(inQ, outQ, cfg), logFunc}
}


func (p *LogProc) Run() {
    log.Debug("LogProc Starting ... ")
    p.MustStop = false
    for !p.MustStop {
        log.Debug("LogProc Reading")
        e := <- p.InQ
        p.logFunc("LogProc: " + e.ToString())

        if p.OutQ != nil {
            log.Debug("LogProc Pushing")
            p.OutQ<-e
        }

        // Stats
        p.StatsAddMesg()
        p.PrintStats("Log", 50000)

    }

    log.Info("LogProc Stopping!?")
}
