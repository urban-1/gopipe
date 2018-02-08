package proc

import (
    "testing"
    . "github.com/urban-1/gopipe/tests"
)


func TestLogProcOutput(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "blah"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }


}

func TestLogProcInfo(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "info"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}


func TestLogProcDebug(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "debug"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}


func TestLogProcWarn(t *testing.T) {
    in,  out := GetChannels()
	in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "warn"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}
