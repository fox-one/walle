package main

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/walle/cmd/walle-cli/config"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/service/broker"
	"github.com/fox-one/walle/service/render"
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
	return render.Simple()
}
