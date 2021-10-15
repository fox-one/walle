package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/cmd/walle-agent/config"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/internal/build"
	"github.com/fox-one/walle/pkg/cmdutil"
	"github.com/fox-one/walle/worker/blaze"
	"github.com/fox-one/walle/worker/cashier"
	"github.com/fox-one/walle/worker/messenger"
	"github.com/fox-one/walle/worker/runner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

var flags struct {
	debug   bool
	config  string
	migrate bool
	version bool
}

func init() {
	pflag.BoolVar(&flags.debug, "debug", false, "enable debug mode")
	pflag.StringVar(&flags.config, "config", "", "custom config file")
	pflag.BoolVar(&flags.migrate, "migrate", false, "migrate database")
	pflag.BoolVarP(&flags.version, "version", "v", false, "print version")
}

func main() {
	ctx := context.Background()
	ver := build.Version

	pflag.Parse()

	log := logger.FromContext(ctx)
	if flags.debug {
		log.Logger.SetLevel(logrus.DebugLevel)
	}

	log.Logger.SetOutput(os.Stderr)

	var cfg config.Config
	if err := cmdutil.LoadConfig(
		&cfg,
		flags.config,
		"./config.yaml",
		cmdutil.HomePath(".walle_agent.yaml"),
	); err != nil {
		log.Fatalln("load config failed", err)
	}

	if flags.version {
		fmt.Println("Walle Agent", ver)
		return
	}

	database := provideDatabase(cfg)
	if flags.migrate {
		if err := db.Migrate(database); err != nil {
			log.Fatalln("migrate failed", err)
		}

		return
	}

	client := provideMixinClient(cfg)
	brokers := provideBrokerStore(database, cfg)
	brokerz := provideBrokerService(client)
	orders := provideOrderStore(database)
	orderz := provideOrderService(client, cfg)
	messages := provideMessageStore(database)
	messagez := provideMessageService(client)
	transfers := provideTransferStore(database)
	property := providePropertyStore(database)
	render := provideRender()

	var g errgroup.Group

	// blaze
	{
		w := blaze.New(client, brokers, brokerz, messages, render)
		g.Go(func() error { return w.Run(ctx) })
	}

	// cashier
	{
		w := cashier.New(brokers, transfers, func(broker *core.Broker) core.TransferService {
			return provideTransferService(broker, cfg)
		})
		g.Go(func() error { return w.Run(ctx) })
	}

	// messenger
	{
		w := messenger.New(messages, messagez)
		g.Go(func() error { return w.Run(ctx) })
	}

	// runner
	{
		w := runner.New(orders, orderz, transfers, property)
		g.Go(func() error { return w.Run(ctx) })
	}

	log.Printf("walle agent with version %q launched!\n", ver)

	if err := g.Wait(); err != nil {
		log.Fatalln("program terminated:", err)
	}
}
