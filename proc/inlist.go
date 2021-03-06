/*
   - INLIST: Check if a field of the event exists in a list and store the result
   (true/false) into another field. The list we check against can be provided by
   config (static) or can be regularly read from a file. In any case, the items
   in the list are strings and thus every the data field is converted to string
   to be checked against the list. The main function/purpose of this plugin is
   to verify against lists that change regularly (ex IP blacklist) and thus the
   analysis has to take place at the correct time (the time of the event) and
   cannot be performed in later time!
*/
package proc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urban-1/gopipe/core"
)

func init() {
	log.Info("Registering InListProc")
	core.GetRegistryInstance()["InListProc"] = NewInListProc
}

type InListProc struct {
	*core.ComponentBase
	List          map[string]bool
	ListLock      *sync.Mutex
	FilePath      string
	InField       string
	OutField      string
	ReloadMinutes int
}

func NewInListProc(inQ chan *core.Event, outQ chan *core.Event, cfg core.Config) core.Component {
	log.Info("Creating InListProc")

	// Set this modules log level

	list := map[string]bool{}
	tmp, ok1 := cfg["list"].([]interface{})

	if ok1 {
		// Convert list to map for faster lookup
		for _, item := range tmp {
			list[item.(string)] = true
		}
	} else {
		list = nil
	}

	fpath, ok2 := cfg["filepath"].(string)

	// If none provided, complain
	if !ok1 && !ok2 {
		panic("InList requires a list either in 'list' or in 'filepath'")
	}

	// If filepath provided, ignore config list
	if ok2 {
		log.Info("Clearing config list since file path is present")
		list = nil
	}

	r, ok := cfg["reload_minutes"].(float64)
	reload := -1
	if ok {
		reload = int(r)
	}

	m := &InListProc{core.NewComponentBase(inQ, outQ, cfg),
		list, &sync.Mutex{}, fpath,
		cfg["in_field"].(string),
		cfg["out_field"].(string),
		reload}

	m.Tag = "PROC-INLIST"
	return m
}

func (p *InListProc) Signal(signal string) {
	log.Infof("InListProc Received signal '%s'", signal)
	switch signal {
	case "reload":
		if p.FilePath == "" {
			log.Error("InListProc IGNORING signal 'reload' - no filepath configured")
			return
		}
		p.loadList()
	default:
		log.Infof("InListProc UNKNOW signal '%s'", signal)
	}
}

func (p *InListProc) Run() {
	log.Debug("InListProc Starting ... ")

	// Spawn the loader
	if p.List == nil {
		if p.ReloadMinutes > 0 {
			// Periodic reloading
			go func(p *InListProc) {
				p.loadList()
				time.Sleep(time.Duration(p.ReloadMinutes) * time.Minute)
			}(p)
		} else {
			// Once off loading...
			p.loadList()
		}
	}

	p.MustStop = false
	cfg_error := false

	for !p.MustStop {
		log.Debug("InListProc Reading")
		e, err := p.ShouldRun()
		if err != nil {
			continue
		}

		what, ok := e.Data[p.InField]
		if !ok {
			// This is a user error, maybe error once?
			if !cfg_error {
				log.Error("Cannot find field ", p.InField)
				cfg_error = true
			}
		}

		whatstr := fmt.Sprintf("%v", what)
		p.ListLock.Lock()
		if _, ok := p.List[whatstr]; ok {
			e.Data[p.OutField] = true
		} else {
			e.Data[p.OutField] = false
		}
		p.ListLock.Unlock()

		p.OutQ <- e

		// Stats
		p.StatsAddMesg()
		p.PrintStats()

	}

	log.Info("InListProc Stopping!?")
}

func (p *InListProc) loadList() {
	f, err := os.Open(p.FilePath)
	if err != nil {
		log.Error("INLIST: Could not load list file")
		return
	}

	p.ListLock.Lock()

	log.Warn("INLIST: Reading file")
	reader := bufio.NewReader(f)

	count := 0

	p.List = map[string]bool{}

	line, _, err := reader.ReadLine()
	for err != io.EOF {
		if string(line) == "" {
			continue
		}
		p.List[string(line)] = true
		count += 1
		line, _, err = reader.ReadLine()
	}

	log.Info("INLIST: Done! Loaded ", count, " items!")

	f.Close()
	p.ListLock.Unlock()
}
