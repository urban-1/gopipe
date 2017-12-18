# Proc: Add Time

Every event gets a timestamp the moment it is created. This however lives in the
metadata and is not visible in the event's data (or systems output). This timestamp
can be added in the data under a key name with this compoment.

Example config:

```
{
    "module": "AddTimeProc",
    "field_name": "_timestamp"
}
```

Timestamp will be in `Data["_timestamp"]`
