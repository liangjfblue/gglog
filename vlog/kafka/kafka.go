/*
@Time : 2020/5/10 22:49
@Author : liangjiefan
*/
package kafka

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"

	"github.com/liangjfblue/gglog/vlog"
)

type kafkaLog struct {
	opts           vlog.LogOptions
	channel        chan string
	_syncProducer  sarama.SyncProducer
	_asyncProducer sarama.AsyncProducer
	stopChan       chan struct{}
	isStop         bool
	once           sync.Once
}

var (
	_serverIp string
	_hostName string
)

func init() {
	_serverIp, _ = externalIP()
	_hostName, _ = os.Hostname()
}

func NewKafkaLog(opts ...vlog.LogOption) vlog.Log {
	options := vlog.NewOptions(opts...)

	return &kafkaLog{
		opts:     options,
		channel:  make(chan string, 10),
		stopChan: make(chan struct{}, 1),
		isStop:   false,
	}
}

func (s *kafkaLog) Name() string {
	return s.opts.Name
}

func (s *kafkaLog) Init(opts ...vlog.LogOption) {
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		s.initKafkaLog()
	})
}

func (s *kafkaLog) Options() vlog.LogOptions {
	return s.opts
}

func (s *kafkaLog) String() string {
	return "kafkaLog"
}

const (
	levelD = iota + 1
	levelI
	levelW
	levelE
)

const (
	Debug  = "debug"
	Info   = "info"
	Warn   = "warn"
	Error  = "error"
	Access = "access"
	IAVGD  = "iavgd"
)

func (s *kafkaLog) Debug(format string, args ...interface{}) {
	if s.opts.Level <= levelD {
		s.toChannel(Debug, fmt.Sprintf(format, args...))
	}
}

func (s *kafkaLog) Info(format string, args ...interface{}) {
	if s.opts.Level <= levelI {
		s.toChannel(Info, fmt.Sprintf(format, args...))
	}
}

func (s *kafkaLog) Warn(format string, args ...interface{}) {
	if s.opts.Level <= levelW {
		s.toChannel(Warn, fmt.Sprintf(format, args...))
	}
}

func (s *kafkaLog) Error(format string, args ...interface{}) {
	if s.opts.Level <= levelE {
		s.toChannel(Error, fmt.Sprintf(format, args...))
	}
}

func (s *kafkaLog) Access(format string, args ...interface{}) {
	if s.opts.OpenAccessLog == true {
		s.toChannel(Access, fmt.Sprintf(format, args...))
	}
}

func (s *kafkaLog) InterfaceAvgDuration(format string, args ...interface{}) {
	if s.opts.OpenInterfaceAvgDurationLog == true {
		s.toChannel(IAVGD, fmt.Sprintf(format, args...))
	}
}

func (s *kafkaLog) FlushLog() {
	log.Println("nothing to do")
}

func (s *kafkaLog) Run() {
	if s.opts.IsSync {
		s.sendLogToKafkaSync()
	} else {
		s.sendLogToKafkaAsync()
	}
}

func (s *kafkaLog) Stop() {
	if s.isStop {
		log.Println("had stop")
		return
	}
	s.stopChan <- struct{}{}
}

func (s *kafkaLog) initKafkaLog() {
	var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 3
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = 100 * time.Millisecond

	if s.opts.IsSync {
		s._syncProducer, err = sarama.NewSyncProducer(s.opts.BrokerAddrs, config)
		if err != nil {
			log.Printf("create producer error, because is %s\n", err.Error())
		}
	} else {
		s._asyncProducer, err = sarama.NewAsyncProducer(s.opts.BrokerAddrs, config)
		if err != nil {
			log.Printf("create producer error, because is %s\n", err.Error())
		}
	}
}

func (s *kafkaLog) sendLogToKafkaSync() {
	defer func() {
		close(s.channel)
		close(s.stopChan)

		if err := s._syncProducer.Close(); err != nil {
			log.Printf("close producer fail, because is %s\n", err.Error())
		}
	}()

	for {
		select {
		case content, ok := <-s.channel:
			if !ok {
				log.Println("chan is closed")
				return
			}
			key := s.opts.Key + "_" + time.Now().Format("2006-01-02 15:04:05")

			msg := sarama.ProducerMessage{
				Topic: s.opts.Topic,
				Value: sarama.StringEncoder(content),
				Key:   sarama.StringEncoder(key),
			}

			log.Println(msg)
			partition, offset, err := s._syncProducer.SendMessage(&msg)
			if err != nil {
				log.Printf("send msg to kafka fail, because is %s", err.Error())
			}

			logMsg := fmt.Sprintf("parttion, offset: %d %d", partition, offset)
			log.Print(logMsg)
		}
	}
}

func (s *kafkaLog) sendLogToKafkaAsync() {
	defer func() {
		if err := s._asyncProducer.Close(); err != nil {
			log.Printf("close producer fail, because is %s\n", err.Error())
		}
	}()

	for {
		if content, ok := <-s.channel; ok {
			key := s.opts.Key + "_" + time.Now().Format("2006-01-02 15:04:05")

			msg := sarama.ProducerMessage{
				Topic: s.opts.Topic,
				Value: sarama.StringEncoder(content),
				Key:   sarama.StringEncoder(key),
			}

			s._asyncProducer.Input() <- &msg

			select {
			case <-s._asyncProducer.Successes():
			case err := <-s._asyncProducer.Errors():
				//TODO retry to send log to kafka
				log.Println("produced message error: ", err)
			default:
				log.Println("produced message close / asyncClose")
			}
		} else {
			log.Print("read data form chan error")
		}
	}
}

func (s *kafkaLog) toChannel(level string, desc string) {
	now := time.Now().UnixNano() / 1000000
	_, file, line, _ := runtime.Caller(2)
	slash := strings.LastIndex(file, "/")
	if slash >= 0 {
		file = file[slash+1:]
	}

	str := fmt.Sprintf("{\"ip\":\"%s\", \"location\":\"%s:%d\", \"tm\":%d, \"level\":\"%s\", \"desc\":\"%s\", \"hostname\":\"%s\"}",
		_serverIp, file, line, now, level, desc, _hostName)

	if !s.isStop {
		select {
		case s.channel <- str:
		case <-s.stopChan:
			s.isStop = true
		}
	}
}
