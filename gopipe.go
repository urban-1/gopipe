package main

import (
    "os"
    "encoding/json"
    "io/ioutil"
    "time"

    "github.com/urfave/cli"
    log "github.com/sirupsen/logrus"

    "gopipe/core"
    _ "gopipe/input"
    _ "gopipe/output"
    _ "gopipe/proc"
)

func init() {
    customFormatter := new(log.TextFormatter)
    customFormatter.FullTimestamp = true
    log.SetFormatter(customFormatter)
}


func instanceFromConfig(cfg core.Config, ch1 chan *core.Event, ch2 chan *core.Event, reg core.Registry) (core.Component, error) {

    module_name, ok := cfg["module"].(string)
    if !ok {
        log.Error("Missing 'module' (module name) from configuration")
        return nil, cli.NewExitError("Missing 'module' (module name) from configuration", -3)
    }

    log.Info("Loading ", module_name)

    mod_constructor, ok := reg[module_name]
    if !ok {
        log.Error("Unknown module '", module_name, "'")
        return nil, cli.NewExitError("Unknown module '" + module_name + "'", -4)
    }

    log.Info("Loaded!")

    return mod_constructor(ch1, ch2, cfg), nil
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
        cli.BoolFlag{
            Name:  "debug, d",
            Usage: "Enable debug mode",
        },
    }


    app.Action = func(c *cli.Context) error {
        if (c.String("config") == "") {
            const msg = "You need to provide config file..."
            log.Error(msg)
            return cli.NewExitError(msg, -1)
        }

        if c.Bool("debug") {
            log.SetLevel(log.DebugLevel)
        }

        DN , _ := os.Getwd()
        log.Info("Running from directory '", DN, "'")
        log.Info("Loding configuration from '", c.String("config"), "'")
        raw, err := ioutil.ReadFile(c.String("config"))
        if err != nil {
            log.Error(err.Error())
            return cli.NewExitError(err.Error(), -2)
        }

        var CFG core.Config
        if err := json.Unmarshal(raw, &CFG); err != nil {
            log.Error(err.Error())
            return cli.NewExitError(err.Error(), -2)
        }


        tmp_q_len, ok := CFG["main"].(core.Config)["channel_size"].(float64)
        if !ok {
            log.Warn("No buffer length (channel_size) set in main section. Assuming 1")
            tmp_q_len = 1
        } else {
            log.Info("Using buffer length: ", tmp_q_len)
        }

        Q_LEN := int(tmp_q_len)


        // Load registry
        reg := core.GetRegistryInstance()

        // In module
        in, ok := CFG["in"].(core.Config)
        if !ok {
            log.Error("You need to define 'in' section in your config")
            return cli.NewExitError("You need to define 'in' section in your config", -2)
        }

        // Store our modules
        var mods []core.Component
        // Store all channels
        var chans []chan *core.Event

        // Create the first Q (ouput)
        tmpch := make(chan *core.Event, Q_LEN)
        chans = append(chans, tmpch)

        // Append in module
        tmp, err := instanceFromConfig(in, nil, tmpch, reg)
        if err != nil {
            return err
        }
        mods = append(mods, tmp)

        proc, ok := CFG["proc"].([]interface{})
        for index, cfg := range proc {
            // Create a new one for every output

            chans = append(chans, make(chan *core.Event, Q_LEN))

            log.Info("Loading processor module ", index)
            tmp, err = instanceFromConfig(
                cfg.(core.Config),
                chans[len(chans)-2],
                chans[len(chans)-1],
                reg)

            if err != nil {
                return err
            }
            mods = append(mods, tmp)
        }

        // Output module
        out, ok := CFG["out"].(core.Config)
        if !ok {
            log.Error("You need to define 'out' section in your config")
            return cli.NewExitError("You need to define 'out' section out your config", -2)
        }
        tmp, err = instanceFromConfig(out, chans[len(chans)-1], nil, reg)
        if err != nil {
            return err
        }
        mods = append(mods, tmp)
        log.Info("Created ", len(chans), " channels")

        // Start all (reverse order)
        for _, mod := range mods {
            go mod.Run()
        }

        // Now loop forever
        for {
            time.Sleep(time.Duration(1000)*time.Millisecond)
            // log.Warn("All Channels Empty")
        }

        return nil
    }

    app.Run(os.Args)


}