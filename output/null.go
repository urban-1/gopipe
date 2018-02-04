/*
    - NULL: Just empty that channel! NOTE: This component CANNOT be used in the
    processing stages, since whereever put, it acts like a blackhole...
 */
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
    m := &NullOutput{*core.NewComponentBase(inQ, outQ, cfg)}

    m.Tag = "OUT-NULL"

    return m
}

func  (p *NullOutput) Signal(string) {}

func (p *NullOutput) Run() {
    p.MustStop = false
    log.Debug("NullOutput Starting ... ")
    for !p.MustStop {
        log.Debug("NullOutput Reading")
        _, err := p.ShouldRun()
        if err != nil {
            continue
        }

        // Stats
        p.StatsAddMesg()
        p.PrintStats()
    }
    log.Debug("NullOutput Stopping")
}
