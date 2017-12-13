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
    TCPInput
}

func NewTCPCSVInput(inQ chan Event, outQ chan Event, cfg Config) Component {
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

    return &TCPCSVInput{
            TCPInput{
                *NewComponentBase(inQ, outQ, cfg),
                &CSVLineCodec{
                    headers,
                    sep},
                cfg["listen"].(string), uint32(cfg["port"].(float64))}}
}
