/*
   - LOG: Use the logger (logrus) to print the event. This component does not
   modify the events in any way. NOTE: This can also be used as an OUTPUT component!
*/
package proc

import (
	"github.com/Knetic/govaluate"
	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering IfProc")
	core.GetRegistryInstance()["if"] = NewIfProc

	log.Info("Registering ElseProc")
	core.GetRegistryInstance()["else"] = NewElseProc

	log.Info("Registering EndIfProc")
	core.GetRegistryInstance()["endif"] = NewEndIfProc
}

type IfProc struct {
	*core.ComponentBase
	Expr *govaluate.EvaluableExpression
}

func NewIfProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating IfProc")

	// Set this modules log level
	cond, ok := cfg["condition"].(string)
	if !ok {
		panic("If module needs a condition")
	}

	expression, err := govaluate.NewEvaluableExpression(cond)
	if err != nil {
		panic("If module failed to evaluate condition")
	}

	m := &IfProc{core.NewComponentBase(inQ, outQ, cfg), expression}
	m.Tag = "PROC-IF"
	return m
}

func (p *IfProc) Signal(string) {}

// The if module modifies the BoolStack of an event. This is used to flag skips
func (p *IfProc) Run() {
	log.Debug("IfProc Starting ... ")
	p.MustStop = false
	for !p.MustStop {
		log.Debug("IfProc Reading")
		e := <-p.InQ

		// Evaluate the expression against the data of the event!
		result, err := p.Expr.Evaluate(e.Data)
		if err != nil {
			log.Warn(p.Tag, ": ", err.Error())
		}

		log.Debug("If expression evaluates to: ", result)

		// We now push the result to the stack:
		// -  if we are in a false group, push false
		if e.ShouldRun.Size() > 0 {
			state, _ := e.ShouldRun.Top()
			if !state {
				// Override to true so we push a false at the end...
				result = false
			}
		}

		//  - if the condition was true then ShouldRun=true
		//  - else ShouldRun=false
		e.ShouldRun.Push(result.(bool))

		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("IfProc Stopping!?")
}

type ElseProc struct {
	*core.ComponentBase
}

func NewElseProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating ElseProc")
	m := &ElseProc{core.NewComponentBase(inQ, outQ, cfg)}
	m.Tag = "PROC-ELSE"
	return m
}

func (p *ElseProc) Signal(string) {}

// The else module reverse the effect of the if module!
func (p *ElseProc) Run() {
	log.Debug("ElseProc Starting ... ")
	p.MustStop = false
	for !p.MustStop {
		log.Debug("ElseProc Reading")
		e := <-p.InQ

		result, err := e.ShouldRun.Pop()
		if err != nil {
			log.Error("User configuration error... Check your IF/ELSEs")
		}

		// Check parent state, if false push false
		if e.ShouldRun.Size() > 0 {
			state, _ := e.ShouldRun.Top()
			if !state {
				result = true
			}
		}
		e.ShouldRun.Push(!result)

		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("ElseProc Stopping!?")
}

type EndIfProc struct {
	*core.ComponentBase
}

func NewEndIfProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating EndIfProc")
	m := &EndIfProc{core.NewComponentBase(inQ, outQ, cfg)}
	m.Tag = "PROC-ENDIF"
	return m
}

func (p *EndIfProc) Signal(string) {}

// The endif module just removes the current state from the ShouldRun stack
func (p *EndIfProc) Run() {
	log.Debug("EndIfProc Starting ... ")
	p.MustStop = false
	for !p.MustStop {
		log.Debug("EndIfProc Reading")
		e := <-p.InQ

		// In any case pop one out...
		_, _ = e.ShouldRun.Pop()

		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("EndIfProc Stopping!?")
}
