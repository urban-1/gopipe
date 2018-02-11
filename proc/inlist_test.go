package proc

import (
	log "github.com/sirupsen/logrus"
	. "github.com/urban-1/gopipe/core"
	. "github.com/urban-1/gopipe/tests"
	"testing"
)

func getInList(in chan *Event, out chan *Event) Component {
	comp := NewInListProc(in, out, GetConfig(`{
		"in_field": "port",
		"out_field": "port_block",
		"reload_minutes": -1,
		"list": ["8080", "443", "23230", "14572", "17018", "1.9"]
	}`))

	return comp
}

func TestInListProc(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"port": "443"}`)

	comp := getInList(in, out)
	go comp.Run()

	e := <-out
	if !e.Data["port_block"].(bool) {
		t.Error("InList did not match!")
		log.Error(e.Data)
	}
}

func TestInListProcIntStr(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"port": 443}`)

	comp := getInList(in, out)
	go comp.Run()

	e := <-out
	if !e.Data["port_block"].(bool) {
		t.Error("InList did not match!")
		log.Error(e.Data)
	}
}

func TestInListProcFloatStr(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"port": 1.9}`)

	comp := getInList(in, out)
	go comp.Run()

	e := <-out
	if !e.Data["port_block"].(bool) {
		t.Error("InList did not match!")
		log.Error(e.Data)
	}
}

func TestInListShouldNotRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"port": "443"}`, false)

	comp := getInList(in, out)
	go comp.Run()

	e := <-out
	if _, ok := e.Data["port_block"]; ok {
		t.Error("InList run when it shouldn't")
		log.Error(e.Data)
	}
}

func TestInListShouldRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"port": "443"}`, true)

	comp := getInList(in, out)
	go comp.Run()

	e := <-out
	if !e.Data["port_block"].(bool) {
		t.Error("InList did not run when it should...")
		log.Error(e.Data)
	}
}
