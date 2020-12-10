package main

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/walle/cmd/walle-agent/config"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/render/terminal"
	"github.com/fox-one/walle/service/broker"
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
	return terminal.Simple()
}
