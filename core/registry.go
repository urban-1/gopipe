//
// - Component global registry: This is a map in the format `componentName =>
// Constructor`. The name is the class name which is also used in the configuration
// "module" field
//
package core

import (
	log "github.com/sirupsen/logrus"
)

// Dont even try to pronounce it...
type Registry = map[string]func(chan *Event, chan *Event, map[string]interface{}) Component

// Create singleton registry
var registry Registry

// Singleton implementation that returns the Global registry
func GetRegistryInstance() Registry {
	if (registry == nil) {
		registry = make(Registry)
		log.Info("Created...")
	}

	return registry
}
