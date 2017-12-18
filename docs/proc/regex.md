# Proc: Regex

Given a regex with named captures, convert each event from a text one to a data
one (using the "message" field, which is where Str codecs store their output)

Example:

```
{
    "module": "RegexProc",
    "regexes": [
        "(?mi)(?P<host>[0-9a-z]+) (?P<port>[0-9]+): (?P<hostEvent>.*)"
    ]
}
```

The above will convert the log-line:

```
hostname27667 31881: Message-14214
```

To:

```
{"host":"hostname27667","hostEvent":"Message-14214","port":"31881"}
```

Note that the `port` is still a string but it can be converted to `int` with the
use of `CastProc`.


Full regex syntax available at: https://github.com/google/re2/wiki/Syntax
