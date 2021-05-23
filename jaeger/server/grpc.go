package server

import (
	"context"
	"play-arch/proto"
)

type rpcSvc struct {
}

func (svc rpcSvc) SayHello(context.Context, *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Greeting: "hello",
	}, nil
}
