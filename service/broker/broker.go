package broker

import (
	"context"
	"encoding/json"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/walle/core"
)

func New(client *mixin.Client) core.BrokerService {
	return &brokerService{
		client: client,
	}
}

type brokerService struct {
	client    *mixin.Client
	members   []string
	threshold uint8
}

func (s *brokerService) Create(ctx context.Context, name, pin string) (*core.Broker, error) {
	key := mixin.GenerateEd25519Key()

	user, keystore, err := s.client.CreateUser(ctx, key, name)
	if err != nil {
		return nil, err
	}

	// update pin
	c, _ := mixin.NewFromKeystore(keystore)
	if err := c.ModifyPin(ctx, "", pin); err != nil {
		return nil, err
	}

	data, _ := json.Marshal(keystore)
	b := &core.Broker{
		WalletID: user.UserID,
		Pin:      pin,
		Data:     data,
	}

	return b, nil
}
