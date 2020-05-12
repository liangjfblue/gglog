package vlog

import (
	"time"
)

var (
	defaultOption = LogOptions{
		Name:                        "vlog",
		LogDir:                      "./logs",
		Level:                       1,
		OpenAccessLog:               true,
		OpenInterfaceAvgDurationLog: false,

		EnableLogHeader: true,
		EnableLogLink:   false,
		FlushInterval:   time.Duration(2) * time.Second,

		//kafka
		BrokerAddrs: []string{"127.0.0.1:9092"},
		Topic:       "gglog-topic",
		Partition:   1,
		Key:         "gglog-key",
		IsSync:      true,

		//aliyun
		Endpoint:     "127.0.0.1",
		Project:      "aliyun-log",
		LogStore:     "gglog-aliyun-log",
		AccessId:     "",
		AccessSecret: "",
		AliyunTopic:  "127.0.0.1",
		AliyunSource: "topic",
	}
)
