package input

import (
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering TCPStrInput")
    GetRegistryInstance()["TCPStrInput"] = NewTCPStrInput
}

type TCPStrInput struct {
    *TCPJSONInput
}

func NewTCPStrInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating TCPStrInput")

    // Defaults...
    m := TCPStrInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

    // Change to CSV
    m.Decoder = &StringLineCodec{}

    return &m
}
