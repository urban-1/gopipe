/*
    - LOG: Use the logger (logrus) to print the event. This component does not
    modify the events in any way. NOTE: This can also be used as an OUTPUT component!
 */
package proc

import (
	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

// Register our component when the module is included
func init() {
	log.Info("Registering LogProc")
	core.GetRegistryInstance()["LogProc"] = NewLogProc
}

// Base struct "extending" ComponentBase
type LogProc struct {
	*core.ComponentBase
	logFunc func(args ...interface{})
}

// Constructor function
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

	// Create an instance
	m := &LogProc{core.NewComponentBase(inQ, outQ, cfg), logFunc}

	// Assign a unique tag to this component. This is used mainly for logging
	m.Tag = "PROC-LOG"

	// instance ready...
	return m
}

// Handle signals for this component
func (p *LogProc) Signal(string) {}

// Our component's main function 
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
			p.OutQ <- e
		}

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("LogProc Stopping!?")
}
