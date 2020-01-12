# RateLimit

Best performance rate limit library, provide promise lock limiter and standard-like limiter. Comes with standard-http compatible middleware.

![](https://github.com/tspn/rate-limit/raw/master/Screen%20Shot%202563-01-13%20at%2000.19.49.png)

The images above shows performance of library compares with standard library Limiter. By x-axis is concurrent of http request and y-axis is % of error. Lower is better.

### Get
```bash
$ go get github.com/tspn/rate_limit
```

### Usage
- Http middleware with promise every request will get served
```go
func main(){
  rate := 1000.0 // 1000 request per second
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(writer, "Hello")
  })

  pipeline := rate_limit.New(rate).PromiseMiddleware(mux)
  http.ListenAndServe(":8080", pipeline)
}
```
