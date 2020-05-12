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
	}
)
