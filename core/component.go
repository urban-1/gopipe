package core

import (
    "fmt"
)

// == Aliases the name to work for casting too?!?! Dont know dont ask
type Config = map[string]interface{}

type Component struct {
    config Config
    mustStop bool
    inQ chan Event
    outQ chan Event
}

func (c *Component) Run() {
    fmt.Print("Not implemented")
}

func (c *Component) Stop() {
    c.mustStop = true
}

func (c *Component) CleanUp() {}
