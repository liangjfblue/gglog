/*
@Time : 2020/5/10 22:49
@Author : liangjiefan
*/
package kafka

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"

	"github.com/liangjfblue/gglog/vlog"
)

type kafkaLog struct {
	opts    vlog.LogOptions
	channel chan string
	once    sync.Once
}

func NewKafkaLog(opts ...vlog.LogOption) vlog.Log {
	options := vlog.NewOptions(opts...)

	return &kafkaLog{
		opts:    options,
		channel: make(chan string, 1000),
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

func (s *kafkaLog) Debug(format string, args ...interface{}) {
	if s.opts.Level <= levelD {
		s.channel <- fmt.Sprintf(format, args...)
	}
}

func (s *kafkaLog) Info(format string, args ...interface{}) {
	if s.opts.Level <= levelI {
		s.channel <- fmt.Sprintf(format, args...)
	}
}

func (s *kafkaLog) Warn(format string, args ...interface{}) {
	if s.opts.Level <= levelW {
		s.channel <- fmt.Sprintf(format, args...)
	}
}

func (s *kafkaLog) Error(format string, args ...interface{}) {
	if s.opts.Level <= levelE {
		s.channel <- fmt.Sprintf(format, args...)
	}
}

func (s *kafkaLog) Access(format string, args ...interface{}) {
	if s.opts.OpenAccessLog == true {
		s.channel <- fmt.Sprintf(format, args...)
	}
}

func (s *kafkaLog) InterfaceAvgDuration(format string, args ...interface{}) {
	if s.opts.OpenInterfaceAvgDurationLog == true {
		s.channel <- fmt.Sprintf(format, args...)
	}
}

func (s *kafkaLog) FlushLog() {
	log.Println("nothing to do")
}

func (s *kafkaLog) initKafkaLog() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 3
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	if s.opts.IsSync {
		producer, err := sarama.NewSyncProducer(s.opts.BrokerAddrs, config)
		if err != nil {
			log.Printf("create producer error, because is %s\n", err.Error())
		}
		go s.sendLogToKafkaSync(producer)
	} else {
		producer, err := sarama.NewAsyncProducer(s.opts.BrokerAddrs, config)
		if err != nil {
			log.Printf("create producer error, because is %s\n", err.Error())
		}
		go s.sendLogToKafkaAsync(producer)
	}
}

func (s *kafkaLog) sendLogToKafkaSync(producer sarama.SyncProducer) {
	defer func() {
		if err := producer.Close(); err != nil {
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

			partition, offset, err := producer.SendMessage(&msg)
			if err != nil {
				log.Printf("send msg to kafka fail, because is %s", err.Error())
			}

			logMsg := fmt.Sprintf("Msg is stored in parttion, offset: %d %d", partition, offset)
			log.Print(logMsg)
		} else {
			log.Print("read data form channel error")
		}
	}
}

func (s *kafkaLog) sendLogToKafkaAsync(producer sarama.AsyncProducer) {
	defer func() {
		if err := producer.Close(); err != nil {
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

			producer.Input() <- &msg

			select {
			case <-producer.Successes():
			case err := <-producer.Errors():
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
