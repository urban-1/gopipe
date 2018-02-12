package input

import (
	"bytes"
	"encoding/json"
	"github.com/urban-1/gopipe/core"
	"github.com/urban-1/gopipe/output"
	. "github.com/urban-1/gopipe/tests"
	"testing"
	"time"
)

func TestUDPJSON(t *testing.T) {
	in, out := GetChannels()
	mid := make(chan *core.Event, 1)

	out <- GetEvent(`{"a": 1}`)

	cin := NewUDPJSONInput(nil, in, GetConfig(`
		{"listen": "127.0.0.1", "port": 10000}
	`))

	cout := output.NewUDPJSONOutput(out, mid, GetConfig(`
		{"target": "127.0.0.1", "port": 10000}
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

func TestUDPCSV(t *testing.T) {
	in, out := GetChannels()
	mid := make(chan *core.Event, 1)

	out <- GetEvent(`{"a": 1}`)

	cin := NewUDPCSVInput(nil, in, GetConfig(`
		{"listen": "127.0.0.1", "port": 10001, "headers": ["a"]}
	`))

	cout := output.NewUDPCSVOutput(out, mid, GetConfig(`
		{"target": "127.0.0.1", "port": 10001, "headers": ["a"]}
	`))

	go cin.Run()
	time.Sleep(time.Duration(1) * time.Second)
	go cout.Run()

	// Test UDP as middle stage
	e := <-mid

	// Here we are dealing with json.... (from GetEvent)
	tmp, _ := e.Data["a"].(json.Number).Int64()
	if tmp != 1 {
		t.Error("UDP MID error: I was expecting a: 1")
		t.Error(e.Data)
	}

	// Test socket
	e = <-in
	// Here we are dealing with int64, since CSVLineCodec.Convert is set
	if e.Data["a"].(int64) != 1 {
		t.Error("UDP IO error: I was expecting a: 1")
		t.Error(e.Data)
	}
}

func TestUDPRaw(t *testing.T) {
	in, out := GetChannels()
	mid := make(chan *core.Event, 1)

	bin := []byte("1234567890")
	out <- GetRawEvent(bin)

	cin := NewUDPRawInput(nil, in, GetConfig(`
		{"listen": "127.0.0.1", "port": 10002}
	`))

	cout := output.NewUDPRawOutput(out, mid, GetConfig(`
		{"target": "127.0.0.1", "port": 10002}
	`))

	go cin.Run()
	time.Sleep(time.Duration(1) * time.Second)
	go cout.Run()

	// Test UDP as middle stage
	e := <-mid

	// Here we are dealing with json.... (from GetEvent)
	if bytes.Compare(e.Data["bytes"].([]byte), bin) != 0 {
		t.Error("UDP MID error: I was expecting a: 1")
		t.Error(e.Data)
	}

	// Test socket
	e = <-in
	// Here we are dealing with int64, since CSVLineCodec.Convert is set
	if bytes.Compare(e.Data["bytes"].([]byte), bin) != 0 {
		t.Error("UDP IO error: I was expecting a: 1")
		t.Error(e.Data)
	}
}

func TestUDPStr(t *testing.T) {
	in, out := GetChannels()
	mid := make(chan *core.Event, 1)

	msg := "[32241.135047] cfg80211:  DFS Master region: unset"
	out <- GetEvent(`{"message": "` + msg + `"}`)

	cin := NewUDPStrInput(nil, in, GetConfig(`
		{"listen": "127.0.0.1", "port": 10003}
	`))

	cout := output.NewUDPStrOutput(out, mid, GetConfig(`
		{"target": "127.0.0.1", "port": 10003}
	`))

	go cin.Run()
	time.Sleep(time.Duration(1) * time.Second)
	go cout.Run()

	// Test UDP as middle stage
	e := <-mid

	// Here we are dealing with json.... (from GetEvent)
	if e.Data["message"].(string) != msg {
		t.Error("UDP MID error: I was expecting a: 1")
		t.Error(e.Data)
	}

	// Test socket
	e = <-in
	// Here we are dealing with int64, since CSVLineCodec.Convert is set
	if e.Data["message"].(string) != msg {
		t.Error("UDP IO error: I was expecting a: 1")
		t.Error(e.Data)
	}
}
