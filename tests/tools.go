package tests

import (
    "fmt"
    "encoding/json"

    . "github.com/urban-1/gopipe/core"

    log "github.com/sirupsen/logrus"
)


func init() {
    log.SetLevel(log.DebugLevel)
}

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

func GetEventRun(s string, run bool) *Event {
    e := map[string]interface{}{}
    err := json.Unmarshal([]byte(s), &e)
    if err != nil {
        fmt.Printf("User error: cannot create mock event")
        panic("User error: cannot create mock event")
    }

    ret := NewEvent(e)
	ret.ShouldRun.Push(run)
	return ret
}
