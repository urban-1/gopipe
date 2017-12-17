//
// - Events: The core data representation of this framework. These are passed
// around between components using Go channels
//
package core

import (
    "time"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)


// Basic event struct containing the Data and the received time
type Event struct {
    Timestamp time.Time
    Data map[string]interface{}
}

// Create a new event with the given data
func NewEvent(data map[string]interface{}) *Event {
    return &Event{time.Now(), data}
}

// Get the string replresentation of this event
func (e *Event) ToString() string {
    b, err := json.Marshal(e.Data)
    if err != nil {
        log.Error("Invalid JSON while converting event to string...")
        return ""
    }

    return string(b)
}

// Get the []byte representation of this event
func (e *Event) GetBytes() []byte{
    b, err := json.Marshal(e.Data)
    if err != nil {
        log.Error("Invalid JSON while converting event to string...")
        return []byte{}
    }

    return b
}
