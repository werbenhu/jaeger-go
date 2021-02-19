# jaeger-go
A simple go tool for OpenTracing jaeger.

# install
To install this package, you need to install Go and set your Go workspace first.

```go get -u github.com/werbenhu/jaeger-go```

For more information about span, please refer to

opentracing-go https://github.com/opentracing/opentracing-go

opentracing.io https://opentracing.io/docs/overview/inject-extract/


# initialization
```
opt := &jaeger.Opt{
    ServiceName: "srv-two",
    HostPort: "127.0.0.1:6831",
    SamplerPerSecond: 3,
}
    
jaegerClient := jaeger.New(opt)
defer jaegerClient.Close()
```

# local-span

#### root-span
```
rootSpan := jaeger.NewSpan(context.Background(), "operation-name")
rootSpan.SetTag("tag-one", "tag-info")
rootSpan.SetBaggageItem("bag-one", "bag-info")
//...
rootSpan.Finish()
```

#### sub-span

```
rootSpan := jaeger.NewRootSpan("operation-name")
subSpan := rootSpan.Sub("sub-operation-name") 
brotherSpan := rootSpan.Sub("brother-operation-name")

brotherSpan.Finish()
grandsonSpan := subSpan.Sub("grandson-operation-name") 

grandsonSpan.Finish()
subSpan.Finish()
rootSpan.Finish()
```

# http-client-span
httpClient get span's carrier, and put it into http request header 
```
// new a http request
client := &http.Client{}
req, _ := http.NewRequest("GET","http://localhost:9002/server_two", nil)

// new a http req span
sub := span.Sub("http-req-span")
header := sub.GetHttpHeader()

// put trace-id(from carrier) into http header
for key, value := range header {
	req.Header.Add(key, value[0])
}

// send request
resp, _ := client.Do(req)
defer resp.Body.Close()
body, _ := ioutil.ReadAll(resp.Body)

//finish request span
sub.Finish()
fmt.Printf(string(body))
return sub
```

# http-server-span
httpServer get the header and new a remote sub-span by trace id
```
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/werbenhu/jaeger-go"
)

func main() {
    opt := &jaeger.Opt{
    	ServiceName: "srv-two",
    	HostPort: "127.0.0.1:6831",
    }
    
    jaegerClient := jaeger.New(opt)
    defer jaegerClient.Close()

	r := gin.Default()
	r.GET("/server_two", func(c *gin.Context) {
		//span := jaeger.NewSpanByTraceId(c.GetHeader("uber-trace-id"), "server-two-http-root")
		span := jaeger.NewSpanByHttpHeader(&c.Request.Header, "server-two-http-root")
		c.JSON(200, gin.H{
			"message": "server two response",
		})
		span.Finish()
	})
	r.Run("0.0.0.0:9002")
}
```

# mq span deliver
You can also deliver the span through mq(nsq,kafka,rabbitmq),

Producer put the trace-id or http-header into a message body.

Consumer get trace-id or http-header from the message and then new a span by trace-id or http-header.

