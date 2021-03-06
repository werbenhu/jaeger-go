package main

import (
	"github.com/gin-gonic/gin"
	"github.com/werbenhu/jaeger-go"
	"time"
)

func selfCall(span *jaeger.Span) *jaeger.Span {
	sub := span.Sub("self-two-call-1")
	time.Sleep(time.Second)

	other := sub.Sub("self-two-call-2")
	time.Sleep(time.Second)

	other.Finish()
	sub.Finish()
	return sub
}

func main() {

	opt := &jaeger.Opt{
		ServiceName: "srv-two",
		HostPort: "127.0.0.1:6831",
		SamplerPerSecond: 3,
	}

	jaegerClient := jaeger.New(opt)
	defer jaegerClient.Close()
	InitEvent()

	r := gin.Default()
	r.GET("/server_two", func(c *gin.Context) {
		span := jaeger.NewSpanByTraceId(c.GetHeader("uber-trace-id"), "server-two-http-root")
		jaeger.NewSpanByHttpHeader(&c.Request.Header, "")
		selfCall(span)
		c.JSON(200, gin.H{
			"message": "server two response",
		})
		span.Finish()
	})
	r.Run("0.0.0.0:9002")
}
