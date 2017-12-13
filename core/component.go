package core

import (
    // log "github.com/sirupsen/logrus"
)

// == Aliases the name to work for casting too?!?! Dont know dont ask
type Config = map[string]interface{}

/**
 * Component's interface to reference any type
 */
type Component interface {
    Run()
	Stop()
}

/**
 * ComponentBase has all core functions that EVERY component must have
 */
type ComponentBase struct {
    InQ chan Event
    OutQ chan Event
    Config Config
    MustStop bool
}

func NewComponentBase(inQ chan Event, outQ chan Event, cfg Config) *ComponentBase {
    return &ComponentBase{inQ, outQ, cfg, false}
}


func (p *ComponentBase) Stop() {
    p.MustStop = true
}
