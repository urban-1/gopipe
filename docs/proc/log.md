# Proc: Log

Use the logger (logrus) to print the event. This component does not
modify the events in any way. **This can also be used as an OUTPUT component!**

Its main purpose is to be used for debugging and allow us to inspect events
mid-flight.

Example config:

```
{
    "module": "LogProc",
    "level": "info"
}
```

Supported log levels are:

-   `debug`
-   `info`
-   `warn`
