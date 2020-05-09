package gglog

import (
	"sync"
)

type GGLog interface {
	Name() string
	Init()
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Access(format string, args ...interface{})
	InterfaceAvgDuration(format string, args ...interface{})
	FlushLog()
}

type Option func(*Options)

type gglog struct {
	opts Options

	once sync.Once
}

func NewGGLog(opts ...Option) GGLog {
	options := newOptions(opts...)

	return &gglog{
		opts: options,
	}
}

func (l *gglog) Name() string {
	return l.opts.Name
}

func (l *gglog) Init() {
	l.opts.Log.Init()
}

func (l *gglog) Debug(format string, args ...interface{}) {
	for _, fn := range l.opts.Before {
		fn()
	}
	l.opts.Log.Debug(format, args...)
	for _, fn := range l.opts.After {
		fn()
	}
}

func (l *gglog) Info(format string, args ...interface{}) {
	for _, fn := range l.opts.Before {
		fn()
	}
	l.opts.Log.Info(format, args...)
	for _, fn := range l.opts.After {
		fn()
	}
}

func (l *gglog) Warn(format string, args ...interface{}) {
	for _, fn := range l.opts.Before {
		fn()
	}
	l.opts.Log.Warn(format, args...)
	for _, fn := range l.opts.After {
		fn()
	}
}

func (l *gglog) Error(format string, args ...interface{}) {
	for _, fn := range l.opts.Before {
		fn()
	}
	l.opts.Log.Error(format, args...)
	for _, fn := range l.opts.After {
		fn()
	}
}

func (l *gglog) Access(format string, args ...interface{}) {
	for _, fn := range l.opts.Before {
		fn()
	}
	l.opts.Log.Access(format, args...)
	for _, fn := range l.opts.After {
		fn()
	}
}

func (l *gglog) InterfaceAvgDuration(format string, args ...interface{}) {
	for _, fn := range l.opts.Before {
		fn()
	}
	l.opts.Log.InterfaceAvgDuration(format, args...)
	for _, fn := range l.opts.After {
		fn()
	}
}

func (l *gglog) FlushLog() {
	for _, fn := range l.opts.Before {
		fn()
	}
	l.opts.Log.FlushLog()
	for _, fn := range l.opts.After {
		fn()
	}
}
