# TODO

-   Embedded web server for status reporting (json? Prometheus?)
-   Port to https://github.com/Jeffail/gabs
-   Check if we need buffered writers or OS will do the job
-   Split lines option in Str readers? (mainly to support multiple messages in a
    single UDP packet...)
-   Conditionals based on https://github.com/Knetic/govaluate/blob/master/README.md
-   Find out a more "automatic" way to defin codec and its parameter from config
    such that extending with new codecs does not require new component structs.
    Maybe via a codec Registry?
-   ADD TASKS! (background jobs not touching events)
