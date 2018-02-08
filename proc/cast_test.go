package proc

import (
    "testing"
	"reflect"
    . "github.com/urban-1/gopipe/tests"
)


func TestCast(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"a": "1"}`)

    comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["int"]}`))
    go comp.Run()

    e := <-out
    if reflect.TypeOf(e.Data["a"]).Name() != "int64" {
        t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["a"]).Name())
    }
}


func TestCastShouldNotRun(t *testing.T) {
	in,  out := GetChannels()
    in <- GetEventRun(`{"a": "1"}`, false)

    comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["int"]}`))
    go comp.Run()

    e := <-out
    if reflect.TypeOf(e.Data["a"]).Name() == "int64" {
        t.Error("Cast run... it shouldn't")
    }
}

func TestCastShouldRun(t *testing.T) {
	in,  out := GetChannels()
    in <- GetEventRun(`{"a": "1"}`, true)

    comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["int"]}`))
    go comp.Run()

    e := <-out
    if reflect.TypeOf(e.Data["a"]).Name() != "int64" {
	    // Has the new value!!! raise error
        t.Error("Cast: didn't run when it should...")
	}
}
