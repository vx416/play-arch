package client

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewSentinelClient(address []string) (*redis.Client, error) {

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: address,
		DB:            0,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
