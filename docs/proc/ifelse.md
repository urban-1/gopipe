# Proc: if/else

This component gives you more control over which processing steps are running
and when.

Example:
```
"proc": [
    ...
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
    ...
]
```

The above will hash the field `host` only if the `port_block` is set to true,
else it will append timestamp in the events' data. (I know it doesn't make sense
but just an example...)

NOTE: **In theory** if/else statements can be nested, however, this has not been
tested with more than 2 levels :) you have been warned
