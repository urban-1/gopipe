# Proc: Cast

**FULL CREDIT:** https://github.com/tsaikd/gogstash/blob/master/filter/typeconv/filtertypeconv.go

Convert fields to different data types. Supported target types are:

-   string, str
-   int
-   float

Example config:

```
{
    "module": "CastProc",
    "fields": ["port", "rate"],
    "types": ["int", "float"]
}
```
