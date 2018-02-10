package proc

import (
	. "github.com/urban-1/gopipe/tests"
	"testing"
)

func TestIf(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": 1.0}`)

	// This is MENTAL: TODO: Add a warning in the docs once you figure it out...
	// ... seems to fundamentally broken :(
	// comp := NewIfProc(in, out, GetConfig(`{"condition": "json_to_int64(a) == 1"}`))
	comp := NewIfProc(in, out, GetConfig(`{"condition": "json_to_float64(a) == 1"}`))
	// comp := NewIfProc(in, out, GetConfig(`{"condition": "json_to_float64(a) < 2"}`))
	go comp.Run()

	e := <-out
	shouldRun, err := e.ShouldRun.Top()
	if err != nil {
		t.Error("If did not push a value:")
		t.Error(e.ShouldRun)
	}
	if !shouldRun {
		t.Error("If did not add expected value (true)")
		t.Error(e.ShouldRun)
	}
}

func TestIfStringMistake(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": "1"}`)

	comp := NewIfProc(in, out, GetConfig(`{"condition": "a == 1"}`))
	go comp.Run()

	e := <-out
	shouldRun, err := e.ShouldRun.Top()
	if err != nil {
		t.Error("If did not push a value:")
		t.Error(e.ShouldRun)
	}
	if shouldRun {
		t.Error("If did not add expected value (false)")
		t.Error(e.ShouldRun)
	}
}

// Correct condition, but disabled due to parent block
func TestIfParentDisabled(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": 1}`, false)

	comp := NewIfProc(in, out, GetConfig(`{"condition": "a == 1"}`))
	go comp.Run()

	e := <-out
	shouldRun, _ := e.ShouldRun.Top()
	if e.ShouldRun.Size() < 2 {
		t.Error("If did not push a value:")
		t.Error(e.ShouldRun)
	}
	if shouldRun {
		t.Error("If did not add expected value (false)")
		t.Error(e.ShouldRun)
	}
}

// Correct condition and enabled parent block
func TestIfParentEnabled(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": 1}`, true)

	comp := NewIfProc(in, out, GetConfig(`{"condition": "a == 1"}`))
	go comp.Run()

	e := <-out
	shouldRun, _ := e.ShouldRun.Top()
	if e.ShouldRun.Size() < 2 {
		t.Error("If did not push a value:")
		t.Error(e.ShouldRun)
	}
	if !shouldRun {
		t.Error("If did not add expected value (true)")
		t.Error(e.ShouldRun)
	}
}

// -----------------------------------------------------------------------------

func TestElse(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": 1}`, false)

	comp := NewElseProc(in, out, GetConfig(`{}`))
	go comp.Run()

	e := <-out
	shouldRun, _ := e.ShouldRun.Top()
	if !shouldRun {
		t.Error("Else did not add expected value (true)")
		t.Error(e.ShouldRun)
	}
}

func TestElse2(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": 1}`, true)

	comp := NewElseProc(in, out, GetConfig(`{}`))
	go comp.Run()

	e := <-out
	shouldRun, _ := e.ShouldRun.Top()
	if shouldRun {
		t.Error("Else did not add expected value (false)")
		t.Error(e.ShouldRun)
	}
}


// Correct condition, but disabled due to parent block
func TestElseParentDisabled(t *testing.T) {
	in, out := GetChannels()

	ein := GetEventRun(`{"a": 1}`, false)

	// Current block state...
	ein.ShouldRun.Push(false)

	in <- ein

	comp := NewElseProc(in, out, GetConfig(`{}`))
	go comp.Run()

	e := <-out
	shouldRun, _ := e.ShouldRun.Top()
	if e.ShouldRun.Size() < 2 {
		t.Error("Else poped a value?!")
		t.Error(e.ShouldRun)
	}
	if shouldRun {
		t.Error("Else did not add expected value (false)")
		t.Error(e.ShouldRun)
	}
}

// Correct condition, but disabled due to parent block
func TestElseParentEnabled(t *testing.T) {
	in, out := GetChannels()

	ein := GetEventRun(`{"a": 1}`, true)

	// Current block state...
	ein.ShouldRun.Push(false)

	in <- ein

	comp := NewElseProc(in, out, GetConfig(`{}`))
	go comp.Run()

	e := <-out
	shouldRun, _ := e.ShouldRun.Top()
	if e.ShouldRun.Size() < 2 {
		t.Error("Else poped a value?!")
		t.Error(e.ShouldRun)
	}
	if !shouldRun {
		t.Error("Else did not add expected value (false)")
		t.Error(e.ShouldRun)
	}
}

// -----------------------------------------------------------------------------

func TestEndIf(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"a": 1}`, false)

	comp := NewEndIfProc(in, out, GetConfig(`{}`))
	go comp.Run()

	e := <-out
	if e.ShouldRun.Size() > 0 {
		t.Error("End if did not remove from stack")
		t.Error(e.ShouldRun)
	}
}

func TestEndIfConfigError(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"a": 1}`)

	comp := NewEndIfProc(in, out, GetConfig(`{}`))
	go comp.Run()

	e := <-out
	if e.ShouldRun.Size() > 0 {
		t.Error("End if did not remove from stack")
		t.Error(e.ShouldRun)
	}
}
