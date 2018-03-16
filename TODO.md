# TODO

## General

-   Figure out TravisCI builds...
-   Consider porting JSON to https://github.com/Jeffail/gabs
-   Check if we need buffered writers or OS will do the job
-   Split lines option in Str readers? (mainly to support multiple messages in a
    single UDP packet...)
-   Find out a more "automatic" way to define codec and its parameter from config
    such that extending with new codecs does not require new component structs.
    Maybe via a codec Registry?
-   Complete tests aiming for 85%+
-   Stress and memleak test
-   Create a prometheus metrics end-point
-   External component loading on runtime (given a folder path) if possible so
    custom modules can be easily created and used
-   Allow multiple input components
-   Allow for YAML config

## Component Ideas

-   Add OS command input component (periodically running something). Note that
    this can be implemented atm with a `task` that pipes to `netcat` if you are
    using TCP/UDP input.
-   SQL output component, maybe via https://github.com/volatiletech/sqlboiler
-   ElasticSearch output component
-   InfluxDB maybe?
-   Kafka output component
