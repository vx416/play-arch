package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	logFile := filepath.Join(dir, "../tmp/log")

	if logF := os.Getenv("LOGFILE"); logF != "" {
		logFile = logF
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = append(zapConfig.OutputPaths, logFile)
	logger, err = zapConfig.Build()
	if err != nil {
		log.Panicf("build logger failed, err:%+v", err)
		os.Exit(1)
	}
	logger.Info("logger recorde on " + logFile)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/ping", ping)
	server := http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:13333",
	}

	logger.Info("run server on 0.0.0.0:13333")
	if err := server.ListenAndServe(); err != nil {
		logger.Error("run server failed", zap.Error(err))
	}
}

func ping(resp http.ResponseWriter, req *http.Request) {
	logger.Info("access log", zap.Time("timestamp", time.Now()),
		zap.String("client_ip", req.RemoteAddr), zap.String("method", req.Method),
		zap.String("uri", req.RequestURI))
	resp.Write([]byte("pong!"))
}

func home(resp http.ResponseWriter, req *http.Request) {
	logger.Info("access log", zap.Time("timestamp", time.Now()),
		zap.String("client_ip", req.RemoteAddr), zap.String("method", req.Method),
		zap.String("uri", req.RequestURI))
	resp.Write([]byte("hello world!"))
}
