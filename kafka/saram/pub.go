package saram

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func NewSyncProducer(cfg *sarama.Config) (sarama.SyncProducer, error) {

	if cfg == nil {
		cfg = sarama.NewConfig()
		// cfg.Producer.Compression = sarama.CompressionSnappy
	}
	cfg.Producer.Return.Successes = true

	return sarama.NewSyncProducer([]string{"localhost:9092"}, cfg)
}

func PubMessages(topics []string, msg string, n int) error {
	pro, err := NewSyncProducer(nil)
	if err != nil {
		return err
	}

	for _, t := range topics {
		go func(topic string) {
			for i := 1; i <= n; i++ {
				_, _, proErr := pro.SendMessage(&sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(fmt.Sprintf("%s-%d", msg, i)),
				})
				if err == nil && proErr != nil {
					err = proErr
				}
			}
		}(t)
	}

	return err
}
