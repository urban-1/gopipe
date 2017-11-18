package main

import (
    //"fmt"
    log "github.com/sirupsen/logrus"
    "gopipe/core"
    _ "gopipe/input"
)

func init() {
    customFormatter := new(log.TextFormatter)
    customFormatter.FullTimestamp = true
    log.SetFormatter(customFormatter)
}

func main() {

    e := core.NewDataEvent()
    reg := core.GetRegistryInstance()
    log.Info(len(reg))
    log.Info(e.Type())
    reg["TCPInput"](nil, nil, nil)
    //fmt.Printf("hello, world\n")
}
