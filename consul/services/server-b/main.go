package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	port    = flag.Int("port", 90, "server port")
	portA   = flag.Int("port_a", 8811, "server a port")
	appName = flag.String("app_name", "server_b_1", "application name")
)

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

	data, err := sendReqToServerA()
	if err != nil {
		log.Panicf("err:%+v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		io.WriteString(resp, "failed")
		return
	}

	io.WriteString(resp, "Hello, I am "+*appName+". This is my brother's greeting: "+string(data))
	resp.WriteHeader(http.StatusOK)
}

func sendReqToServerA() ([]byte, error) {
	resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(*portA))
	if err != nil {
		log.Panicf("err:%+v", err)
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
