/*
    - LOG: Use the logger (logrus) to print the event. This component does not
    modify the events in any way. NOTE: This can also be used as an OUTPUT component!
 */
package proc

import (
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering LogProc")
    GetRegistryInstance()["LogProc"] = NewLogProc
}

type LogProc struct {
    *ComponentBase
    logFunc func(args ...interface{})
}

func NewLogProc(inQ chan *Event, outQ chan *Event, cfg Config) Component {
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
    m := &LogProc{NewComponentBase(inQ, outQ, cfg), logFunc}
    m.Tag = "PROC-LOG"
    return m
}


func (p *LogProc) Run() {
    log.Debug("LogProc Starting ... ")
    p.MustStop = false
    for !p.MustStop {
        log.Debug("LogProc Reading")
        e, err := p.ShouldRun()
        if err != nil {
            continue
        }
        p.logFunc("LogProc: " + e.ToString())

        if p.OutQ != nil {
            log.Debug("LogProc Pushing")
            p.OutQ<-e
        }

        // Stats
        p.StatsAddMesg()
        p.PrintStats()

    }

    log.Info("LogProc Stopping!?")
}
