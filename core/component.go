package core

import (
    "fmt"
)

type Component struct {
    config map[string]interface{}
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
