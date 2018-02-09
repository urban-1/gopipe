/*
   The package having the processing components.

   Every processing component reads, modifies messages and pushes to the output
   channel.

   - ADDTIME: This component gets the Event's timestamp and embeds it into the
   data (under the given key/name)
*/
package proc

import (
	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering AddTimeProc")
	core.GetRegistryInstance()["AddTimeProc"] = NewAddTimeProc
}

type AddTimeProc struct {
	*core.ComponentBase
	FieldName string
	InSeconds bool
}

func NewAddTimeProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating AddTimeProc")
	field_name, ok := cfg["field_name"].(string)
	if !ok {
		field_name = "timestamp"
	}

	in_s, ok := cfg["in_seconds"].(bool)
	if !ok {
		in_s = false
	}

	m := &AddTimeProc{core.NewComponentBase(inQ, outQ, cfg), field_name, in_s}
	m.Tag = "PROC-ADDTIME"
	return m
}

func (p *AddTimeProc) Signal(string) {}

func (p *AddTimeProc) Run() {
	log.Debug("AddTimeProc Starting ... ")
	p.MustStop = false

	factor := int64(1000000)
	if p.InSeconds {
		factor = 1000000000
	}
	for !p.MustStop {
		log.Debug("AddTimeProc Reading")
		e, err := p.ShouldRun()
		if err != nil {
			continue
		}

		e.Data[p.FieldName] = uint64(e.Timestamp.UnixNano() / factor)
		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("AddTimeProc Stopping!?")
}
