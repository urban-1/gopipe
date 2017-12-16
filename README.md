# GoPipe

A pipe-line in go

## Why

- Learning Go
- Network oriented
- Allows the user to write more logic, still based on config

## Architecture

inputs/ouputs and procs

## Example Configs

### UDP FlowReplicator

    ./gopipe -c ./etc/flowreplicator.json

### LPM

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

-   Only one input: This can be fixed but there is no need atm

## Developers
