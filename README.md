# GoPipe
<img src="https://travis-ci.org/urban-1/gopipe.svg?branch=master"/>

A processing pipe-line (similar to logstash) written in Go

## Quick-Start

There is a Makefile you can use to quickly build and start working on this
project:

```
$ make help

Available targets:

  rdkafka         Build librdkafka locally
  gopipe          Build gopipe
  tests           Run tests (includes fmt and vet)
  show_coverage   Show coverage results
```



## Architecture

Our goal is to define a pipeline where events are received and processed in a
modular way. Events can contain:

-   Raw unformated bytes (`Data["raw"]`)
-   String contents (`Data["message"]`)
-   Decoded JSON value (`Data` as `map[string]interface{}`)

We call our modules "Components" and these fall into three categories:

-   Input: Components that generate events and push them into a Go channel
-   Proc: Processing components that modify event data. They read and write events
    from and to Go channels
-   Output: there are just reading events and output in some form. They usually
    do not push them back into a Go channel, but instead discard them.

In essence, all components are the "same" (implement the same interface). The
only difference is which channels are made available to them.

## Whys and Hows?

1.  Well... I have done something similar in C++ for processing netflow packets
and thought since Go is (really) fast and concurrent is a perfect match for such
an application.

2. Why not use something already out there: We could extend an existing framework
however, this is a Go learning exercise to replicate the C++ code...

3. How is that different? We focus on a systems perspective and we want this
framework to be more network/data oriented rather than `log`s oriented:
    -   Raw data handling - see the `flowreplicator.json` configuration.
    -   Allows the user to write more logic, still based on config (if/else support)
    -   (Experimental) Support tasks to feed and update processing module's data

4. What are the future plans: We plan to maintain and extend this
until we fully port our C++ code... Maintenance will continue but we kinda
hope we will extend as needed with the help of the community.


## Components

### Inputs

-   **[TCP](docs/input/tcp.md)**: Supporting raw, string, CSV and JSON
-   **[UDP](docs/input/udp.md)**: Supporting raw, string, CSV and JSON
-   **[Kafka](docs/input/kafka.md)**: Supporting raw, string, CSV and JSON

### Processing

-   **[Add field](docs/proc/addfield.md)**: Add a new field based on static value
    or expression
-   **[Add time](docs/proc/addtime.md)**: Adds timestamp to the data
-   **[Cast](docs/proc/cast.md)**: Converts fields to different data types
-   **[Drop field](docs/proc/dropfield.md)**: Removes fields from data
-   **[If/Else](docs/proc/ifelse.md)**: Control flow with if/else/endif
-   **[In List](docs/proc/inlist.md)**: Checks a field against a list of values
-   **[Log](docs/proc/log.md)**: Logs the events' data to stdout
-   **[Longest Prefix Match](docs/proc/lpm.md)**: Performs LPM and attaches meta-data to the events' data
-   **[MD5](docs/proc/md5.md)**: Hash event's fields
-   **[Regex](docs/proc/regex.md)**: Convert string events into data ones
-   **[Sampler](docs/proc/sampler.md)**: Selectively forward events (one every X)

### Output

-   **[File](docs/output/file.md)**: Supporting CSV and JSON
-   **[Null](docs/output/null.md)**: Blackholes events
-   **[UDP](docs/output/udp.md)**: Supporting raw and string


## Example Configs

### UDP FlowReplicator

Will replicate and optionally sample UDP packtes:

```
{
    "main": {
        "num_cpus": 2,
        "log_level": 1,
        "channel_size": 50000,
        "stats_every": 100000
    },
    "in": {
        "module": "UDPRawInput",
        "listen": "0.0.0.0",
        "port": 9090
    },
    "proc": [
        {
            "module": "SamplerProc",
            "every": 2
        },
        {
            "module": "UDPRawOutput",
            "target": "127.0.0.1",
            "port": 9091
        }
    ],
    "out": {
        "module": "UDPRawOutput",
        "target": "127.0.0.1",
        "port": 9092
    }
}
```

### Format log lines + if/else demo

Receives lines over TCP, parses them into data fields, adds timestamp, converts
some data to different data-types, discards the original message, hashes some
fields, etc

