# Output: File

The file output component is responsible for opening, writing and rotating files.
It uses `LineCodec`s to serialize the events. At the moment it supports CSV and
JSON formats. However this can easily be extended to more formats.

## `FileJSONOutput`

Example config:

```
"out": {
    "module": "FileJSONOutput",
    "rotate_seconds": 60,
    "folder": "/tmp",
    "file_name_format": "gopipe-20060102-150405.json"
}
```
Where:

-   `folder`: Is the output directory
-   `file_name_format`: defines the naming pattern of each log file. This will
    be parsed with `<time>.Format(<file_name_format>)` to form the final filename
-   `rotate_seconds`: Every how many seconds you want to rotate the file

## `FileCSVOutput`

Similar to the JSON output, however in CSV we need to define the column names
since the data are written in ordered mode. This compoment will not output any
headers, example:

```
"out": {
    "module": "FileCSVOutput",
    "rotate_seconds": 60,
    "folder": "/tmp",
    "file_name_format": "gopipe-20060102-150405.csv",
    "headers": ["_timestamp", "host", "port", "host_hash", "payload"],
    "separator": ","
}
```

The separator is optional and the default character is (`,`/`comma`)

NOTE: **This Components can be used as processing compoments too,** they are aware
of the output channels and they can pass through events. This is useful for
logging and sampling purposes.
