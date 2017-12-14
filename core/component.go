package core

import (
    "fmt"
    "time"
    log "github.com/sirupsen/logrus"
)

// == Aliases the name to work for casting too?!?! Dont know dont ask
type Config = map[string]interface{}

/**
 * Component's interface to reference any type
 */
type Component interface {
    Run()
	Stop()
}

/**
 * Parsing stats
 */
type ComponentStats struct {
    MsgCount uint64
    MsgCountOld uint64
    MsgRate uint64
    LastUpdate int64
}

func NewComponentStats() ComponentStats {
    return ComponentStats{0,0,0,0}
}


func (c *ComponentStats) DebugStr() string {
    return fmt.Sprintf("count=%d, rate=%d", c.MsgCount, c.MsgRate)

}

func (c *ComponentStats) AddMessage() {

    now := time.Now().Unix()

    if c.LastUpdate == 0 {
        c.LastUpdate = now
        c.MsgCount += 1
        c.MsgCountOld = 0
        return
    }

    c.MsgCount += 1

    // 5 second interval stats
    if now - c.LastUpdate > 3 {
        c.MsgRate = (c.MsgCount - c.MsgCountOld) / (uint64)(now - c.LastUpdate)
        c.MsgCountOld = c.MsgCount
        c.LastUpdate = now
    }
}

func (c *ComponentStats) Reset () {
    c.LastUpdate = 0
    c.MsgCount = 0
    c.MsgCountOld = 0
    c.MsgRate = 0
}

/**
 * ComponentBase has all core functions that EVERY component must have
 */
type ComponentBase struct {
    InQ chan *Event
    OutQ chan *Event
    Config Config
    MustStop bool
    Stats ComponentStats
}

func NewComponentBase(inQ chan *Event, outQ chan *Event, cfg Config) *ComponentBase {
    return &ComponentBase{inQ, outQ, cfg, false, NewComponentStats()}
}


func (p *ComponentBase) Stop() {
    p.MustStop = true
}

func  (p *ComponentBase) StatsAddMesg() {
    p.Stats.AddMessage()
}

func  (p *ComponentBase) PrintStats(name string, every uint64) {
    if p.Stats.MsgCount % every != 0 {
        return
    }

    inQLen := -1
    if p.InQ != nil {
        inQLen = len(p.InQ)
    }

    outQLen := -1
    if p.OutQ != nil {
        outQLen = len(p.OutQ)
    }

    log.Info(name, "> iq=", inQLen, ", oq=", outQLen, ", ", p.Stats.DebugStr())

}
