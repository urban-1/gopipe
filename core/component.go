package core

import (
)

// == Aliases the name to work for casting too?!?! Dont know dont ask
type Config = map[string]interface{}

type Component interface {
    Run()
	Stop()
}
