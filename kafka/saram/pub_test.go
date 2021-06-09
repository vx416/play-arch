package saram

import (
	"testing"
	"unsafe"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func BenchmarkPub(b *testing.B) {
	b.Log(len(bigStr))
	b.Log(unsafe.Sizeof(bigStr))
	p, err := NewSyncProducer(nil)
	if err != nil {
		b.Fatal(err)
	}
	cfg := sarama.NewConfig()
	cfg.Producer.Compression = sarama.CompressionLZ4
	p2, err := NewSyncProducer(cfg)
	if err != nil {
		b.Fatal(err)
	}
	cfg = sarama.NewConfig()
	cfg.Producer.Compression = sarama.CompressionSnappy
	p3, err := NewSyncProducer(cfg)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("compress_snappy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p3.SendMessage(&sarama.ProducerMessage{
				Topic: "test_3",
				Value: sarama.StringEncoder(bigStr),
			})
		}
	})

	b.Run("compress_lz4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p2.SendMessage(&sarama.ProducerMessage{
				Topic: "test_2",
				Value: sarama.StringEncoder(bigStr),
			})
		}
	})

	b.Run("no_compress", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p.SendMessage(&sarama.ProducerMessage{
				Topic: "test_1",
				Value: sarama.StringEncoder(bigStr),
			})
		}
	})
}

func TestPub(t *testing.T) {
	p, err := NewSyncProducer(nil)
	assert.NoError(t, err)
	p.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("hello"),
	})
}

func BenchmarkConsumer(b *testing.B) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.ChannelBufferSize = 256
	// consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, cfg)
	// if err != nil {
	// 	b.Fatal(err)
	// }
}
