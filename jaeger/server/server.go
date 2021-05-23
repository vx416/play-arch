package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"play-arch/jaeger/trace"
	"play-arch/proto"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
)

func RunServerA() error {
	tracer, close, err := trace.Init("server-a")
	if err != nil {
		return err
	}
	defer close.Close()
	opentracing.SetGlobalTracer(tracer)

	traceMid := traceMiddleware(tracer)
	grpcClient, err := NewGrpcClient(tracer)
	if err != nil {
		return err
	}

	pub, err := NewPublisher("localhost:9092")
	if err != nil {
		return err
	}
	http.HandleFunc("/ping", traceMid(func(w http.ResponseWriter, r *http.Request) {
		response, err := HttpPing(r.Context(), "localhost:8082")
		if err != nil {

		}
		io.WriteString(w, fmt.Sprintf("%s -> %s", "server-a", response))
	}))
	http.HandleFunc("/grpc-ping", traceMid(func(w http.ResponseWriter, r *http.Request) {
		resp, err := grpcClient.SayHello(r.Context(), &proto.HelloRequest{})
		if err != nil {

		}
		io.WriteString(w, fmt.Sprintf("rec:%s", resp.Greeting))
	}))
	http.HandleFunc("/pub", traceMid(func(w http.ResponseWriter, r *http.Request) {
		Publish(r.Context(), "test_topic", "hello", tracer, pub)

		io.WriteString(w, "pub to test-topic")

	}))
	log.Printf("Listening on localhost:8081")
	return http.ListenAndServe(":8081", nil)
}

func RunServerB() error {
	tracer, close, err := trace.Init("server-b")
	if err != nil {
		return err
	}
	defer close.Close()
	opentracing.SetGlobalTracer(tracer)

	traceMid := traceMiddleware(tracer)
	http.HandleFunc("/ping", traceMid(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server-b")
	}))

	log.Printf("Listening on localhost:8082")
	return http.ListenAndServe(":8082", nil)
}

func RunGrpcA() error {
	tracer, close, err := trace.Init("grpc-a")
	if err != nil {
		return err
	}
	defer close.Close()
	opentracing.SetGlobalTracer(tracer)
	svc := grpc.NewServer(grpc.UnaryInterceptor(ServerInterceptor(tracer)))
	proto.RegisterGreeterServer(svc, rpcSvc{})
	ln, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		return err
	}
	return svc.Serve(ln)
}

func RunConsumer() error {
	tracer, close, err := trace.Init("consumer-a")
	if err != nil {
		return err
	}
	defer close.Close()
	opentracing.SetGlobalTracer(tracer)

	consumer, err := NewConsumerGroup("test", "localhost:9092")
	ctx, _ := context.WithCancel(context.Background())

	c := &Consumer{ready: make(chan bool)}

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Sarama err:%+v", err)
		}
	}()

	for {
		err := consumer.Consume(ctx, []string{"test_topic"}, c)
		log.Printf("err:%+v", err)
	}

}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Ready() {
	<-consumer.ready
}

func (consumer *Consumer) Setup(seesion sarama.ConsumerGroupSession) error {
	// close(consumer.ready)
	log.Println("Sarama consumer setup")

	return nil
}

func (consumer *Consumer) Cleanup(seesion sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(seesion sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for {
		ctx := context.Background()
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, offset = %d, topic = %s", string(message.Value), message.Offset, message.Topic)
			// log.Printf("Message claimed: header:%+v", message.Headers)
			// default:
			rw := HeaderReaderWriter{message.Headers}

			spanContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, rw)
			if err != nil {
				fmt.Println("consume extrac span faile, err ", err.Error())
			} else {
				span := opentracing.GlobalTracer().StartSpan(
					"consume "+message.Topic,
					ext.RPCServerOption(spanContext),
					opentracing.Tag{Key: string(ext.Component), Value: "consumer"},
					ext.SpanKindConsumer,
				)
				ctx = opentracing.ContextWithSpan(ctx, span)
				fmt.Println("consume span")
				ConsumeMsg(ctx, message)
				span.Finish()
			}

			// seesion.MarkMessage(message, "")

			// time.Sleep(30 * time.Microsecond)
		}
	}
}

func ConsumeMsg(ctx context.Context, msg *sarama.ConsumerMessage) {
	span, _ := opentracing.StartSpanFromContext(ctx, "consume msg")
	defer span.Finish()

	fmt.Println(string(msg.Value))
}
