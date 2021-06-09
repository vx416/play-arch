package client

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewMSClient(address []string) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		ClusterSlots: func(c context.Context) ([]redis.ClusterSlot, error) {
			return []redis.ClusterSlot{
				{
					Start: 0,
					End:   5460,
					Nodes: []redis.ClusterNode{{"master1", address[0]}, {"slave1", address[1]}},
				},
				{
					Start: 5461,
					End:   10922,
					Nodes: []redis.ClusterNode{{"master2", address[2]}, {"slave2", address[3]}},
				},
				{
					Start: 10923,
					End:   16383,
					Nodes: []redis.ClusterNode{{"master3", address[4]}, {"slave3", address[1]}},
				},
			}, nil
		},
		RouteByLatency: true,
		ReadOnly:       true,
		Password:       "mypass",
		MaxRedirects:   3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return client, nil
}
