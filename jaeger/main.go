package main

import (
	"flag"
	"play-arch/jaeger/server"
)

var svcType *string = flag.String("s", "grpc-a", "server type")

func main() {
	flag.Parse()

	switch *svcType {
	case "http-a":
		server.RunServerA()
	case "http-b":
		server.RunServerB()
	case "grpc-a":
		server.RunGrpcA()
	case "consume-a":
		server.RunConsumer()
	}
}
