package proc

import (
	"encoding/json"
	. "github.com/urban-1/gopipe/tests"
	"reflect"
	"testing"
)

func TestCastIntFromString(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": "1"}`)

	comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["int"]}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["a"]).Name())
	}
}

func TestCastIntFromNumber(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": 1}`)

	comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["int"]}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["a"]).Name())
	}
}

func TestCastIntFromIntFloat(t *testing.T) {
	in, out := GetChannels()
	ein := GetEvent(`{"a": 1, "b": 1, "c": 1, "d": 1}`)
	ein.Data["a"], _ = ein.Data["a"].(json.Number).Int64()
	ein.Data["a"] = int32(ein.Data["a"].(int64))
	ein.Data["b"] = int16(ein.Data["a"].(int32))
	ein.Data["c"] = int8(ein.Data["a"].(int32))
	ein.Data["d"] = int(ein.Data["a"].(int32))
	ein.Data["e"] = float64(ein.Data["a"].(int32))
	ein.Data["f"] = float32(ein.Data["a"].(int32))
	in <- ein

	comp := NewCastProc(in, out, GetConfig(`{
		"fields":["a", "b", "c", "d", "e", "f"],
		"types": ["int", "int", "int", "int", "int", "int"]
	}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["a"]).Name())
	}
	if reflect.TypeOf(e.Data["b"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["b"]).Name())
	}
	if reflect.TypeOf(e.Data["c"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["c"]).Name())
	}
	if reflect.TypeOf(e.Data["d"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["d"]).Name())
	}
	if reflect.TypeOf(e.Data["e"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["c"]).Name())
	}
	if reflect.TypeOf(e.Data["f"]).Name() != "int64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["d"]).Name())
	}
}

func TestCastFloatFromIntFloat(t *testing.T) {
	in, out := GetChannels()
	ein := GetEvent(`{"a": 1, "b": 1, "c": 1, "d": 1}`)
	ein.Data["a"], _ = ein.Data["a"].(json.Number).Float64()
	ein.Data["a"] = int64(ein.Data["a"].(float64))
	ein.Data["b"] = int32(ein.Data["a"].(int64))
	ein.Data["c"] = int16(ein.Data["a"].(int64))
	ein.Data["d"] = int8(ein.Data["a"].(int64))
	ein.Data["e"] = int(ein.Data["a"].(int64))
	ein.Data["f"] = float32(ein.Data["a"].(int64))
	in <- ein

	comp := NewCastProc(in, out, GetConfig(`{
		"fields":["a", "b", "c", "d", "e", "f"],
		"types": ["float", "float", "float", "float", "float", "float"]
	}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["a"]).Name())
	}
	if reflect.TypeOf(e.Data["b"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["b"]).Name())
	}
	if reflect.TypeOf(e.Data["c"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["c"]).Name())
	}
	if reflect.TypeOf(e.Data["d"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["d"]).Name())
	}
	if reflect.TypeOf(e.Data["e"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["c"]).Name())
	}
	if reflect.TypeOf(e.Data["f"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["d"]).Name())
	}
}

func TestCastFloatFromString(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": "1.0"}`)

	comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["float"]}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["a"]).Name())
	}
}

func TestCastFloatFromNumber(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": 1.0}`)

	comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["float"]}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() != "float64" {
		t.Error("Cast did not add anything!! Type is:")
		t.Error(reflect.TypeOf(e.Data["a"]).Name())
	}
}

func TestCastShouldNotRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": "1"}`, false)

	comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["int"]}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() == "int64" {
		t.Error("Cast run... it shouldn't")
	}
}

func TestCastShouldRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": "1"}`, true)

	comp := NewCastProc(in, out, GetConfig(`{"fields":["a"], "types": ["int"]}`))
	go comp.Run()

	e := <-out
	if reflect.TypeOf(e.Data["a"]).Name() != "int64" {
		// Has the new value!!! raise error
		t.Error("Cast: didn't run when it should...")
	}
}
