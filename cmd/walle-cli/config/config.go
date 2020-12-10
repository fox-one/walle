package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/config"
	"github.com/fox-one/pkg/store/db"
)

type Config struct {
	App    App       `json:"app,omitempty"`
	DB     db.Config `json:"db,omitempty"`
	Broker Broker    `json:"broker,omitempty"`
}

type App struct {
	*mixin.Keystore
}

type Broker struct {
	PinSecret string `json:"pin_secret,omitempty"`
}

func Load(cfgFile string, cfg *Config) error {
	config.AutomaticLoadEnv("WALLE")
	if err := config.LoadYaml(cfgFile, cfg); err != nil {
		return err
	}

	return nil
}
