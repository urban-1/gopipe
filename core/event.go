package core

import (
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

const (
    EVENT_STR = iota
    EVENT_DATA = iota
)


type Event interface {
    Type() int
    ToString() string
}

type DataEvent struct {
    mode int
    Data map[string]interface{}
}

type StrEvent struct {
    mode int
    Message string
}

func NewDataEvent(data map[string]interface{}) *DataEvent {
    // m := new(Event)
    // m.Mode = TYPE_STR
    // m.Message = ""
    // return m
    return &DataEvent{EVENT_DATA, data}
}

func (e *DataEvent) ToString() string {
    b, err := json.Marshal(e.Data)
    if err != nil {
        log.Error("Invalid JSON while converting event to string...")
        return ""
    }

    return string(b)
}

func (e *DataEvent) Type() int {
    return e.mode
}



func NewStrEvent(s string) *StrEvent {
    return &StrEvent{EVENT_STR, s}
}

func (e *StrEvent) ToString() string {
    return e.Message
}

func (e *StrEvent) Type() int {
    return e.mode
}
