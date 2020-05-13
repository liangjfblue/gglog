/**
 *
 * @author liangjf
 * @create on 2020/5/12
 * @version 1.0
 */
package aliyun

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/liangjfblue/gglog/utils"

	"github.com/gogo/protobuf/proto"

	sls "github.com/aliyun/aliyun-log-go-sdk"

	"github.com/aliyun/aliyun-log-go-sdk/producer"

	"github.com/liangjfblue/gglog/vlog"
)

type AliyunLog struct {
	opts   vlog.LogOptions
	onceDo sync.Once

	callback *Callback
	producer *producer.Producer
	slsLog   *sls.Log
}

var (
	_serverIp string
	_hostName string
)

func init() {
	_serverIp, _ = utils.ExternalIP()
	_hostName, _ = os.Hostname()
}

func NewAliyunLog(opts ...vlog.LogOption) vlog.Log {
	options := vlog.NewOptions(opts...)

	return &AliyunLog{
		opts: options,
	}
}

func (a *AliyunLog) Name() string {
	return "aliyun-log"
}

func (a *AliyunLog) Init(opts ...vlog.LogOption) {
	for _, o := range opts {
		o(&a.opts)
	}

	a.onceDo.Do(func() {
		a.initProducer()
	})

}

func (a *AliyunLog) Run() {
	a.producer.Start()
}

func (a *AliyunLog) Stop() {
	a.producer.Close(500)
}

func (a *AliyunLog) Debug(format string, args ...interface{}) {
	if a.opts.Level <= utils.LevelD {
		a.sendLog(utils.Debug, fmt.Sprintf(format, args...))
	}
}

func (a *AliyunLog) Info(format string, args ...interface{}) {
	if a.opts.Level <= utils.LevelI {
		a.sendLog(utils.Info, fmt.Sprintf(format, args...))
	}
}

func (a *AliyunLog) Warn(format string, args ...interface{}) {
	if a.opts.Level <= utils.LevelW {
		a.sendLog(utils.Warn, fmt.Sprintf(format, args...))
	}
}

func (a *AliyunLog) Error(format string, args ...interface{}) {
	if a.opts.Level <= utils.LevelE {
		a.sendLog(utils.Error, fmt.Sprintf(format, args...))
	}
}

func (a *AliyunLog) Access(format string, args ...interface{}) {
	if a.opts.OpenAccessLog == true {
		a.sendLog(utils.Access, fmt.Sprintf(format, args...))
	}
}

func (a *AliyunLog) InterfaceAvgDuration(format string, args ...interface{}) {
	if a.opts.OpenInterfaceAvgDurationLog == true {
		a.sendLog(utils.IAVGD, fmt.Sprintf(format, args...))
	}
}

func (a *AliyunLog) FlushLog() {

}

func (a *AliyunLog) initProducer() {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = a.opts.Endpoint
	producerConfig.AccessKeyID = a.opts.AccessId
	producerConfig.AccessKeySecret = a.opts.AccessSecret
	producerConfig.Retries = 5
	sls.GlobalForceUsingHTTP = true

	a.callback = new(Callback)
	a.producer = producer.InitProducer(producerConfig)
	a.slsLog = &sls.Log{}
}

func (a *AliyunLog) sendLog(level string, desc string) {
	now := time.Now().UnixNano() / 1000000
	_, file, line, _ := runtime.Caller(2)
	slash := strings.LastIndex(file, "/")
	if slash >= 0 {
		file = file[slash+1:]
	}

	a.slsLog = &sls.Log{}
	a.slsLog.Time = proto.Uint32(uint32(time.Now().Unix()))
	a.slsLog.Contents = []*sls.LogContent{
		{
			Key:   proto.String("hostname"),
			Value: proto.String(_hostName),
		},
		{
			Key:   proto.String("ip"),
			Value: proto.String(_serverIp),
		},
		{
			Key:   proto.String("location"),
			Value: proto.String(file + ":" + fmt.Sprint(line)),
		},
		{
			Key:   proto.String("tm"),
			Value: proto.String(fmt.Sprint(now)),
		},
		{
			Key:   proto.String("level"),
			Value: proto.String(level),
		},
		{
			Key:   proto.String("desc"),
			Value: proto.String(desc),
		},
	}

	err := a.producer.SendLogWithCallBack(
		a.opts.Project, a.opts.LogStore,
		a.opts.AliyunTopic, a.opts.AliyunSource, a.slsLog, a.callback,
	)
	if err != nil {
		log.Println("send log err: ", err.Error())
	}
}

type Callback struct {
}

func (callback *Callback) Success(result *producer.Result) {

}

func (callback *Callback) Fail(result *producer.Result) {
	log.Println(result.IsSuccessful())        // 获得发送日志是否成功
	log.Println(result.GetErrorCode())        // 获得最后一次发送失败错误码
	log.Println(result.GetErrorMessage())     // 获得最后一次发送失败信息
	log.Println(result.GetReservedAttempts()) // 获得producerBatch 每次尝试被发送的信息
	log.Println(result.GetRequestId())        // 获得最后一次发送失败请求Id
	log.Println(result.GetTimeStampMs())      // 获得最后一次发送失败请求时间
}
