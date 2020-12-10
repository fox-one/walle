package main

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/walle/cmd/walle-agent/config"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/render/messenger"
	"github.com/fox-one/walle/render/terminal"
	"github.com/fox-one/walle/service/broker"
	"github.com/fox-one/walle/service/message"
	"github.com/fox-one/walle/service/order"
	"github.com/fox-one/walle/service/transfer"
)

func provideMixinClient(cfg config.Config) *mixin.Client {
	client, err := mixin.NewFromKeystore(cfg.App.Keystore)
	if err != nil {
		panic(err)
	}

	return client
}

func provideBrokerService(client *mixin.Client) core.BrokerService {
	return broker.New(client)
}

func provideRender() core.Render {
	r := terminal.Verbose()
	return messenger.Wrap(r)
}

func provideOrderService(client *mixin.Client, cfg config.Config) core.OrderService {
	return order.New(client, order.Config{VerifyKey: cfg.Multisig.VerifyKey})
}

func provideTransferService(broker *core.Broker, cfg config.Config) core.TransferService {
	return transfer.New(broker.MixinClient(), transfer.Config{
		Pin:       broker.Pin,
		Members:   cfg.Multisig.Members,
		Threshold: cfg.Multisig.Threshold,
	})
}

func provideMessageService(client *mixin.Client) core.MessageService {
	return message.New(client)
}
