package core

import (
    log "github.com/sirupsen/logrus"
)

type Registry = map[string]func(chan *Event, chan *Event, map[string]interface{}) Component

// Create singleton registry
var registry Registry

func GetRegistryInstance() Registry {
    if (registry == nil) {
        registry = make(Registry)
        log.Info("Created...")
    }

    return registry
}
