package glog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/liangjfblue/gglog/utils"

	"github.com/liangjfblue/gglog/vlog"

	"github.com/liangjfblue/gglog/vlog/glog/lib"
)

type glog struct {
	opts vlog.LogOptions

	GLog *lib.Logger
	once sync.Once
}

func New(opts ...vlog.LogOption) vlog.Log {
	options := vlog.NewOptions(opts...)

	return &glog{
		opts: options,
	}
}

func (s *glog) Name() string {
	return s.opts.Name
}

func (s *glog) Init(opts ...vlog.LogOption) {
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		s.initGLog()
	})
}

func (s *glog) Options() vlog.LogOptions {
	return s.opts
}

func (s *glog) String() string {
	return "vlog"
}

func (s *glog) Debug(format string, args ...interface{}) {
	if s.opts.Level <= utils.LevelD {
		s.GLog.DebugDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *glog) Info(format string, args ...interface{}) {
	if s.opts.Level <= utils.LevelI {
		s.GLog.InfoDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *glog) Warn(format string, args ...interface{}) {
	if s.opts.Level <= utils.LevelW {
		s.GLog.WarningDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *glog) Error(format string, args ...interface{}) {
	if s.opts.Level <= utils.LevelE {
		s.GLog.ErrorDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *glog) Access(format string, args ...interface{}) {
	if s.opts.OpenAccessLog == true {
		s.GLog.AccessDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *glog) InterfaceAvgDuration(format string, args ...interface{}) {
	if s.opts.OpenInterfaceAvgDurationLog == true {
		s.GLog.InterfaceAvgDurationDepth(1, fmt.Sprintf(format, args...))
	}
}

func (s *glog) FlushLog() {
	s.GLog.Flush()
}

func (s *glog) Run() {
	log.Println("nothing to do")
}

func (s *glog) Stop() {
	log.Println("nothing to do")
}

func (s *glog) initGLog() {
	if s.opts.LogDir == "" {
		panic("logDir is empty")
	}

	if _, err := os.Stat(s.opts.LogDir); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(s.opts.LogDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	s.GLog = lib.NewLogger().
		LogDir(s.opts.LogDir).
		EnableLogHeader(s.opts.EnableLogHeader).
		EnableLogLink(s.opts.EnableLogLink).
		FlushInterval(s.opts.FlushInterval).
		HeaderFormat(func(buf *bytes.Buffer, l lib.Severity, ts time.Time, pid int, file string, line int) {
			switch l {
			case lib.InfoLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][INFO]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case lib.DebugLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][DEBUG]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case lib.WarnLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][WARN]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case lib.ErrorLog:
				_, _ = fmt.Fprintf(buf, "[%s][%s:%d][ERROR]: ", ts.Format("2006-01-02 15:04:05"), file, line)
			case lib.AccessLog:
				_, _ = fmt.Fprintf(buf, "[%s]\t", ts.Format("2006-01-02 15:04:05"))
			case lib.InterfaceAvgDurationLog:
				_, _ = fmt.Fprintf(buf, "[%s]\t", ts.Format("2006-01-02 15:04:05"))
			}
		}).
		FileNameFormat(fileNameFormatFunc).
		Init()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		_, _ = os.Stdout.WriteString("flushing vlog\n")
		lib.Flush()
		_, _ = os.Stdout.WriteString("flush vlog done\n")
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
	case lib.SevDebug:
		tag = "debug"
	case lib.SevInfo:
		tag = "info"
	case lib.SevWarn:
		tag = "warning"
	case lib.SevError:
		tag = "error"
	case lib.SevAccess:
		tag = "access"
	case lib.SevFatal:
		tag = "fatal"
	case lib.SevInterfaceAvgDuration:
		tag = "iavgd"
	}

	return tag
}

func fileNameFormatFunc(severityLevel string, ts time.Time) string {
	return fmt.Sprintf(
		"%s.vlog.%04d-%02d-%02d",
		logTag(severityLevel), ts.Year(), ts.Month(), ts.Day())
}
