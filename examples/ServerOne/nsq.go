package main

/*
@Time : 2020/09/04
@Author : werben
@File : nsq.go
@Desc : 消息队列nsq的封装
*/

import (
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var nsqConfig *nsq.Config

var producer *nsq.Producer

//枷锁，防止协程间资源竞争
var mutexNsq sync.Mutex

type MessageHandler func(message string) error

//封装nsq handler对象
type DefaultHandler struct {
	//自定义消息处理函数
	messageHandler MessageHandler
}

func (handler *DefaultHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	return handler.messageHandler(string(m.Body))
}

//获取nsg config
func getNsqConfig() *nsq.Config {
	mutexNsq.Lock()
	defer mutexNsq.Unlock()

	if nil == nsqConfig {
		nsqConfig = nsq.NewConfig()
	}
	return nsqConfig
}

//获取生产者
func getProducer() *nsq.Producer {
	config := getNsqConfig()
	mutexNsq.Lock()
	defer mutexNsq.Unlock()
	if producer == nil {
		pro, err := nsq.NewProducer("218.91.230.20:4150", config)
		logger := log.New(os.Stderr, "", log.Flags())
		pro.SetLogger(logger, nsq.LogLevelError)
		if err != nil {
			log.Println(err)
			return nil
		}
		producer = pro
	}
	return producer
}

//封装nsq生产接口
func Produce(topic string, message string) {
	producer := getProducer()
	//defer producer.Stop()
	if nil != producer {
		err := producer.Publish(topic, []byte(message))
		if err != nil {
			log.Println(err)
		}
	}
}

//封装nsq消费接口
func Consume(topic string, channel string, handler MessageHandler) {
	go func(topic string, channel string, handler MessageHandler) {
		config := getNsqConfig()
		consumer, err := nsq.NewConsumer(topic, channel, config)
		logger := log.New(os.Stderr, "", log.Flags())
		consumer.SetLogger(logger, nsq.LogLevelError)
		defer consumer.Stop()

		if err != nil {
			log.Println(err)
		}
		consumer.AddHandler(&DefaultHandler{
			messageHandler : handler,
		})

		err = consumer.ConnectToNSQLookupd("218.91.230.20:4161")
		if err != nil {
			log.Println(err)
		}

		//消费接口，要一直运行
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
		<- signalCh
	}(topic, channel, handler)
}

