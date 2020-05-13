package vlog

import (
	"time"
)

type LogOptions struct {
	Name                        string
	LogDir                      string
	Level                       int32
	OpenAccessLog               bool
	OpenInterfaceAvgDurationLog bool

	//glog
	EnableLogHeader bool
	EnableLogLink   bool
	FlushInterval   time.Duration

	//kafka
	BrokerAddrs []string
	Partition   int
	Topic       string
	Key         string
	IsSync      bool
	CallBack    ICallBack

	//aliyun
	Endpoint     string
	Project      string
	LogStore     string
	AccessId     string
	AccessSecret string
	AliyunTopic  string
	AliyunSource string
}

func NewOptions(opts ...LogOption) LogOptions {
	opt := defaultOption

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Name(name string) LogOption {
	return func(o *LogOptions) {
		o.Name = name
	}
}

func LogDir(logDir string) LogOption {
	return func(o *LogOptions) {
		o.LogDir = logDir
	}
}

func Level(level int32) LogOption {
	return func(o *LogOptions) {
		o.Level = level
	}
}

func OpenAccessLog(openAccessLog bool) LogOption {
	return func(o *LogOptions) {
		o.OpenAccessLog = openAccessLog
	}
}

func OpenInterfaceAvgDurationLog(openInterfaceAvgDurationLog bool) LogOption {
	return func(o *LogOptions) {
		o.OpenInterfaceAvgDurationLog = openInterfaceAvgDurationLog
	}
}

func EnableLogHeader(enableLogHeader bool) LogOption {
	return func(o *LogOptions) {
		o.EnableLogHeader = enableLogHeader
	}
}

func EnableLogLink(enableLogLink bool) LogOption {
	return func(o *LogOptions) {
		o.EnableLogLink = enableLogLink
	}
}

func FlushInterval(flushInterval time.Duration) LogOption {
	return func(o *LogOptions) {
		o.FlushInterval = flushInterval
	}
}

func BrokerAddrs(brokerAddrs []string) LogOption {
	return func(o *LogOptions) {
		o.BrokerAddrs = brokerAddrs
	}
}

func Topic(topic string) LogOption {
	return func(o *LogOptions) {
		o.Topic = topic
	}
}

func Partition(partition int) LogOption {
	return func(o *LogOptions) {
		o.Partition = partition
	}
}

func Key(key string) LogOption {
	return func(o *LogOptions) {
		o.Key = key
	}
}

func IsSync(isSync bool) LogOption {
	return func(o *LogOptions) {
		o.IsSync = isSync
	}
}
func CallBack(callBack ICallBack) LogOption {
	return func(o *LogOptions) {
		o.CallBack = callBack
	}
}

func Endpoint(endpoint string) LogOption {
	return func(o *LogOptions) {
		o.Endpoint = endpoint
	}
}

func Project(project string) LogOption {
	return func(o *LogOptions) {
		o.Project = project
	}
}

func LogStore(logStore string) LogOption {
	return func(o *LogOptions) {
		o.LogStore = logStore
	}
}

func AccessId(accessId string) LogOption {
	return func(o *LogOptions) {
		o.AccessId = accessId
	}
}

func AccessSecret(accessSecret string) LogOption {
	return func(o *LogOptions) {
		o.AccessSecret = accessSecret
	}
}

func AliyunTopic(aliyunTopic string) LogOption {
	return func(o *LogOptions) {
		o.AliyunTopic = aliyunTopic
	}
}

func AliyunSource(aliyunSource string) LogOption {
	return func(o *LogOptions) {
		o.AliyunSource = aliyunSource
	}
}
