package test

import (
	"testing"

	"github.com/liangjfblue/gglog"
)

func TestGLog(t *testing.T) {
	ggLog := gglog.NewGGLog(
		gglog.Name("test-gglog"),
		gglog.LogDir("./test-logs"),
		gglog.Level(1),
		gglog.OpenInterfaceAvgDurationLog(true),
	)

	ggLog.Init()

	ggLog.Debug("debug...")
	ggLog.Info("info...")
	ggLog.Warn("info...")
	ggLog.Error("info...")
	ggLog.Access("access...")

	ggLog.FlushLog()
}
