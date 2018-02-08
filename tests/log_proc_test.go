package main

import (
    "testing"
    "gopipe/proc"

)


func TestLogProc(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := proc.NewLogProc(in, out, GetConfig(`{"level": "error"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}
