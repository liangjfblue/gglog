package gglog

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/liangjfblue/gglog/glog"
)

type gglog struct {
	opts Options

	GLog *glog.Logger
	once sync.Once
}

func newGGLog(opts ...Option) GGLog {
	options := newOptions(opts...)

	return &gglog{
		opts: options,
	}
}

func (s *gglog) Name() string {
	return s.opts.Name
}

func (s *gglog) Init(opts ...Option) {
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		s.initGLog()
	})
}

func (s *gglog) Options() Options {
	return s.opts
}

func (s *gglog) String() string {
	return "gglog"
}

const (
	levelD = iota + 1
	levelI
	levelW
	levelE
)

func (s *gglog) Debug(format string, args ...interface{}) {
	if s.opts.Level <= levelD {
		s.GLog.DebugDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *gglog) Info(format string, args ...interface{}) {
	if s.opts.Level <= levelI {
		s.GLog.InfoDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *gglog) Warn(format string, args ...interface{}) {
	if s.opts.Level <= levelW {
		s.GLog.WarningDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *gglog) Error(format string, args ...interface{}) {
	if s.opts.Level <= levelE {
		s.GLog.ErrorDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *gglog) Access(format string, args ...interface{}) {
	if s.opts.OpenAccessLog == true {
		s.GLog.AccessDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *gglog) InterfaceAvgDuration(format string, args ...interface{}) {
	if s.opts.OpenInterfaceAvgDurationLog == true {
		s.GLog.InterfaceAvgDurationDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *gglog) FlushLog() {
	s.GLog.Flush()
}

func (s *gglog) initGLog() {
	if s.opts.LogDir == "" {
		panic("logDir is empty")
	}

	if _, err := os.Stat(s.opts.LogDir); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(s.opts.LogDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	s.GLog = glog.NewLogger().
		LogDir(s.opts.LogDir).
		EnableLogHeader(s.opts.EnableLogHeader).
		EnableLogLink(s.opts.EnableLogLink).
		FlushInterval(s.opts.FlushInterval).
		HeaderFormat(func(buf *bytes.Buffer, l glog.Severity, ts time.Time, pid int, file string, line int) {
			switch l {
			case glog.InfoLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][INFO]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.DebugLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][DEBUG]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.WarnLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][WARN]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.ErrorLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][ERROR]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case glog.AccessLog:
				_, _ = fmt.Fprintf(buf, "[%s]\t", ts.Format("2006-01-02 15:04:05"))
			case glog.InterfaceAvgDurationLog:
				_, _ = fmt.Fprintf(buf, "[%s]\t", ts.Format("2006-01-02 15:04:05"))
			}
		}).
		FileNameFormat(fileNameFormatFunc).
		Init()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		_, _ = os.Stdout.WriteString("flushing log\n")
		glog.Flush()
		_, _ = os.Stdout.WriteString("flush log done\n")
		signal.Reset(os.Interrupt)

		proc, err := os.FindProcess(syscall.Getpid())
		if err != nil {
			panic(err)
		}

		err = proc.Signal(os.Interrupt)
		if err != nil {
			panic(err)
		}
	}()
}

func logTag(severityLevel string) string {
	tag := "info"
	switch severityLevel {
	case glog.SevDebug:
		tag = "debug"
	case glog.SevInfo:
		tag = "info"
	case glog.SevWarn:
		tag = "warning"
	case glog.SevError:
		tag = "error"
	case glog.SevAccess:
		tag = "access"
	case glog.SevFatal:
		tag = "fatal"
	case glog.SevInterfaceAvgDuration:
		tag = "iavgd"
	}

	return tag
}

func fileNameFormatFunc(severityLevel string, ts time.Time) string {
	return fmt.Sprintf(
		"%s.log.%04d-%02d-%02d",
		logTag(severityLevel), ts.Year(), ts.Month(), ts.Day())
}
