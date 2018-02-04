# Proc: In List

As the name suggest, checks if a field of the event exists in a list and stores
the result (true/false) into another field.

The list we check against can be provided by config (static) or can be regularly
read from a file. In both cases, the items in the list are strings and thus every
the data field is converted to string (using `%v`) in order to be checked against
the list.

The main function/purpose of this plugin is to verify against lists that change
regularly (ex IP blacklist) and thus the analysis has to take place at the
correct time (the time of the event) and cannot be performed in later time!

Example static configuration:

```
{
    "module": "InListProc",
    "in_field": "port",
    "out_field": "port_block",
    "reload_minutes": 100000000,
    "list": ["8080", "443", "23230", "14572", "17018"]
},
```

Example with configuration from file:

```
{
    "module": "InListProc",
    "in_field": "port",
    "out_field": "port_block",
    "reload_minutes": 2,
    "filepath": "/tmp/test.list"
},
```

NOTE: If both `filepath` and `list` are specified in the configuration, `list`
takes priority and overwrites the `filepath`

NOTE2: When filepath is used, to disable the auto-reload use
`"reload_minutes": 0` (or less). One can use `signal`s to reload this file.

## Signals

The supported signals are:

-   `reload`: This will attempt to reload the list from `filepath`. If filepath
    is not defined, it will print a warning and ignore the signal
