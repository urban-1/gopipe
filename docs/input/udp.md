# Input: UDP

Listen on a UDP socket for messages. Each packet is processed as a separate
message. (might change)

There are different ways for this module to interpret messages, depending on
which Codec is used:

# `UDPJSONInput`

 This is the default and will try to decode every line into a JSON object.
 Example config:

    {
        "module": "UDPJSONInput",
        "listen": "0.0.0.0",
        "port": 9092
    }

# `UDPCSVInput`

Reads each line and interprets it as CSV. Extra parameters:

    {
        "module": "UDPCSVInput",
        "listen": "0.0.0.0",
        "port": 9092
        "headers": ["hello", "test", "src"],
        "separator": ",",
        "convert": false
    }

If `convert` is set to true, this compoment will try to convert ints and floats.
If set to false, these will be represented as strings. `CastProc` compoment can
be used later to change types


# `UDPStrInput`

Reads each line, converts it to string and stores it in `Data["message"]`

# `UDPRawInput`

Reads each line as byte array and store it in `Data["raw"]`
