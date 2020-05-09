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
	}
)
