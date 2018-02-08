package proc

import (
    "testing"
    . "github.com/urban-1/gopipe/tests"
)


func TestAddFieldValue(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"value": "blah",
			"field_name": "test"
		}
	`))
    go comp.Run()

    e := <-out
    if e.Data["test"] != "blah" {
        t.Error("AddField did not add anything!!")
    }
}

func TestAddFieldExpression(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"expression": "doesnt",
			"field_name": "test"
		}
	`))
    go comp.Run()

    e := <-out
    if e.Data["test"] != "matter" {
        t.Error("AddField: I was expecting test: matter")
		t.Error(e.Data)
    }
}

func TestAddFieldExpression2(t *testing.T) {
    in,  out := GetChannels()
    in <- GetEvent(`{"doesnt": "matter"}`)

    comp := NewAddFieldProc(in, out, GetConfig(`
		{
			"expression": "10*10+3",
			"field_name": "test"
		}
	`))
    go comp.Run()

    e := <-out
    if int(e.Data["test"].(float64)) != 103 {
        t.Error("AddField: I was expecting test:103...")
		t.Error(e.Data)
    }
}
