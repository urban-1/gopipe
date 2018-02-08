package proc

import (
    "testing"
    . "github.com/urban-1/gopipe/tests"
)


func TestDropField(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"a": "1", "b": 100}`)

    comp := NewDropFieldProc(in, out, GetConfig(`{"field_name":"a"}`))
    go comp.Run()

    e := <-out
    if _, ok := e.Data["a"]; ok {
        t.Error("DropField did not drop anything!! Data is:")
		t.Error(e.Data)
    }
}


func TestDropFieldShouldNotRun(t *testing.T) {
	in,  out := GetChannels()
    in <- GetEventRun(`{"a": "1"}`, false)

    comp := NewDropFieldProc(in, out, GetConfig(`{"field_name":"a"}`))
    go comp.Run()

    e := <-out
    if _, ok := e.Data["a"]; !ok {
        t.Error("DropField run... it shouldn't")
    }
}

func TestDropFieldShouldRun(t *testing.T) {
	in,  out := GetChannels()
    in <- GetEventRun(`{"a": "1"}`, true)

    comp := NewDropFieldProc(in, out, GetConfig(`{"field_name":"a"}`))
    go comp.Run()

    e := <-out
    if _, ok := e.Data["a"]; ok {
	    // Has the new value!!! raise error
        t.Error("DropField: didn't run when it should...")
	}
}
