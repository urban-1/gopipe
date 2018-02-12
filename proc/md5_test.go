package proc

import (
	. "github.com/urban-1/gopipe/tests"
	"testing"
)

func TestMd5(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": "1", "b": 100}`)

	comp := NewMd5Proc(in, out, GetConfig(`{
        "in_fields": ["a"],
        "out_fields": ["md5"],
        "salt": "test"
    }`))
	go comp.Run()

	e := <-out
	if _, ok := e.Data["md5"]; !ok {
		t.Error("Md5 did not run... Data is:")
		t.Error(e.Data)
	}
}

func TestMd5ShouldNotRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": "1", "b": 100}`, false)

	comp := NewMd5Proc(in, out, GetConfig(`{
        "in_fields": ["a"],
        "out_fields": ["md5"],
        "salt": "test"
    }`))
	go comp.Run()

	e := <-out
	if _, ok := e.Data["md5"]; ok {
		t.Error("Md5 run... it shouldn't")
	}
}

func TestMd5ShouldRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": "1", "b": 100}`, true)

	comp := NewMd5Proc(in, out, GetConfig(`{
        "in_fields": ["a"],
        "out_fields": ["md5"],
        "salt": "test"
    }`))
	go comp.Run()

	e := <-out
	if _, ok := e.Data["md5"]; !ok {
		// Has the new value!!! raise error
		t.Error("Md5: didn't run when it should...")
		t.Error(e.Data)
	}
}
