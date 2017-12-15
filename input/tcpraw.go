package input

import (
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering TCPRawInput")
    GetRegistryInstance()["TCPRawInput"] = NewTCPRawInput
}

type TCPRawInput struct {
    *TCPJSONInput
}

func NewTCPRawInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating TCPRawInput")

    // Defaults...
    m := TCPRawInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

    // Change to CSV
    m.Decoder = &RawLineCodec{}

    return &m
}
