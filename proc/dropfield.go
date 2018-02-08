/*
    - DROP: Remove a field from the event's data
 */
package proc

import (
	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering DropFieldProc")
	core.GetRegistryInstance()["DropFieldProc"] = NewDropFieldProc
}

type DropFieldProc struct {
	*core.ComponentBase
	FieldName string
}

func NewDropFieldProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating DropFieldProc")
	field_name, ok := cfg["field_name"].(string)
	if !ok {
		field_name = "timestamp"
	}
	m := &DropFieldProc{core.NewComponentBase(inQ, outQ, cfg), field_name}
	m.Tag = "PROC-DROPFIELD"
	return m
}

func (p *DropFieldProc) Signal(string) {}

func (p *DropFieldProc) Run() {
	log.Debug("DropFieldProc Starting ... ")
	p.MustStop = false
	for !p.MustStop {
		log.Debug("DropFieldProc Reading")
		e, err := p.ShouldRun()
		if err != nil {
			continue
		}

		delete(e.Data, p.FieldName)
		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("DropFieldProc Stopping!?")
}
