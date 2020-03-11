# gglog
a useful log, depend on glog

## demo
you can see the demo in test/glog_test.go

## use step

- 1.NewGGLog()
- 2.Init()
- 3.FlushLog()    //or you can set the config param FlushInterval to control flush the log to file

## log format
> [date time][code line][level]:msg

    [2020-03-11 17:19:22][testing.go:909][INFO]: info...
    [2020-03-11 17:31:42][testing.go:909][INFO]: info...
    [2020-03-11 17:53:59][testing.go:909][INFO]: info...
    [2020-03-11 18:01:35][testing.go:909][INFO]: info...

