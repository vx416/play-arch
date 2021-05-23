package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"play-arch/jaeger/trace"
	"play-arch/proto"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func HttpPing(ctx context.Context, hostPort string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ping-send")
	defer span.Finish()

	url := fmt.Sprintf("http://%s/ping", hostPort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if err := trace.Inject(span, req); err != nil {
		return "", err
	}
	return do(req)
}

func do(req *http.Request) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}
	return string(body), nil
}

func NewGrpcClient(tracer opentracing.Tracer) (proto.GreeterClient, error) {
	conn, err := grpc.Dial("127.0.0.1:8001",
		grpc.WithUnaryInterceptor(ClientInterceptor(tracer)),
		grpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}

	return proto.NewGreeterClient(conn), nil
}
