package core

import (
    "time"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

const (
    EVENT_STR = iota
    EVENT_DATA = iota
)



type Event struct {
    mode int
    Timestamp time.Time
    Data map[string]interface{}
}


func NewEvent(data map[string]interface{}) *Event {
    // m := new(Event)
    // m.Mode = TYPE_STR
    // m.Message = ""
    // return m
    return &Event{EVENT_DATA, time.Now(), data}
}

func (e *Event) ToString() string {
    b, err := json.Marshal(e.Data)
    if err != nil {
        log.Error("Invalid JSON while converting event to string...")
        return ""
    }

    return string(b)
}

func (e *Event) Type() int {
    return e.mode
}

func (e *Event) GetBytes() []byte{
    b, err := json.Marshal(e.Data)
    if err != nil {
        log.Error("Invalid JSON while converting event to string...")
        return []byte{}
    }

    return b
}
