package main

import (
	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/cmd/walle-agent/config"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/store/broker"
	"github.com/fox-one/walle/store/encrypt"
)

func provideDatabase(cfg config.Config) *db.DB {
	return db.MustOpen(cfg.DB)
}

func provideBrokerStore(db *db.DB, cfg config.Config) core.BrokerStore {
	enc, err := encrypt.New(cfg.Broker.PinSecret)
	if err != nil {
		panic(err)
	}

	return broker.New(db, enc)
}
