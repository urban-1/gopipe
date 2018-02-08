package proc

import (
    "testing"
	"bytes"
	"os"
    . "github.com/urban-1/gopipe/tests"
	log "github.com/sirupsen/logrus"
)


func TestLogProcOutput(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "blah"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }


}

func TestLogProcInfo(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "info"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}


func TestLogProcDebug(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "debug"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}


func TestLogProcWarn(t *testing.T) {
    in,  out := GetChannels()
	in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewLogProc(in, out, GetConfig(`{"level": "warn"}`))
    go comp.Run()

    e := <-out
    if e.Data["doesnt"] != "matter" {
        t.Error("Log Proc did not output!")
    }
}


func TestLogShouldNotRun(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEventRun(`{"doesnt": "matter"}`, false)

	var buf bytes.Buffer
    log.SetOutput(&buf)
    defer func() {
        log.SetOutput(os.Stderr)
    }()

	comp := NewLogProc(in, out, GetConfig(`{"level": "debug"}`))
    go comp.Run()

    <-out
	if bytes.Contains(buf.Bytes(), []byte(`LogProc: {`)) {
        t.Error("Log Proc run... it shouldn't")
		t.Error(buf.String())
	}
}

func TestLogShouldRun(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEventRun(`{"doesnt": "matter"}`, true)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	comp := NewLogProc(in, out, GetConfig(`{"level": "warn"}`))
    go comp.Run()

    <-out
	if !bytes.Contains(buf.Bytes(), []byte(`LogProc: {`)) {
        t.Error("Log Proc run... it shouldn't")
		t.Error(buf.String())
	}
}
