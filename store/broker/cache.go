package broker

import (
	"context"
	"time"

	"github.com/fox-one/walle/core"
	"github.com/patrickmn/go-cache"
)

func Cache(brokers core.BrokerStore) core.BrokerStore {
	return &cacheStore{
		BrokerStore: brokers,
		brokers:     cache.New(time.Hour, time.Minute),
	}
}

type cacheStore struct {
	core.BrokerStore
	brokers *cache.Cache
}

func (c *cacheStore) Find(ctx context.Context, walletID string) (*core.Broker, error) {
	if b, ok := c.brokers.Get(walletID); ok {
		return b.(*core.Broker), nil
	}

	b, err := c.BrokerStore.Find(ctx, walletID)
	if err != nil {
		return nil, err
	}

	c.brokers.SetDefault(walletID, b)
	return b, nil
}
