/*
    The main `gopipe` binary.

    Reads config, creates modules and channels and starts the whole pipeline...

    Given thet all modules are reusable, one can use the rest of the packages as
    library of components and write their own "main". This could accept different
    config formats, customize the logger, etc
 */
package main

import (
    "encoding/json"
    "io/ioutil"
    "os/signal"
    "net/http"
    "os/exec"
    "syscall"
    "time"
    "fmt"
    "os"

    "github.com/urfave/cli"
    log "github.com/sirupsen/logrus"

    "gopipe/core"
    _ "gopipe/input"
    _ "gopipe/output"
    _ "gopipe/proc"
)

// Store our modules
var mods []core.Component

func init() {
    customFormatter := new(log.TextFormatter)
    customFormatter.FullTimestamp = true
    log.SetFormatter(customFormatter)
}

// Given the configuration of a module, the channels and the registry, create and
// return an instance
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

// Loop for ever while sleeping for interval_seconds in every iteration
// and execut a command (discarding output)
func runTask(name string, parts []string, interval_seconds uint64, signals []interface{}) {
    cmd, args := parts[0], parts[1:]
    for {
    	if err := exec.Command(cmd, args...).Run(); err != nil {
            log.Error("Failed to run command '"+name+"': " + err.Error())
            return
    	}
    	log.Debug("Command '"+name+"' run successfully...")

        // Signal other components
        for _, signal := range signals {
            mod := int(signal.(core.Config)["mod"].(float64))
            sig := signal.(core.Config)["signal"].(string)
            log.Infof("Invoking signal '%s' on component %d", sig, mod)
            mods[mod].Signal(sig)
        }
        time.Sleep(time.Duration(interval_seconds)*time.Second)
    }
}

func apiStatus(w http.ResponseWriter, r *http.Request) {
    var err error

    log.Info("ACCESS ", r.URL.Path)
    ret := map[string]interface{}{
        "components": []interface{}{},
    }

    for _, mod := range mods {
        ret["components"] = append(ret["components"].([]interface{}), mod.GetStatsJSON())
    }

    var content []byte

    if content, err = json.Marshal(ret); err != nil {
        log.Error("Failed to parse component stats to JSON")
    }
    fmt.Fprintf(w, "%s", string(content))
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
        log.Info("Loading configuration from '", c.String("config"), "'")
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


        // Channel buffer size
        tmp_f64, ok := CFG["main"].(core.Config)["channel_size"].(float64)
        if !ok {
            log.Warn("No buffer length (channel_size) set in main section. Assuming 1")
            tmp_f64 = 1
        } else {
            log.Info("Using buffer length: ", tmp_f64)
        }

        Q_LEN := int(tmp_f64)

        // Printing frequency
        tmp_f64, ok = CFG["main"].(core.Config)["stats_every"].(float64)
        if ok {
            core.STATS_EVERY = uint64(tmp_f64)
        }

        // Load registry
        reg := core.GetRegistryInstance()

        // In module
        in, ok := CFG["in"].(core.Config)
        if !ok {
            log.Error("You need to define 'in' section in your config")
            return cli.NewExitError("You need to define 'in' section in your config", -2)
        }

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

        // Start the HTTP server
        tmpport, ok :=  CFG["main"].(core.Config)["apiport"].(float64)
        var apiport string
        if !ok {
            apiport = "9090"
        } else {
            apiport = fmt.Sprintf("%v",tmpport)
        }
        http.HandleFunc("/status", apiStatus) // set router

        go func() error {
            err = http.ListenAndServe(":"+apiport, nil) // set listen port
            if err != nil {
                log.Error("Failed to bind TCP port for API server")
                log.Fatal("ListenAndServe: ", err)
                return err
            }
	    return nil
        }()


        // Spawn all tasks before starting processing
        tasks, ok := CFG["tasks"].([]interface{})
        for _, task := range tasks {
            go runTask(
                task.(core.Config)["name"].(string),
                core.InterfaceToStringArray(task.(core.Config)["command"].([]interface{})),
                uint64(task.(core.Config)["interval_seconds"].(float64)),
                task.(core.Config)["signals"].([]interface {}))
        }

        // Start all
        for _, mod := range mods {
            go mod.Run()
        }

        chExit := make(chan os.Signal, 1)
        chInst := make(chan os.Signal, 1)
    	signal.Notify(chExit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
        signal.Notify(chInst, syscall.SIGUSR1, syscall.SIGUSR1)

        // Now loop forever
        run := true
        for run {
            select {
            case <-chExit:
                log.Info("gopipe stoping components...")

                // Kill the input
                mods[0].Stop()

                log.Info("gopipe waiting for queues to empty...")
                for i, c := range chans {
                    for len(c) > 0 {
                        log.Info("Waiting on channel", i, " len=", len(c))
                        time.Sleep(time.Duration(1000)*time.Millisecond)
                    }
                }

                for _, mod := range mods {
                    mod.Stop()
                }

        		log.Info("gopipe exiting...")
                run = false
                break
            case sig := <-chInst:
                switch sig {
                case syscall.SIGUSR1:
                    for _, mod := range mods {
                        mod.MustPrintStats()
                    }
                case syscall.SIGUSR2:
                    //handle SIGTERM
                }
            default:
                time.Sleep(time.Duration(1000)*time.Millisecond)
                // log.Warn("All Channels Empty")
        	}

        }


        return nil
    }

    app.Run(os.Args)


}