```
{
    "main": {
        "num_cpus": 2,
        "log_level": 1,
        "channel_size": 50000,
        "stats_every": 100000
    },
    "in": {
        "module": "TCPStrInput",
        "listen": "0.0.0.0",
        "port": 9092,
        "headers": ["hello", "test", "src"],
        "separator": ",",
        "convert": false
    },
    "proc": [
        {
            "module": "RegexProc",
            "regexes": [
                "(?mi)(?P<host>[0-9a-z]+) (?P<port>[0-9]+): (?P<hostEvent>.*)"
            ]
        },
        {
            "module": "DropFieldProc",
            "field_name": "message"
        },
        {
            "module": "CastProc",
            "fields": ["port"],
            "types": ["int"]
        },
        {
            "module": "InListProc",
            "in_field": "port",
            "out_field": "port_block",
            "reload_minutes": 100000000,
            "list": ["8080", "443", "23230", "14572", "17018"]
        },
        {"module": "if",  "condition": "port_block == true "},
            {
                "module": "Md5Proc",
                "in_fields": ["host"],
                "out_fields": ["host_hash"],
                "salt": "andPepper!"
            },
        {"module": "else"},
            {
                "module": "AddTimeProc",
                "field_name": "_timestamp"
            },
        {"module": "endif"},
        {
            "module": "LogProc",
            "level": "info"
        }
    ],
    "out": {
        "module": "FileJSONOutput",
        "rotate_seconds": 60,
        "folder": "/tmp",
        "file_name_format": "gopipe-20060102-150405.json"
    }
}
```

### Longest Prefix Match

Receive on a TCP socket listening for JSON line:
```
{
    "main": {
        "num_cpus": 2,
        "log_level": 1,
        "channel_size": 50000,
        "stats_every": 10000000
    },
    "in": {
        "module": "TCPJSONInput",
        "listen": "0.0.0.0",
        "port": 9092
    },
    "proc": [
        {
            "module": "AddTimeProc",
            "field_name": "_timestamp"
        },
        {
            "module": "LPMProc",
            "filepath": "/tmp/prefix-asn.txt",
            "reload_minutes": 1440,
            "in_fields": ["src", "dst"],
            "out_fields": [
                {"newkey": "sky_{{in_field}}_prefix", "metakey": "prefix"},
                {"newkey": "sky_{{in_field}}_asn", "metakey": "asn"}
            ]
        }
    ],
    "out": {
        "module": "FileJSONOutput",
        "rotate_seconds": 60,
        "folder": "/tmp",
        "file_name_format": "gopipe-20060102-150405.json"
    }
}
```

### Tasks

The following config part defines a task that runs every 10 seconds. Usually you
would like to update file sources for `InListProc` and `LPMProc` components...
In such cases the idea is that you have a small shell-script somewhere in your
system that will update your local files. Then you need to "invoke" a reload to
load the new data in memory:

```
...
    ],
    "out": {
        "module": "FileJSONOutput",
        "rotate_seconds": 60,
        "folder": "/tmp",
        "file_name_format": "gopipe-20060102-150405.json"
    },
    "tasks": [
        {
            "name": "LSing...",
            "command": ["ls", "-al"],
            "interval_seconds": 10,
            "signals": [
                {"mod": 4, "signal": "reload"}
            ]
        },
        {
            "name": "LSing...2",
            "command": ["ls", "-al"],
            "interval_seconds": 10,
            "signals": []
        }
    ]
...
```

Above we define two tasks. The difference between them is that the first one
will signal a component if it runs successfully. The signal `"reload"` is going
to be sent to component `4` and is up to the component to handle it.

The component index is defined as the order of this component in config
**including input components**. Given that at the moment we only support one
input, component `4` above is the 3rd in `proc` section.

## Limitations

-   Only one input is supported at the moment but this might change
-   A bit immature framework :) we need more components
-   JSON: Decoding with `UseNumber()` is needed for correct output, however,
    it breaks `govaluate` so when comparing you have to use `json_to_float64()`.
    See `TestIf` for example...

## Developers

Hello! The idea is that we can provide JSON-configurable pipeline processing
capability for the end user. However, we do need more components for various
jobs and maybe codecs!

-   Components should be extremely easy to implement. Use `proc/log.go` as a
    starting point (~60 LOC) to implement you component.

-   Codecs: Have a quick look into `linecodecs.go`. One can easily implement new
    line encoders/decoders. These can then be plugged into input/output modules

Not sure with what to help? have a look at [TODO.md](TODO.md) As always,
comments, suggestions, documentation, bug reports, etc are more than
welcome :)
