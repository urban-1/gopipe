# Proc: Drop field

This compoment drops a field from the events' data. This is particularly useful
when you want to get rid of the original string message after it is parsed. Therefore,
you would most commonly use it after a `RegexProc`. Additionaly, it can be used
to clean up temporary fields used for expressions and calculations.

Example config:

```
{
    "module": "DropFieldProc",
    "field_name": "message"
},
```
