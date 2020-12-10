package main

import (
	"context"
	"os"
	"path"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/walle/cmd/walle-cli/config"
	"github.com/fox-one/walle/internal/build"
	"github.com/fox-one/walle/pkg/cmd/root"
	"github.com/mitchellh/go-homedir"
)

func main() {
	ctx := context.Background()
	ver := build.Version

	log := logger.FromContext(ctx)
	log.Logger.SetOutput(os.Stderr)

	configFile := lookupConfigFilePath()
	log.Debugf("load config from %q", configFile)

	var cfg config.Config
	if err := config.Load(configFile, &cfg); err != nil {
		log.Fatalln("load config failed", err)
	}

	b := newBuilder(cfg)
	cmd := root.NewCmd(b, ver)
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalln(err)
	}
}

func lookupConfigFilePath() string {
	// lookup ./config.yaml
	{
		filename := "./config.yaml"
		if exists(filename) {
			return filename
		}
	}

	// lookup ~/.walle.yaml
	{
		dir, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		filename := path.Join(dir, ".walle_cli.yaml")
		if exists(filename) {
			return filename
		}
	}

	return ""
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
