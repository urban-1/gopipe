/*
    This package contains all output plugins.

    Output plugins are mainly used in the "out" section of the config, however,
    they can also be used in any processing step to split the flow output of
    the framework. The output compoment should support that which means that it
    should be aware of the `outQ` and check if is nil ("out" section) or not
    (processing section):

        // Check if we are being used in proc!
        if p.OutQ != nil {
            p.OutQ<-e
        }

    - File: Output to timestamped files with regular (time-based) rotation.
 */
package output

import (
    "os"
    "time"
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering FileJSONOutput")
    GetRegistryInstance()["FileJSONOutput"] = NewFileJSONOutput

    log.Info("Registering FileCSVOutput")
    GetRegistryInstance()["FileCSVOutput"] = NewFileCSVOutput
}

type FileJSONOutput struct {
    *ComponentBase
    LastRotate int64
    Folder string
    Pattern string
    RotateSeconds int
    Fd *os.File
    Encoder LineCodec
}

func NewFileJSONOutput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating FileJSONOutput")

    folder := "/tmp"
    if tmp, ok := cfg["folder"].(string); ok {
        folder = tmp
    }

    pattern := "gopipe-20060102-150405.unknown"
    if tmp, ok := cfg["file_name_format"].(string); ok {
        pattern = tmp
    }

    rotate_seconds := 60
    if tmp, ok := cfg["rotate_seconds"].(float64); ok {
        rotate_seconds = int(tmp)
    }

    m := &FileJSONOutput{NewComponentBase(inQ, outQ, cfg),
        0, folder, pattern, rotate_seconds, nil,
        &JSONLineCodec{}}

    m.Tag = "OUT-FILE-JSON"

    return m
}


// Check and rotate the output file if needed
func (p *FileJSONOutput) checkRotate() {
    now := time.Now().Unix()
    if int(now - p.LastRotate) >= p.RotateSeconds {
        p.getNewFile()
    }
}

// Create a new file, close the old file if required
func (p *FileJSONOutput) getNewFile() {

    if p.Fd != nil {
        log.Debug("Closing old file")
        p.Fd.Sync()
        p.Fd.Close()
    }

    now := time.Now()
    fname := now.Format(p.Pattern)
    fname = p.Folder + "/" + fname

    log.Info("Creating ", fname)

    tmp, err := os.Create(fname)
    if err != nil {
        panic("Failed to open output file - Check permissions of " + p.Folder)
    }
    p.Fd = tmp
    p.LastRotate = now.Unix()

}


func (p *FileJSONOutput) Run() {
    p.MustStop = false
    log.Debug("FileJSONOutput Starting ... ")
    p.getNewFile()

    var data []byte
    var err error

    for !p.MustStop {
        p.checkRotate()

        log.Debug("FileJSONOutput Reading")
        e := <- p.InQ

        data, err = p.Encoder.ToBytes(e.Data)
        if err != nil {
            log.Error("Failed to encode data: ", err.Error())
            continue
        }

        _, err = p.Fd.Write(data)

        if err != nil {
            log.Error("Failed to write data: ", err.Error())
        }

        // Check if we are being used in proc!
        if p.OutQ != nil {
            p.OutQ<-e
        }

        // Stats
        p.StatsAddMesg()
        p.PrintStats()
    }
    log.Debug("FileJSONOutput Stopping")
}


// File CSV implementation
type FileCSVOutput struct {
    *FileJSONOutput
}

func NewFileCSVOutput(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating FileCSVOutput")

    headers := []string{}
    if tmp, ok := cfg["headers"].([]interface{}); ok {
        headers = InterfaceToStringArray(tmp)
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

    m.Tag = "OUT-FILE-CSV"

    // Change to CSV
    m.Encoder = &CSVLineCodec{headers, sep, convert}

    return &m
}
