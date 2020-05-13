# gglog

English | [简体中文](README_ZH.MD)

A lightweight log library, including local logs, kafka logs, Alibaba Cloud logs... 
Of course, plug-in access to the log library you want

## demo
you can see the demo in gglog_test.go

## feature
- Pluggable
- Concurrent security
- lightweight

## use step
### local log
- 1.NewGGLog()
- 2.Init()
- 3.FlushLog()    //or you can set the config param FlushInterval to control flush the log to file

### kafka log
- 1.NewGGLog()
- 2.Init()
- 3.Run()

### aliyun log
- 1.NewGGLog()
- 2.Init()
- 3.Run()

You can also use your own log library by implementing the Log interface

## log format
### log local
> [date time][code line][level]:msg

    [2020-03-11 17:19:22][testing.go:909][INFO]: info...
    [2020-03-11 17:31:42][testing.go:909][INFO]: info...
    [2020-03-11 17:53:59][testing.go:909][INFO]: info...
    [2020-03-11 18:01:35][testing.go:909][INFO]: info...

### kafka log
> [ip][code line][date time][level][desc][hostname]

    {"ip":"127.0.0.1", "location":"gglog.go:53", "tm":1589254203824, "level":"info", "desc":"I...", "hostname":"DESKTOP-7LEL6NV"}
    {"ip":"127.0.0.1", "location":"gglog.go:63", "tm":1589254203824, "level":"warn", "desc":"W...", "hostname":"DESKTOP-7LEL6NV"}
    {"ip":"127.0.0.1", "location":"gglog.go:73", "tm":1589254203824, "level":"error", "desc":"E...", "hostname":"DESKTOP-7LEL6NV"}

### aliyun log
> [ip][code line][date time][level][desc][hostname]

    {"ip":"127.0.0.1", "location":"gglog.go:53", "tm":1589254203824, "level":"info", "desc":"I...", "hostname":"DESKTOP-7LEL6NV"}
    {"ip":"127.0.0.1", "location":"gglog.go:63", "tm":1589254203824, "level":"warn", "desc":"W...", "hostname":"DESKTOP-7LEL6NV"}
    {"ip":"127.0.0.1", "location":"gglog.go:73", "tm":1589254203824, "level":"error", "desc":"E...", "hostname":"DESKTOP-7LEL6NV"}
