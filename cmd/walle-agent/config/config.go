package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
)

type (
	Config struct {
		App      App       `json:"app,omitempty"`
		DB       db.Config `json:"db,omitempty"`
		Broker   Broker    `json:"broker,omitempty"`
		Multisig Multisig  `json:"multisig,omitempty"`
	}

	App struct {
		*mixin.Keystore
	}

	Broker struct {
		PinSecret string `json:"pin_secret,omitempty"`
	}

	Multisig struct {
		Members   []string `json:"members,omitempty"`
		Threshold uint8    `json:"threshold,omitempty"`
		VerifyKey string   `json:"verify_key,omitempty"`
	}
)
