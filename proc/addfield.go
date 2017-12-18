/*
    The package having the processing components.

    Every processing component reads, modifies messages and pushes to the output
    channel.

    - ADD: Add a field to the event's data. TODO: Add expression support if easy
    instead of static value
 */
package proc

import (
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering AddFieldProc")
    GetRegistryInstance()["AddFieldProc"] = NewAddFieldProc
}

type AddFieldProc struct {
    *ComponentBase
    FieldName string
    Value interface{}
}

func NewAddFieldProc(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating AddFieldProc")
    field_name, ok := cfg["field_name"].(string)
    if !ok {
        panic("ADD Field: Field name is required")
    }
    value, ok := cfg["value"]
    if !ok {
        panic("ADD Field: Field value is required")
    }
    m := &AddFieldProc{NewComponentBase(inQ, outQ, cfg), field_name, value}
    m.Tag = "PROC-ADDFIELD"
    return m
}

// TODO: Add expression support if easy
func (p *AddFieldProc) Run() {
    log.Debug("AddFieldProc Starting ... ")
    p.MustStop = false
    for !p.MustStop {
        log.Debug("AddFieldProc Reading")
        e := <- p.InQ

        e.Data[p.FieldName] = p.Value
        p.OutQ<-e

        // Stats
        p.StatsAddMesg()
        p.PrintStats()

    }

    log.Info("AddFieldProc Stopping!?")
}
