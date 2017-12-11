package main

import (
    "os"
    "encoding/json"
    "io/ioutil"

    "github.com/urfave/cli"
    log "github.com/sirupsen/logrus"

    "gopipe/core"
    _ "gopipe/input"
)

func init() {
    customFormatter := new(log.TextFormatter)
    customFormatter.FullTimestamp = true
    log.SetFormatter(customFormatter)
}



func main() {
    app := cli.NewApp()
    // app.Name = "gopipe"
    app.Usage = "Pipeline processing in Go!"
    app.Version = "0.0.1"

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "config, c",
            Usage: "Load configuration from `FILE` (required)",
        },
    }

    app.Action = func(c *cli.Context) error {
        if (c.String("config") == "") {
            const msg = "You need to provide config file..."
            log.Error(msg)
            return cli.NewExitError(msg, -1)
        }

        DN , _ := os.Getwd()
        log.Info("Running from directory '", DN, "'")
        log.Info("Loding configuration from '", c.String("config"), "'")
        raw, err := ioutil.ReadFile(c.String("config"))
        if err != nil {
            log.Error(err.Error())
            cli.NewExitError(err.Error(), -2)
        }

        var CFG core.Config
        json.Unmarshal(raw, &CFG)
        log.Info(CFG["in"].(core.Config)["module"])

        e := core.NewDataEvent()
        reg := core.GetRegistryInstance()
        log.Info(len(reg))
        log.Info(e.Type())
        reg["TCPInput"](nil, nil, nil)
        return nil
    }

    app.Run(os.Args)


}
