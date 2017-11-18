package core

import (
    log "github.com/sirupsen/logrus"
)

// Create singleton registry
var registry map[string]func(chan Event, chan Event, map[string]interface{})

func GetRegistryInstance() map[string]func(chan Event, chan Event, map[string]interface{}) {
    if (registry == nil) {
        registry = make(map[string]func(chan Event, chan Event, map[string]interface{}))
        log.Info("Created...")
    }

    return registry
}
