package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasterSlave(t *testing.T) {
	ctx := context.Background()
	addrs := []string{":16379", ":36379", ":16378", ":36378", ":26379", ":26378"}

	client, err := NewMSClient(addrs)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	err = client.Set(ctx, "1345sdfg234", 1, -1).Err()
	assert.NoError(t, err)
	res, err := client.Get(ctx, "test").Result()
	assert.NoError(t, err)
	t.Log(res)
}
