package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func NewPublisher(address ...string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.ClientID = "testing"
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	produer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for msg := range produer.Successes() {
			log.Printf("producer: publish message on topic(%s)-(%d)", msg.Topic, msg.Partition)
		}
	}()

	go func() {
		for err := range produer.Errors() {
			log.Printf("producer: publish message on topic(%s)", err.Err)
			produer.Input() <- err.Msg
		}
	}()

	return produer, nil

}

func Publish(ctx context.Context, topic, msg string, tracer opentracing.Tracer, produer sarama.AsyncProducer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "pub "+topic, ext.SpanKindProducer)
	defer span.Finish()

	header := NewHeader()
	err := tracer.Inject(span.Context(), opentracing.TextMap, header)
	if err != nil {
		fmt.Println("inject pub spna failed, err ", err.Error())
	}
	fmt.Printf("HEADER:%+v\n", header.GetHeader())

	pm := &sarama.ProducerMessage{
		Headers: header.GetHeader(),
	}
	pm.Topic = topic
	pm.Value = sarama.StringEncoder(msg)
	produer.Input() <- pm
	log.Print("published!!!")
	return nil
}

func NewConsumerGroup(groupID string, address ...string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.ClientID = "test"
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Member.UserData = []byte("test_member_1")
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Net.MaxOpenRequests = 5
	config.RackID = "testinstanceid"

	return sarama.NewConsumerGroup(address, groupID, config)
}

func NewHeader() *HeaderReaderWriter {
	return &HeaderReaderWriter{headers: make([]*sarama.RecordHeader, 0, 1)}
}

//MDReaderWriter metadata Reader and Writer
type HeaderReaderWriter struct {
	headers []*sarama.RecordHeader
}

func (c HeaderReaderWriter) GetHeader() []sarama.RecordHeader {
	header := make([]sarama.RecordHeader, len(c.headers))
	for i := range c.headers {
		header[i] = *c.headers[i]
	}
	return header
}

// ForeachKey implements ForeachKey of opentracing.TextMapReader
func (c HeaderReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for _, header := range c.headers {
		err := handler(string(header.Key), string(header.Value))
		if err != nil {
			return err
		}
	}
	return nil
}

// Set implements Set() of opentracing.TextMapWriter
func (c *HeaderReaderWriter) Set(key, val string) {
	fmt.Println("SET ", key, "l", val)
	key = strings.ToLower(key)
	c.headers = append(c.headers, &sarama.RecordHeader{Key: []byte(key), Value: []byte(val)})
}
