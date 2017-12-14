package output

import (
    "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering NullOutput")
    core.GetRegistryInstance()["NullOutput"] = NewNullOutput
}

type NullOutput struct {
    core.ComponentBase
}

func NewNullOutput(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
    log.Info("Creating NullOutput")
    return &NullOutput{*core.NewComponentBase(inQ, outQ, cfg)}
}

func (p *NullOutput) Stop() {
    p.MustStop = true
}

func (p *NullOutput) Run() {
    p.MustStop = false
    log.Debug("NullOutput Starting ... ")
    for !p.MustStop {
        log.Debug("NullOutput Reading")
        <-p.InQ

        // Stats
        p.StatsAddMesg()
        p.PrintStats("Null", 50000)
    }
    log.Debug("NullOutput Stopping")
}
