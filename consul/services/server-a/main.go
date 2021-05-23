package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	port    = flag.Int("port", 80, "server port")
	appName = flag.String("app_name", "server_a_1", "application name")
)

func init() {
	flag.Parse()
}

func main() {
	http.HandleFunc("/", logMiddleware(root))
	http.HandleFunc("/ping", logMiddleware(pingPong))

	if err := http.ListenAndServe(":"+strconv.Itoa(*port), nil); err != nil {
		log.Fatalf("run server failed, err:%+v", err)
	}
}

func logMiddleware(next func(resp http.ResponseWriter, req *http.Request)) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("the app %s receive the request", *appName)
		next(resp, req)
	}
}

func pingPong(resp http.ResponseWriter, req *http.Request) {
	ip := req.RemoteAddr
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, ip+" pong!")
}

func root(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
	io.WriteString(resp, "Hello, I am "+*appName+".")
}
