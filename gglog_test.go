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
	"time"

	"github.com/liangjfblue/gglog/vlog/aliyun"

	"github.com/liangjfblue/gglog/vlog/kafka"

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

func TestKafkaLogSync(t *testing.T) {
	l := NewGGLog(
		Log(kafka.NewKafkaLog(
			vlog.BrokerAddrs([]string{"172.16.7.16:9092", "172.16.7.16:9093", "172.16.7.16:9094"}),
			vlog.Topic("kafka-vlog"),
			vlog.Key("kafka-key-test"),
			vlog.IsSync(true),
		)),
	)
	l.Init()
	defer l.Stop()
	go l.Run()

	start := time.Now()
	for i := 0; i < 200; i++ {
		l.Debug("D...")
		l.Info("I...")
		l.Warn("W...")
		l.Error("E...")
		l.Access("A...")
	}
	t.Log(time.Now().Sub(start))
}

func TestKafkaLogAsync(t *testing.T) {
	l := NewGGLog(
		Log(kafka.NewKafkaLog(
			vlog.BrokerAddrs([]string{"172.16.7.16:9092", "172.16.7.16:9093", "172.16.7.16:9094"}),
			vlog.Topic("kafka-vlog"),
			vlog.Key("kafka-key-test"),
			vlog.IsSync(false),
		)),
	)
	l.Init()
	defer l.Stop() //you can stop when you want to stop kafka log
	go l.Run()

	for i := 0; i < 100; i++ {
		l.Debug("D...")
		l.Info("I...")
		l.Warn("W...")
		l.Error("E...")
		l.Access("A...")
	}

	time.Sleep(5 * time.Second)
}

func TestKafkaLogAsyncRace(t *testing.T) {
	l := NewGGLog(
		Log(kafka.NewKafkaLog(
			vlog.BrokerAddrs([]string{"172.16.7.16:9092", "172.16.7.16:9093", "172.16.7.16:9094"}),
			vlog.Topic("kafka-vlog"),
			vlog.Key("kafka-key-test"),
			vlog.IsSync(false),
		)),
	)
	l.Init()
	defer l.Stop() //you can stop when you want to stop kafka log
	go l.Run()

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			l.Debug("D...")
			l.Info("I...")
			l.Warn("W...")
			l.Error("E...")
			l.Access("A...")
		})
	}

	time.Sleep(time.Second * 5)
}

func TestAliyunLog(t *testing.T) {
	l := NewGGLog(
		Log(aliyun.NewAliyunLog(
			vlog.Endpoint("cn-shenzhen.log.aliyuncs.com"),
			vlog.Project("gglog"),
			vlog.LogStore("test-gglog"),
			vlog.AccessId("xxx"),     //use your AccessId
			vlog.AccessSecret("xxx"), //use your AccessSecret
			vlog.AliyunTopic("gglog-aliyun"),
			vlog.AliyunSource("gglog-aliyun-test"),
		)),
	)
	l.Init()
	l.Run()
	defer l.Stop()

	for i := 0; i < 2; i++ {
		l.Debug("this is a test for debug")
		l.Info("i am a man")
		l.Warn("something is warning")
		l.Error("now have a wrong")
		l.Access("url:/v1/user/login method:POST")
	}
}
