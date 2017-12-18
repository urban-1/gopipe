# Proc: MD5

Hash a set of given data fields and attach the results back to the event's data.
Optionally, a salt can be provided.

Example:

```
{
    "module": "Md5Proc",
    "in_fields": ["host"],
    "out_fields": ["host_hash"],
    "salt": "andPepper!"
}
```

The result will be in `Data["host_hash"]`! The lists `in_fields` and `out_fields`
should have the same length.
