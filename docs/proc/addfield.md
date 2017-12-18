# Proc: Add Field

Add a new field in the data of the event. In its simplest form the configuration
of this plugin looks like:

```
{
    "module": "AddFieldProc",
    "field_name": "new_field",
    "value": 1000
}
```

Although adding a static value to all the events is useful (for example add the
processing server hostname), it does not have that many applications... It is
beneficial to be able to perform some basic operations/expressions, for example:

```
{
    "module": "AddFieldProc",
    "field_name": "10port",
    "expression": "port * 10"
}
```
In the above example, `port` is a field of our event's data. Any field can be used
in the expression. The expression validator we have used is
https://github.com/Knetic/govaluate and thus we support any expression it does -
check its README for details.
