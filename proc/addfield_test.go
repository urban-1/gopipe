package proc

import (
	. "github.com/urban-1/gopipe/tests"
	"testing"
)

func TestAddFieldValue(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"doesnt": "matter"}`)

	comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"value": "blah",
			"field_name": "test"
		}
	`))
	go comp.Run()

	e := <-out
	if e.Data["test"] != "blah" {
		t.Error("AddField did not add anything!!")
	}
}

func TestAddFieldExpression(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"doesnt": "matter"}`)

	comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"expression": "doesnt",
			"field_name": "test"
		}
	`))
	go comp.Run()

	e := <-out
	if e.Data["test"] != "matter" {
		t.Error("AddField: I was expecting test: matter")
		t.Error(e.Data)
	}
}

func TestAddFieldExpression2(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"doesnt": "matter"}`)

	comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"expression": "10*10+3",
			"field_name": "test"
		}
	`))
	go comp.Run()

	e := <-out
	if int(e.Data["test"].(float64)) != 103 {
		t.Error("AddField: I was expecting test:103...")
		t.Error(e.Data)
	}
}

func TestAddFieldShouldNotRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"doesnt": "matter"}`, false)

	comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"expression": "10*10+3",
			"field_name": "test"
		}
	`))
	go comp.Run()

	e := <-out
	if _, ok := e.Data["test"]; ok {
		// Has the new value!!! raise error
		t.Error("AddField: run when it shouldnt...")
		t.Error(e.Data)
	}
}

func TestAddFieldShouldRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"doesnt": "matter"}`, true)

	comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"value": "matter",
			"field_name": "dark"
		}
	`))
	go comp.Run()

	e := <-out
	if e.Data["dark"] != "matter" {
		t.Error("AddField: I was expecting test: matter")
		t.Error(e.Data)
	}
}
