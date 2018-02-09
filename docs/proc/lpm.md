# Proc: Longest Prefix Match (LPM)

Longest Prefix Match: Loads a list of prefixes (network addresses)
periodically from a file and performs LPM for event data fields. The result
can have metadata which are then exported back to the event's data. For
example the following will figure out the matched prefix and Autonomous
System number (ASN) and attach them to the event's data:

```
{
    "module": "LPMProc",
    "filepath": "/tmp/prefix-asn.txt",
    "reload_minutes": 1440,
    "in_fields": ["src", "dst"],
    "out_fields": [
        {"newkey": "_{{in_field}}_prefix", "metakey": "prefix"},
        {"newkey": "_{{in_field}}_asn", "metakey": "asn"}
    ]
}
```

Data are expected in the following format:

    prefix/len json-meta-data

Example data:

    10.0.0.0/8 {"asn": -1, "owner": "me", "other": "metadata"}
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

Output (for `src == "10.1.1.1"`):

```
{..."_src_asn": -1, "_src_prefix": "10.0.0.0/8" ...}
```

Example script to load and format them:

```
curl http://lg01.xxx.xxx/table.txt | awk -F' ' '{print $1,"{\"asn\": "$2"}"}' > ~/tmp/prefix-asn.txt
```

NOTE: **For concistency, the output fields will be populated ("") even if lookup fails**

NOTE2:To disable the auto-reload use `"reload_minutes": 0` (or less). One can
use `signal`s to reload this file.

## Signals

The supported signals are:

-   `reload`: This will attempt to reload the tree from `filepath`.
