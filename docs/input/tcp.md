# Input: TCP

Listen on a TCP socket for messages. Each line is processed as a separate
message. Maximum line length is 65000 bytes.

There are different ways for this module to interpret messages, depending on
which Codec is used:

# `TCPJSONInput`

 This is the default and will try to decode every line into a JSON object.
 Example config:

    {
        "module": "TCPStrInput",
        "listen": "0.0.0.0",
        "port": 9092
    }

# `TCPCSVInput`

Reads each line and interprets it as CSV. Extra parameters:

    {
        "headers": ["hello", "test", "src"],
        "separator": ",",
        "convert": false
    }

If `convert` is set to true, this compoment will try to convert ints and floats.
If set to false, these will be represented as strings. `CastProc` compoment can
be used later to change types


# `TCPStrInput`

Reads each line, converts it to string and stores it in `Data["message"]`

# `TCPRawInput`

Reads each line as byte array and store it in `Data["raw"]`
