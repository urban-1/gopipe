/*
   The package having the processing components.

   Every processing component reads, modifies messages and pushes to the output
   channel.

   - ADD: Add a field to the event's data. TODO: Add expression support if easy
   instead of static value
*/
package proc

import (
	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering AddFieldProc")
	core.GetRegistryInstance()["AddFieldProc"] = NewAddFieldProc
}

type AddFieldProc struct {
	*core.ComponentBase
	FieldName string
	Value     interface{}
	Expr      *govaluate.EvaluableExpression
}

func NewAddFieldProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating AddFieldProc")
	field_name, ok := cfg["field_name"].(string)
	if !ok {
		panic("ADD Field: Field name is required")
	}

	strexp, ok := cfg["expression"].(string)
	if !ok {
		strexp = ""
	}
	value, ok := cfg["value"]

	expression, err := govaluate.NewEvaluableExpression(strexp)
	if err != nil && !ok {
		panic("Add field: Either expression or value is required")
	}

	if err == nil {
		value = nil
	}

	m := &AddFieldProc{core.NewComponentBase(inQ, outQ, cfg), field_name, value, expression}
	m.Tag = "PROC-ADDFIELD"
	return m
}

func (p *AddFieldProc) Signal(string) {}

// TODO: Add expression support if easy
func (p *AddFieldProc) Run() {
	log.Debug("AddFieldProc Starting ... ")
	p.MustStop = false
	for !p.MustStop {
		log.Debug("AddFieldProc Reading")
		e, err := p.ShouldRun()
		if err != nil {
			continue
		}

		if p.Value == nil {
			result, err := p.Expr.Evaluate(e.Data)
			log.Debug("AddFieldProc EXPR")
			if err != nil {
				log.Warn(p.Tag, ": ", err.Error())
			}
			e.Data[p.FieldName] = result
		} else {
			log.Debug("AddFieldProc VAL")
			e.Data[p.FieldName] = p.Value
		}
		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("AddFieldProc Stopping!?")
}
