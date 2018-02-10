package input

import (
	. "github.com/urban-1/gopipe/tests"
	"github.com/urban-1/gopipe/output"
	"github.com/urban-1/gopipe/core"
	"testing"
	"time"
	"encoding/json"
)

func TestUDP(t *testing.T) {
	in, out := GetChannels()
	mid := make(chan *core.Event, 1)

	out <- GetEvent(`{"a": 1}`)

	cin := NewUDPJSONInput(nil, in, GetConfig(`
		{"listen": "127.0.0.1", "port": 10000}
	`))

	cout := output.NewUDPJSONOutput(out, mid, GetConfig(`
		{"target": "localhost", "port": 10000}
	`))

	go cin.Run()
	time.Sleep(time.Duration(1) * time.Second)
	go cout.Run()

	// Test UDP as middle stage
	e := <-mid
	tmp, _ := e.Data["a"].(json.Number).Int64()
	if tmp != 1 {
		t.Error("UDP MID error: I was expecting a: 1")
		t.Error(e.Data)
	}

	// Test socket
	e = <-in
	tmp, _ = e.Data["a"].(json.Number).Int64()
	if tmp != 1 {
		t.Error("UDP IO error: I was expecting a: 1")
		t.Error(e.Data)
	}
}
