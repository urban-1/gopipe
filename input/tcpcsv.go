package input

import (
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering TCPCSVInput")
    GetRegistryInstance()["TCPCSVInput"] = NewTCPCSVInput
}

type TCPCSVInput struct {
    *TCPJSONInput
}

func NewTCPCSVInput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating TCPCSVInput")

    headers := []string{}
    if tmp, ok := cfg["headers"].([]interface{}); ok {
        headers = Interface2StringArray(tmp)
    }
    log.Infof("  Headers %v", headers)

    sep := ","[0]
    if tmp, ok := cfg["separator"].(string); ok {
        sep = tmp[0]
    }

    convert := true
    if tmp, ok := cfg["convert"].(bool); ok {
        convert = tmp
    }

    // Defaults...
    m := TCPCSVInput{NewTCPJSONInput(inQ, outQ, cfg).(*TCPJSONInput)}

    // Change to CSV
    m.Decoder = &CSVLineCodec{headers, sep, convert}

    return &m
}
