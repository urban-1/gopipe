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

    return &FileJSONOutput{NewComponentBase(inQ, outQ, cfg),
        0, folder, pattern, rotate_seconds, nil,
        &JSONLineCodec{}}
}

func (p *FileJSONOutput) Stop() {
    p.MustStop = true
}

/**
 * Check and rotate the output file if needed
 *
 */
func (p *FileJSONOutput) checkRotate() {
    now := time.Now().Unix()
    if int(now - p.LastRotate) >= p.RotateSeconds {
        p.getNewFile()
    }
}

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
        }
        data = append(data, byte('\n'))
        _, err = p.Fd.Write(data)

        if err != nil {
            log.Error("Failed to write data: ", err.Error())
        }

        // Stats
        p.StatsAddMesg()
        p.PrintStats("JSONFile", 50000)
    }
    log.Debug("FileJSONOutput Stopping")
}
