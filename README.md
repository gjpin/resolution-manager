# Description
A simple CLI to change display resolution and refresh rate. Based on [zergon321's work](https://gist.github.com/zergon321/4914a1af6c3573df47d959b064811f11).

# Usage
- resolution-manager WIDTH HEIGHT FREQUENCY

## Examples
- resolution-manager 3440 1440 144
- resolution-manager 2560 1440 144
- resolution-manager 1920 1080 60

# How to build
```
$env:GOOS="windows"
go build -o resolution-manager.exe
```