package gglog

import "time"

var (
	defaultOption = Options{
		Name:                        "gglog",
		LogDir:                      "./logs",
		Level:                       1,
		OpenAccessLog:               true,
		OpenInterfaceAvgDurationLog: false,

		EnableLogHeader: true,
		EnableLogLink:   false,
		FlushInterval:   time.Duration(2) * time.Second,
	}
)
