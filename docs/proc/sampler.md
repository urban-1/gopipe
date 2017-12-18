# Proc: Sampler

Keep one event every `X` ecents and discart the rest. Admittedly not the most
useful component... but helpful when simulating packet drops or when sampling
the output of the pipeline (quality control? maybe...)

Example:

```
{
    "module": "SamplerProc",
    "every": 2
}
```
