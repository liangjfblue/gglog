/**
 *
 * @author liangjf
 * @create on 2020/5/9
 * @version 1.0
 */
package vlog

type Log interface {
	Name() string
	Init(...LogOption)
	Run()
	Stop()
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Access(format string, args ...interface{})
	InterfaceAvgDuration(format string, args ...interface{})
	FlushLog()
}

type LogOption func(*LogOptions)

type ICallBack interface {
	Success(string)
	Fail(error)
}
