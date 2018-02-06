# Input: Kafka

Listen on a TCP socket for messages. Each line is processed as a separate
message. Maximum line length is 65000 bytes.

There are different ways for this module to interpret messages, depending on
which Codec is used:

# `KafkaJSONInput`

This is the default and will try to decode every line into a JSON object.
Example config:

    "in": {
        "module": "KafkaJSONInput",
        "topics": ["kafka-test"],
        "group": "bigflow",
        "brokers": "localhost:9092",
        "topic_conf": {
             "enable.auto.commit": "true",
             "auto.commit.interval.ms": "1000",
             "auto.offset.reset": "smallest"
        }

    },

The `topic_conf` will be used for the topic default settings


# `KafkaCSVInput`, `KafkaStrInput`, `KafkaRawInput`

As with the TCP/UDP equivalents, some extraparameters are needed in the config
block of this component. However, **non of the above have been tested**
