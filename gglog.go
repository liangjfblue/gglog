package gglog

type GGLog interface {
	// The service name
	Name() string
	// Init initialises options
	Init(...Option)
	// Options returns the current options
	Options() Options
	// The service implementation
	String() string

	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Access(format string, args ...interface{})
	InterfaceAvgDuration(format string, args ...interface{})
	FlushLog()
}

type Option func(*Options)

func NewGGLog(opts ...Option) GGLog {
	return newGGLog(opts...)
}
