package main

import (
	"context"
	"os"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/walle/cmd/walle-agent/config"
	"github.com/fox-one/walle/internal/build"
	"github.com/fox-one/walle/pkg/cmd/root"
	"github.com/fox-one/walle/pkg/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var flags struct {
	debug  bool
	config string
}

func init() {
	pflag.BoolVar(&flags.debug, "debug", false, "enable debug mode")
	pflag.StringVar(&flags.config, "config", "", "custom config file")
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
		cmdutil.HomePath(".walle_cli.yaml"),
	); err != nil {
		log.Fatalln("load config failed", err)
	}

	b := newBuilder(cfg)
	cmd := root.NewCmd(b, ver)
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalln(err)
	}
}
