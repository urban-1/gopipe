package output

import (
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering FileCSVOutput")
    GetRegistryInstance()["FileCSVOutput"] = NewFileCSVOutput
}

type FileCSVOutput struct {
    *FileJSONOutput
}

func NewFileCSVOutput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating FileCSVOutput")

    headers := []string{}
    if tmp, ok := cfg["headers"].(InterfaceArray); ok {
        headers = tmp.ToStringArray()
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

    m := FileCSVOutput{NewFileJSONOutput(inQ, outQ, cfg).(*FileJSONOutput)}

    // Change to CSV
    m.Encoder = &CSVLineCodec{headers, sep, convert}

    return &m
}
