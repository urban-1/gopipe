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
    config core.Config
    inQ chan core.Event
    outQ chan core.Event
    mustStop bool
}

func NewNullOutput(inQ chan core.Event, outQ chan core.Event, cfg core.Config) core.Component {
    log.Info("Creating NullOutput")
    return &NullOutput{cfg, inQ, outQ, false}
}

func (p *NullOutput) Stop() {
    p.mustStop = true
}

func (p *NullOutput) Run() {
    p.mustStop = false
    log.Debug("NullOutput Starting ... ")
    for !p.mustStop {
        log.Debug("NullOutput Reading")
        <-p.inQ
    }
    log.Debug("NullOutput Stopping")
}
