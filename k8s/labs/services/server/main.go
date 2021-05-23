package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type config struct {
	AppName string
	Gateway string
	Addr    string
}

func newCfg() *config {
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName, _ = os.Hostname()
	}
	cfg.AppName = appName
	cfg.Gateway = os.Getenv("GATEWAY_ADDR")
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	cfg.Addr = addr

	return cfg
}

var (
	count int
	cfg   *config = &config{}
)

func main() {
	newCfg()
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", CountMiddleware(echo))
	mux.HandleFunc("/count", CountMiddleware(getCount))

	serv := http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}
	log.Printf("server runing on %s", cfg.Addr)
	if err := serv.ListenAndServe(); err != nil {
		log.Panicf("run server failed, err:%+v", err)
		os.Exit(1)
	}
}

func CountMiddleware(fn func(resp http.ResponseWriter, req *http.Request)) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		count++
		log.Printf("%s receive a request", cfg.AppName)
		fn(resp, req)
	}
}

func echo(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)

	io.WriteString(resp,
		fmt.Sprintf("Hello, %s, I am %s.", req.RemoteAddr, cfg.AppName),
	)

}

func callGateway() {

}

func getCount(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)

	io.WriteString(resp,
		fmt.Sprintf("request count:%d", count),
	)

}
