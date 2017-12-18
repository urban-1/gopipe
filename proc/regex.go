/*
    - REGEX: Given a regex with named captures, convert each event from a text
    one to a data one (using the "message" field, which is where Str codecs store
    their output)
 */
package proc

import (
    "regexp"
    . "gopipe/core"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.Info("Registering RegexProc")
    GetRegistryInstance()["RegexProc"] = NewRegexProc
}

type RegexProc struct {
    *ComponentBase
    Regs []*regexp.Regexp
}

func NewRegexProc(inQ chan *Event, outQ chan *Event, cfg Config) Component {
    log.Info("Creating RegexProc")

    regs := []*regexp.Regexp{}
    tmpres := cfg["regexes"].([]interface{})

    for _, v := range tmpres {
        regs = append(regs, regexp.MustCompile(v.(string)))
    }
    m := &RegexProc{NewComponentBase(inQ, outQ, cfg), regs}
    m.Tag = "REGEX-PROC"
    return m
}


func (p *RegexProc) Run() {
    log.Debug("RegexProc Starting ... ")
    p.MustStop = false
    for !p.MustStop {

        e, err := p.ShouldRun()
        if err != nil {
            continue
        }

        allok := false

        for renum, re := range p.Regs {

            log.Debug("Testing renum=", renum)
            match  := re.FindStringSubmatch(e.Data["message"].(string))
            if match == nil {
                continue
            }

            log.Debug("Matched! renum=", renum)

            // Mark processed
            allok = true

            for i, name := range re.SubexpNames() {

                if i != 0 {
                    e.Data[name] = match[i]
                }
            }
            break
        }

        if !allok {
            log.Warn("Skipping non-mathching line: ", e.Data["message"].(string))
            continue
        }

        p.OutQ<-e

        // Stats
        p.StatsAddMesg()
        p.PrintStats()

    }

    log.Info("RegexProc Stopping!?")
}
