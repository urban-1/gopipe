package proc

import (
	log "github.com/sirupsen/logrus"
	. "github.com/urban-1/gopipe/tests"
	"testing"
)

func TestRegexProc(t *testing.T) {
	in, out := GetChannels()
	in <- GetEvent(`{"message": "up02.somewhere.com 8080: All clean"}`)

	comp := NewRegexProc(in, out, GetConfig(`{
		"regexes": [
			"(?mi)(?P<host>[.0-9a-z]+) (?P<port>[0-9]+): (?P<hostEvent>.*)"
		]
	}`))
	go comp.Run()

	e := <-out
	if e.Data["host"] != "up02.somewhere.com" {
		t.Error("Regex did not match!")
		log.Error(e.Data)
	}

	if _, ok := e.Data["message"]; !ok {
		t.Error("Regex Proc did not output!")
		log.Error(e.Data)
	}

}

func TestRegexShouldNotRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"message": "up02.somewhere.com 8080: All clean"}`, false)

	comp := NewRegexProc(in, out, GetConfig(`{
		"regexes": [
			"(?mi)(?P<host>[.0-9a-z]+) (?P<port>[0-9]+): (?P<hostEvent>.*)"
		]
	}`))
	go comp.Run()

	e := <-out
	if _, ok := e.Data["host"]; ok {
		t.Error("Regex run when it shouldn't")
		log.Error(e.Data)
	}
}

func TestRegexShouldRun(t *testing.T) {
	in, out := GetChannels()
	in <- GetEventRun(`{"message": "up02.somewhere.com 8080: All clean"}`, true)

	comp := NewRegexProc(in, out, GetConfig(`{
		"regexes": [
			"(?mi)(?P<host>[.0-9a-z]+) (?P<port>[0-9]+): (?P<hostEvent>.*)"
		]
	}`))
	go comp.Run()

	e := <-out
	if _, ok := e.Data["host"]; !ok {
		t.Error("Regex did not run when it should")
		log.Error(e.Data)
	}
}
