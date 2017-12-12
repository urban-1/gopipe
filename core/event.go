package core

import (
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

const (
    EVENT_STR = iota
    EVENT_DATA = iota
)

type Event struct {
    mode int
}

func (e *Event) Type() int {
    return e.mode
}

type DataEvent struct {
    Event
    Data map[string]interface{}
}

type StrEvent struct {
    Event
    Message string
}

func NewDataEvent() *DataEvent {
    // m := new(Event)
    // m.Mode = TYPE_STR
    // m.Message = ""
    // return m
    return &DataEvent{Event{EVENT_DATA}, nil}
}


func (e *DataEvent) toString() string {
    b, err := json.Marshal(e.Data)
    if err != nil {
        log.Error("Invalid JSON while converting event to string...")
        return ""
    }

    return string(b)
}

func NewStrEvent() *StrEvent {
    return &StrEvent{Event{EVENT_STR}, ""}
}

func (e *StrEvent) toString() string {
    return e.Message
}
