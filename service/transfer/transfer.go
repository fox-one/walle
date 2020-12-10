package transfer

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/walle/core"
)

type Config struct {
	Pin       string   `valid:"required"`
	Members   []string `valid:"required"`
	Threshold uint8    `valid:"required"`
}

func New(client *mixin.Client, cfg Config) core.TransferService {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	return &transferService{
		client:    client,
		pin:       cfg.Pin,
		members:   cfg.Members,
		threshold: cfg.Threshold,
	}
}

type transferService struct {
	client    *mixin.Client
	pin       string
	members   []string
	threshold uint8
}

func (s *transferService) Handle(ctx context.Context, transfer *core.Transfer) error {
	input := &mixin.TransferInput{
		AssetID:    transfer.AssetID,
		Amount:     transfer.Amount,
		TraceID:    transfer.TraceID,
		Memo:       transfer.Memo,
		OpponentID: transfer.OpponentID,
	}

	if input.OpponentID != "" {
		_, err := s.client.Transfer(ctx, input, s.pin)
		return err
	}

	input.OpponentMultisig.Receivers = s.members
	input.OpponentMultisig.Threshold = s.threshold
	_, err := s.client.Transaction(ctx, input, s.pin)
	return err
}
