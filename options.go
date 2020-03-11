package gglog

import (
	"time"
)

type Options struct {
	Name                        string
	LogDir                      string
	Level                       int32
	OpenAccessLog               bool
	OpenInterfaceAvgDurationLog bool

	EnableLogHeader bool
	EnableLogLink   bool
	FlushInterval   time.Duration
}

func newOptions(opts ...Option) Options {
	opt := defaultOption

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func LogDir(logDir string) Option {
	return func(o *Options) {
		o.LogDir = logDir
	}
}

func Level(level int32) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func OpenAccessLog(openAccessLog bool) Option {
	return func(o *Options) {
		o.OpenAccessLog = openAccessLog
	}
}

func OpenInterfaceAvgDurationLog(openInterfaceAvgDurationLog bool) Option {
	return func(o *Options) {
		o.OpenInterfaceAvgDurationLog = openInterfaceAvgDurationLog
	}
}

func EnableLogHeader(enableLogHeader bool) Option {
	return func(o *Options) {
		o.EnableLogHeader = enableLogHeader
	}
}

func EnableLogLink(enableLogLink bool) Option {
	return func(o *Options) {
		o.EnableLogLink = enableLogLink
	}
}

func FlushInterval(flushInterval time.Duration) Option {
	return func(o *Options) {
		o.FlushInterval = flushInterval
	}
}
