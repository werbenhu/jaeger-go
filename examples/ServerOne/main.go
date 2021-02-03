//
//  @File : main.go.go
//	@Author : WerBen
//  @Email : 289594665@qq.com
//	@Time : 2021/2/1 17:19 
//	@Desc : TODO ...
//

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/werbenhu/jaeger-go"
	"io/ioutil"
	"net/http"
	"time"
)
const (
	TopicName string = "test"
	TopicCh string = "test"
	JaegerHostPort = "218.91.230.203:6831"
)

// http传递trace-id
func twoCall(span *jaeger.Span) *jaeger.Span {

	// 通过http向server-two请求数据
	client := &http.Client{}
	req, _ := http.NewRequest("GET","http://localhost:9002/server_two",nil)

	// 生成一个请求的span
	sub := span.Sub("http-one-req")
	header := sub.GetHttpHeader()

	// 将当前span的trace-id传递到http header中
	for key, value := range header {
		req.Header.Add(key, value[0])
	}

	// 发送请求
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//结束当前请求的span
	sub.Finish()
	fmt.Printf(string(body))
	return sub
}

// nsq消息队列传递trace-id，其他的如kafka等MQ也可以类似
func nsqCall(span *jaeger.Span) *jaeger.Span {

	sub := span.Sub("http-one-req")
	header := sub.GetHttpHeader()

	// 将trace-id封装到消息中，由消息队列，传给消费者
	msg, _ := json.Marshal(header)
	// 生产消息
	Produce(TopicName, string(msg))

	sub.Finish()
	return sub
}

// 这里本地span，不用跨服务器
func selfCall(span *jaeger.Span) *jaeger.Span {
	sub := span.Sub("http-one-self")
	time.Sleep(time.Second)
	sub.Finish()
	return sub
}

func main() {

	//初始化jaeger
	jaegerCli := jaeger.NewJaeger("srv-one", JaegerHostPort)
	defer jaegerCli.Close()
	r := gin.Default()

	// server-one接口，收到请求会给server-two发送http请求
	r.GET("/req_one", func(c *gin.Context) {

		span := jaeger.NewSpan(context.Background(), "server-one-http-root")
		sub := twoCall(span)
		sub = selfCall(sub)

		time.Sleep(2 * time.Second)
		span.Finish()
		c.JSON(200, gin.H{
			"message": "hello_one response",
		})
	})

	// server-one接口，收到请求会给server-two发送消息（nsq消息队列）
	r.GET("/nsq_one", func(c *gin.Context) {
		span := jaeger.NewRootSpan("nsq-one-root")
		sub := nsqCall(span)
		sub = selfCall(sub)

		time.Sleep(2 * time.Second)
		span.Finish()
		c.JSON(200, gin.H{
			"message": "hello_one response",
		})
	})
	r.Run("0.0.0.0:9001")
}