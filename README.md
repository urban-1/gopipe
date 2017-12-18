# GoPipe

A processing pipe-line (similar to logstash) written in Go

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
    -   FUTURE: Allows the user to write more logic, still based on config (if/else support)
    -   FUTURE: Support tasks to feed and update processing module's data

4. What are the future plans: We plan to maintain and extend this
until we fully port our C++ code... Maintenance will continue but we kinda
hope we will extend as needed with the help of the community.


## Components

### Inputs

-   **[TCP](docs/input/TCP.md)**: Supporting raw, string, CSV and JSON
-   **UDP**: Supporting raw, string, CSV and JSON

### Processing

-   **Add time**: Adds timestamp to the data
-   **Cast**: Converts fields to different data types
-   **Drop field**: Removes fields from data
-   **In List**: Checks a field against a list of values
-   **Log**: Logs the events' data to stdout
-   **Longest Prefix Match**: Performs LPM and attaches meta-data to the events' data

### Output

-   **File**: Supporting CSV and JSON
-   **Null**: Blackholes events
-   **UDP**: Supporting raw and string


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

### Format log lines

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

### LPM

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

Data are expected in the following format:

    prefix/len json-meta-data

example:

    160.202.15.0/24 {"asn": 1}
    176.52.166.0/24 {"asn": 1}
    176.52.167.0/24 {"asn": 1}
    198.22.130.0/24 {"asn": 1}
    199.246.102.0/24 {"asn": 1}
    200.52.157.0/24 {"asn": 1}
    202.134.183.0/24 {"asn": 1}
    202.63.238.0/24 {"asn": 1}
    207.227.224.0/22 {"asn": 1}
    207.227.228.0/22 {"asn": 1}

Example script to load and format them:

    curl http://lg01.infra.ring.nlnog.net/table.txt | awk -F' ' '{print $1,"{\"asn\": "$2"}"}' > ~/tmp/prefix-asn.txt



# Limitations

-   Only one input is supported at the moment but this might change
-   A bit immature framework :) we need more components

## Developers

Hello! The idea is that we can provide JSON-configurable pipeline processing
capability for the end user. However, we do need more components for various
jobs and maybe codecs!

-   Components should be extremely easy to implement. Use `proc/log.go` as a
    starting point (~60 LOC) to implement you component.

-   Codecs: Have a quick look into `linecodecs.go`. One can easily implement new
    line encoders/decoders. These can then be plugged into input/output modules

As always, comments, suggestions, documentation, bug reports, etc are more than
welcome :)

[Also, you might have detected that this code has no tests! We are new to Go and
we are still figuring this out...]
