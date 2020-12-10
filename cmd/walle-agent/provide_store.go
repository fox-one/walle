package main

import (
	"github.com/fox-one/pkg/property"
	"github.com/fox-one/pkg/store/db"
	propertystore "github.com/fox-one/pkg/store/property"
	"github.com/fox-one/walle/cmd/walle-agent/config"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/store/broker"
	"github.com/fox-one/walle/store/encrypt"
	"github.com/fox-one/walle/store/message"
	"github.com/fox-one/walle/store/order"
	"github.com/fox-one/walle/store/transfer"
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

func provideOrderStore(db *db.DB) core.OrderStore {
	return order.New(db)
}

func provideMessageStore(db *db.DB) core.MessageStore {
	return message.New(db)
}

func providePropertyStore(db *db.DB) property.Store {
	return propertystore.New(db)
}

func provideTransferStore(db *db.DB) core.TransferStore {
	return transfer.New(db)
}
