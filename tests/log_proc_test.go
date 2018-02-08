package main

import (
    "fmt"
    "testing"
    "encoding/json"

    . "gopipe/core"
    "gopipe/proc"
)

func GetChannels() (chan *Event, chan *Event){

    in := make(chan *Event, 1)
    out := make(chan *Event, 1)

    return in, out
}

func GetConfig(s string) Config {
    cfg := map[string]interface{}{}
    err := json.Unmarshal([]byte(s), &cfg)
    if err != nil {
        fmt.Printf("User error: cannot create mock config")
        panic("User error: cannot create mock config")
    }

    return cfg
}

func GetEvent(s string) *Event {
    e := map[string]interface{}{}
    err := json.Unmarshal([]byte(s), &e)
    if err != nil {
        fmt.Printf("User error: cannot create mock event")
        panic("User error: cannot create mock event")
    }

    return NewEvent(e)
}

func TestLogProc(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    pl := proc.NewLogProc(in, out, GetConfig(`{"level": "info"}`))
    go pl.Run()

    e := <-in
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}
