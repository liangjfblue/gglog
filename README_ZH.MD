# gglog

[![Build Status](https://travis-ci.org/liangjfblue/gglog.svg?branch=master)](https://travis-ci.org/github/liangjfblue/gglog)

[English](README.md) | 简体中文

😎轻量级的日志库，包括本地日志，kafka日志，阿里云日志... 当然，可以通过插件访问所需的日志库

## 例子
你可以在看demo gglog_test.go

## 特性
- 可插拔
- 并发安全
- 轻量级

## 使用步骤
### local log
- 1.NewGGLog()
- 2.Init()
- 3.FlushLog()    //或者您可以设置配置参数FlushInterval来控制将日志刷新到文件

### kafka log
- 1.NewGGLog()
- 2.Init()
- 3.Run()

### aliyun log
- 1.NewGGLog()
- 2.Init()
- 3.Run()

您还可以通过实现Log接口来使用自己的日志库

## 日志格式
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

