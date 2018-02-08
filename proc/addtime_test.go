package proc

import (
    "testing"
    . "github.com/urban-1/gopipe/tests"
)


func TestAddTime(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewAddTimeProc(in, out, GetConfig(`{"field_name":"ts"}`))
    go comp.Run()

    e := <-out
    if _, ok := e.Data["ts"]; !ok {
        t.Error("AddTime did not add anything!!")
    }
}


func TestAddTimeShouldNotRun(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEventRun(`{"doesnt": "matter"}`, false)

    comp := NewAddTimeProc(in, out, GetConfig(`{"field_name":"ts"}`))
    go comp.Run()

    e := <-out
	if _, ok := e.Data["test"]; ok {
	    // Has the new value!!! raise error
        t.Error("AddTime: run when it shouldnt...")
		t.Error(e.Data)
	}
}

func TestAddTimeShouldRun(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEventRun(`{"doesnt": "matter"}`, true)

    comp := NewAddTimeProc(in, out, GetConfig(`{"field_name":"ts"}`))
    go comp.Run()

    e := <-out
	if _, ok := e.Data["test"]; ok {
	    // Has the new value!!! raise error
        t.Error("AddTime: didn't run when it should...")
		t.Error(e.Data)
	}
}
