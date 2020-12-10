package main

import (
	"sync"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/pkg/uuid"
	"github.com/fox-one/walle/cmd/walle-cli/config"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/pkg/cmd/builder"
)

func newBuilder(cfg config.Config) builder.Builder {
	return &cliBuilder{
		cfg:     cfg,
		traceID: uuid.New(),
	}
}

type cliBuilder struct {
	cfg config.Config

	db      *db.DB
	client  *mixin.Client
	traceID string

	mux sync.Mutex
}

func (b *cliBuilder) DB() *db.DB {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.db == nil {
		b.db = db.MustOpen(b.cfg.DB)
	}

	return b.db
}

func (b *cliBuilder) MixinClient() *mixin.Client {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.client == nil {
		b.client = provideMixinClient(b.cfg)
	}

	return b.client
}

func (b *cliBuilder) Brokers() core.BrokerStore {
	return provideBrokerStore(b.DB(), b.cfg.Broker.PinSecret)
}

func (b *cliBuilder) Brokerz() core.BrokerService {
	return provideBrokerService(b.MixinClient())
}

func (b *cliBuilder) Render() core.Render {
	return provideRender()
}

func (b *cliBuilder) Executor() string {
	return ""
}

func (b *cliBuilder) TraceID() string {
	return b.traceID
}
