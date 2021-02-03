package main

/*
@Time : 2020/09/04
@Author : WerBen
@File : event.go
@Desc : dev's event via nsq
*/

import (
	"encoding/json"
	"fmt"
	"github.com/werbenhu/jaeger-go"
	"net/http"
)

const (
	TopicName  string = "test"
	TopicCh string = "test"
)

// 消费者消息处理，处理server-one发送过来的消息
func eventHandler(message string) error {
	fmt.Printf("event msg:%s\n", message)
	var header http.Header
	json.Unmarshal([]byte(message), &header)
	span := jaeger.NewSpanByHttpHeader(&header, "nsq_two")
	selfCall(span)
	span.Finish()
	return nil
}

func InitEvent() {
	jaeger.NewJaeger("srv-nsq", JaegerHostPort)
	//jaegerEx := NewJaeger("srv-nsq")
	//defer jaegerEx.Close()
	Consume(TopicName, TopicCh, eventHandler)
}
