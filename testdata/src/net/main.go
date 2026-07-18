package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
)

func netExercise() {
	_, _ = net.Dial("tcp", "localhost:0")
	_, _ = net.DialTimeout("tcp", "localhost:0", 0)
	_, _ = net.Listen("tcp", ":0")
	_, _ = net.ResolveTCPAddr("tcp", "localhost:0")
	host, port, _ := net.SplitHostPort("localhost:8080")
	_ = host
	_ = port
	_ = net.JoinHostPort("localhost", "8080")
}

func urlExercise() {
	_, _ = url.Parse("https://example.com/path?q=1")
	v := url.Values{}
	v.Set("key", "value")
	v.Add("key", "value2")
	_ = v.Get("key")
	_ = v.Encode()
}

func httpExercise() {
	_, _ = http.Get("http://example.com")
	_, _ = http.Post("http://example.com", "text/plain", nil)
}

func main() {
	netExercise()
	urlExercise()
	httpExercise()
	fmt.Println("net e2e test done")
}
