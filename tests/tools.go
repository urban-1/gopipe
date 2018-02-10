package tests

import (
	"bytes"
	"encoding/json"
	"fmt"

	. "github.com/urban-1/gopipe/core"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func GetChannels() (chan *Event, chan *Event) {

	in := make(chan *Event, 1)
	out := make(chan *Event, 1)

	return in, out
}

func GetConfig(s string) Config {
	cfg := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &cfg)
	if err != nil {
		fmt.Printf("User error: cannot create mock config")
		panic("User error: cannot create mock config")
	}

	return cfg
}

func GetEvent(s string) *Event {
	e := map[string]interface{}{}

	d := json.NewDecoder(bytes.NewReader([]byte(s)))
	d.UseNumber()
	err := d.Decode(&e)

	if err != nil {
		fmt.Printf("User error: cannot create mock event")
		panic("User error: cannot create mock event")
	}

	return NewEvent(e)
}

func GetEventRun(s string, run bool) *Event {
	e := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &e)
	if err != nil {
		panic("GetEventRun: User error: cannot create mock event")
	}

	ret := NewEvent(e)
	ret.ShouldRun.Push(run)
	return ret
}

func GetRawEvent(s []byte) *Event {
	c := &RawLineCodec{}
	json_data, err := c.FromBytes(s)
	if err != nil {
		panic("GetRawEvent: User error: cannot create mock event")
	}

	return NewEvent(json_data)
}
