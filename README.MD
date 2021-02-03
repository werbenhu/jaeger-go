# jaeger-go
A simple go tool for OpenTracing jaeger.


# initialization

    jaegerCli := jaeger.NewJaeger("server-name", JaegerHostPort)
    defer jaegerCli.Close()


# new root span

```
rootSpan := jaeger.NewSpan(context.Background(), "operation-name")
//...
rootSpan.Finish()
```

or 
```
rootSpan := jaeger.NewRootSpan("operation-name")
//...
rootSpan.Finish()
```

# new sub span

```
subSpan := rootSpan.Sub("sub-operation-name") 
//...
subSpan.Finish()
```

# remote http sub span
client get span's carrier, and put it in http request header 
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

server get the header and new remote sub span by trace id
```
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/werbenhu/jaeger-go"
)

const (
	JaegerHostPort = "xxx:6831"
)

func main() {
	jaegerClient := jaeger.NewJaeger("srv-two", JaegerHostPort)
	defer jaegerClient.Close()

	r := gin.Default()
	r.GET("/server_two", func(c *gin.Context) {
		//span := jaeger.NewSpanByTraceId(c.GetHeader("uber-trace-id"), "server-two-http-root")
		span := jaeger.NewSpanByHttpHeader(&c.Request.Header, "")
		selfCall(span)
		c.JSON(200, gin.H{
			"message": "server two response",
		})
		span.Finish()
	})
	r.Run("0.0.0.0:9002")
}
```