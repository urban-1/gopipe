package output

import (
	. "github.com/urban-1/gopipe/tests"
	"testing"
	"time"
)

func TestNull(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"doesnt": "matter"}`)

	comp := NewNullOutput(in, out, GetConfig(`{}`))
	go comp.Run()
	time.Sleep(time.Duration(1) * time.Second)

	if len(out) > 0 {
		t.Error("NullProc did not blackhole!")
	}
}

func TestNullShouldNotRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"doesnt": "matter"}`, false)

	comp := NewNullOutput(in, out, GetConfig(`{}`))
	go comp.Run()
	time.Sleep(time.Duration(1) * time.Second)

	if len(out) < 1 {
		t.Error("NullProc run when it shouldn't")
		t.Error("Out Q: ", len(out))
	}
}

func TestNullShouldRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"doesnt": "matter"}`, true)

	comp := NewNullOutput(in, out, GetConfig(`{}`))
	go comp.Run()
	time.Sleep(time.Duration(1) * time.Second)

	if len(out) > 0 {
		t.Error("NullProc did not run, it should...")
	}
}
