# Output: UDP

Similar with the Output UDP component, each packet is processed as a separate
message. (might change)

There are different ways for this module to interpret messages, depending on
which Codec is used:

# `UDPJSONOutput`

This is the default and will try to decode every line into a JSON object.
Example config:

     {
         "module": "UDPJSONOutput",
         "target": "127.0.0.1",
         "port": 9092
     }

# `UDPCSVOutput`

Reads each line and interprets it as CSV. Extra parameters:

    {
        "module": "UDPCSVOutput",
        "target": "127.0.0.1",
        "port": 9092,
        "headers": ["hello", "test", "src"],
        "separator": ","
    }

Separator is optional...

# `UDPStrOutput`

Encodes `Data["message"]` to bytes and uses it as payload

# `UDPRawOutput`

Uses `Data["raw"]` as payload - without modifying it
