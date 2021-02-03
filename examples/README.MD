
Please install jaeger and nsq 

run server-one

`go run main.go nsq.go`

run server-two

`go run main.go nsq.go event.go`


http tracing test 
http://localhost:9001/req_one

nsq traceing test 
http://localhost:9001/nsq_one