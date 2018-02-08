/*
   - MD5: Hash a set of given data fields and attach the results back to the
   event's data. Optionally, a salt can be provided
*/
package proc

import (
	"crypto/md5"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering Md5Proc")
	core.GetRegistryInstance()["Md5Proc"] = NewMd5Proc
}

type Md5Proc struct {
	*core.ComponentBase
	InFields  []string
	OutFields []string
	Salt      string
}

func NewMd5Proc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating Md5Proc")

	in_fields := []string{}
	if tmp, ok := cfg["in_fields"].([]interface{}); ok {
		in_fields = core.InterfaceToStringArray(tmp)
	}

	out_fields := []string{}
	if tmp, ok := cfg["out_fields"].([]interface{}); ok {
		out_fields = core.InterfaceToStringArray(tmp)
	}

	salt, ok := cfg["salt"].(string)
	if !ok {
		salt = ""
	}

	m := &Md5Proc{core.NewComponentBase(inQ, outQ, cfg), in_fields, out_fields, salt}
	m.Tag = "MD5-LOG"
	return m
}

func (p *Md5Proc) Signal(string) {}

func (p *Md5Proc) Run() {
	log.Debug("Md5Proc Starting ... ")
	p.MustStop = false

	for !p.MustStop {
		// Check if we should run (based on the events' if else state)
		e, err := p.ShouldRun()
		if err != nil {
			continue
		}

		for i, ifield := range p.InFields {
			b, ok := e.Data[ifield].(string)
			if !ok {
				log.Error("Failed to convert field ", ifield, " to string...")
				continue
			}

			md5tmp := md5.Sum([]byte(b + p.Salt))
			e.Data[p.OutFields[i]] = hex.EncodeToString(md5tmp[:])
		}

		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("Md5Proc Stopping!?")
}
