/**
 *
 * @author liangjf
 * @create on 2020/5/9
 * @version 1.0
 */
package gglog

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/liangjfblue/gglog/vlog"

	"github.com/liangjfblue/gglog/vlog/glog"
)

func TestGGlogString(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	l := NewGGLog(
		Log(glog.New(
			vlog.Level(1),
			vlog.Name("test-vlog"),
			vlog.LogDir("./test-logs"),
		)),
		Context(ctx),
		Before(func() {
			fmt.Println("before")
		}),
		After(func() {
			fmt.Println("after")
		}),
	)
	l.Init()

	l.Debug("Debug...")
	l.Info("Info...")
	l.Warn("Warn...")
	l.Error("Error...")
	l.Access("Access...")

	//flush log right now
	l.FlushLog()
}

func TestGGlogOther(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	l := NewGGLog(
		Log(glog.New(
			vlog.Level(1),
			vlog.Name("test-vlog"),
			vlog.LogDir("./test-vlogs"),
		)),
		Context(ctx),
		Before(func() {
			fmt.Println("before")
		}),
		After(func() {
			fmt.Println("after")
		}),
	)
	l.Init()

	l.Debug("Debug...: %s", nil)
	l.Debug("Debug...: %s", []byte{})
	l.Debug("Debug...: %s", errors.New("test error"))

	//flush log right now
	l.FlushLog()
}
