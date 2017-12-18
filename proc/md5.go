/*
    - MD5: Hash a set of given data fields and attach the results back to the
    event's data. Optionally, a salt can be provided
 */
package proc

import (
    "crypto/md5"
    "encoding/hex"
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering Md5Proc")
    GetRegistryInstance()["Md5Proc"] = NewMd5Proc
}

type Md5Proc struct {
    *ComponentBase
    InFields []string
    OutFields []string
    Salt string
}

func NewMd5Proc(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating Md5Proc")

    in_fields := []string{}
    if tmp, ok := cfg["in_fields"].([]interface{}); ok {
        in_fields = InterfaceToStringArray(tmp)
    }

    out_fields := []string{}
    if tmp, ok := cfg["out_fields"].([]interface{}); ok {
        out_fields = InterfaceToStringArray(tmp)
    }

    salt, ok := cfg["salt"].(string)
    if !ok {
        salt = ""
    }

    m := &Md5Proc{NewComponentBase(inQ, outQ, cfg), in_fields, out_fields, salt}
    m.Tag = "MD5-LOG"
    return m
}


func (p *Md5Proc) Run() {
    log.Debug("Md5Proc Starting ... ")
    p.MustStop = false

    for !p.MustStop {
        // Check if we should run (based on the events' if else state)
        e, err := p.ShouldRun()
        if err != nil {
            continue
        }

        for i, ifield := range p.InFields {
            b, ok := e.Data[ifield].(string)
            if !ok {
                log.Error("Failed to convert field ", ifield, " to string...")
                continue
            }

            md5tmp := md5.Sum([]byte(b+p.Salt))
            e.Data[p.OutFields[i]] = hex.EncodeToString(md5tmp[:])
        }

        p.OutQ<-e

        // Stats
        p.StatsAddMesg()
        p.PrintStats()

    }

    log.Info("Md5Proc Stopping!?")
}
