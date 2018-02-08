/*
   - CAST: Change the type of a field. Supported targets are int, float and string
*/
package proc

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering CastProc")
	core.GetRegistryInstance()["CastProc"] = NewCastProc
}

type CastProc struct {
	*core.ComponentBase
	Fields []string
	Types  []string
}

func NewCastProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating CastProc")

	fields := []string{}
	if tmp, ok := cfg["fields"].([]interface{}); ok {
		fields = core.InterfaceToStringArray(tmp)
	}

	types := []string{}
	if tmp, ok := cfg["types"].([]interface{}); ok {
		types = core.InterfaceToStringArray(tmp)
	}

	m := &CastProc{core.NewComponentBase(inQ, outQ, cfg), fields, types}
	m.Tag = "CAST-LOG"
	return m
}

func (p *CastProc) Signal(string) {}

func (p *CastProc) Run() {

	p.MustStop = false
	for !p.MustStop {

		e, err := p.ShouldRun()
		if err != nil {
			continue
		}

		for index, field := range p.Fields {
			value, ok := e.Data[field]
			if !ok {
				continue
			}

			switch p.Types[index] {
			case "string":
				fallthrough
			case "str":
				e.Data[field] = fmt.Sprintf("%v", value)

			case "int":
				switch v := value.(type) {
				case int64:
				case int:
					e.Data[field] = int64(v)
				case int8:
					e.Data[field] = int64(v)
				case int16:
					e.Data[field] = int64(v)
				case int32:
					e.Data[field] = int64(v)
				case float32:
					e.Data[field] = int64(v)
				case float64:
					e.Data[field] = int64(v)
				default:
					if vparse, err := strconv.ParseInt(fmt.Sprintf("%v", v), 0, 64); err == nil {
						e.Data[field] = vparse
					} else if vparse, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64); err == nil {
						e.Data[field] = int64(vparse)
					}
				}
			case "float":
				switch v := value.(type) {
				case float64:
				case int:
					e.Data[field] = float64(v)
				case int8:
					e.Data[field] = float64(v)
				case int16:
					e.Data[field] = float64(v)
				case int32:
					e.Data[field] = float64(v)
				case int64:
					e.Data[field] = float64(v)
				case float32:
					e.Data[field] = float64(v)
				default:
					if vparse, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64); err == nil {
						e.Data[field] = vparse
					}
				}
			}
		}

		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("CastProc Stopping!?")
}
