package core

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/jmoiron/sqlx/types"
)

type (
	Broker struct {
		ID        int64          `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time      `json:"created_at,omitempty"`
		UpdatedAt time.Time      `json:"updated_at,omitempty"`
		WalletID  string         `sql:"size:36" json:"wallet_id,omitempty"`
		UserID    string         `sql:"size:36" json:"user_id,omitempty"`
		Pin       string         `sql:"size:255" json:"pin,omitempty"`
		Data      types.JSONText `sql:"type:varchar(512)" json:"data,omitempty"`
	}

	BrokerStore interface {
		Create(ctx context.Context, broker *Broker) error
		Find(ctx context.Context, walletID string) (*Broker, error)
	}

	BrokerService interface {
		Create(ctx context.Context, name, pin string) (*Broker, error)
	}
)

func (b *Broker) MixinClient() *mixin.Client {
	var store mixin.Keystore
	_ = json.Unmarshal(b.Data, &store)

	client, err := mixin.NewFromKeystore(&store)
	if err != nil {
		panic(err)
	}

	return client
}
