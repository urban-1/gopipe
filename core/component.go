// Core package contains all common structs and functions.
//
// - Config: An alias to `map[string]interface{}`
//
// - Component interface and component base struct
//
package core

import (
    "fmt"
    "time"
    log "github.com/sirupsen/logrus"
)

// Component Global Config: How often to print stats in the log. This can be
// changes in configuration via `stats_every`
var STATS_EVERY uint64 = 100000

// Config alias
type Config = map[string]interface{}

// Component's interface for abstraction
type Component interface {
    Run()
	Stop()
    PrintStats()
}

// Each component's processing stats
type ComponentStats struct {
    MsgCount uint64
    MsgCountOld uint64
    MsgRate uint64
    LastUpdate int64
}

func NewComponentStats() ComponentStats {
    return ComponentStats{0,0,0,0}
}

// Return a string with rate and count (for logging purposes)
func (c *ComponentStats) DebugStr() string {
    return fmt.Sprintf("rate=%-7d count=%d", c.MsgRate, c.MsgCount)

}

// Increments the MsgCount and if required, calculates the MsgRate. the default
// interval is 3 seconds TODO: Global config
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

// Reset the stats back to 0
func (c *ComponentStats) Reset () {
    c.LastUpdate = 0
    c.MsgCount = 0
    c.MsgCountOld = 0
    c.MsgRate = 0
}

// ComponentBase implements core methods that EVERY component must have (avoid
// code duplication)
type ComponentBase struct {
    InQ chan *Event
    OutQ chan *Event
    Config Config
    MustStop bool
    Stats ComponentStats
    Tag string
}

// Create a new component given an input channel, an output channel and the
// component's config
func NewComponentBase(inQ chan *Event, outQ chan *Event, cfg Config) *ComponentBase {
    m := &ComponentBase{inQ, outQ, cfg, false, NewComponentStats(), "Base"}
    return m
}

// By default, just set MustStop to false. Component implementations should be
// taking this into consideration in their Run() methods
func (p *ComponentBase) Stop() {
    p.MustStop = true
}

// Wrapper around p.Stats
func  (p *ComponentBase) StatsAddMesg() {
    p.Stats.AddMessage()
}

// Logs the stats if needed. It will log every STATS_EVERY (default 50K messages)
// and can be disabled with stats_every = 0
func  (p *ComponentBase) PrintStats() {
    if STATS_EVERY == 0 {
        return
    }

    if p.Stats.MsgCount % STATS_EVERY != 0 {
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

    log.Infof("%15s> iq=%-5d oq=%-5d %s", p.Tag, inQLen, outQLen, p.Stats.DebugStr())

}
