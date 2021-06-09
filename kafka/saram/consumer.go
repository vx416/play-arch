package saram

import "github.com/Shopify/sarama"

func PartitionConsumer(topics []string) ([]sarama.PartitionConsumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.ChannelBufferSize = 256
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, cfg)
	if err != nil {
		return nil, err
	}

	pcs := make([]sarama.PartitionConsumer, 0, len(topics))
	for _, topic := range topics {
		pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if err != nil {
			return nil, err
		}
		pcs = append(pcs, pc)
	}

	return pcs, nil
}
