package order

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/4swap-sdk-go/mtg"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/number"
	"github.com/fox-one/walle/core"
)

type Config struct {
	VerifyKey string `valid:"required"`
}

func New(client *mixin.Client, cfg Config) core.OrderService {
	if _, err := govalidator.ValidateStruct(cfg); err != nil {
		panic(err)
	}

	b, err := base64.StdEncoding.DecodeString(cfg.VerifyKey)
	if err != nil {
		panic(err)
	}

	return &orderService{
		client:    client,
		verifyKey: b,
	}
}

type orderService struct {
	client    *mixin.Client
	verifyKey ed25519.PublicKey
}

func (s *orderService) Pull(ctx context.Context, checkpoint time.Time, limit int) ([]*core.Order, error) {
	snapshots, err := s.client.ReadNetworkSnapshots(ctx, "", checkpoint, "ASC", limit)
	if err != nil {
		return nil, err
	}

	orders := make([]*core.Order, 0, len(snapshots))
	for _, snapshot := range snapshots {
		orders = append(orders, convertSnapshot(snapshot, s.verifyKey))
	}

	return orders, nil
}

func convertSnapshot(snapshot *mixin.Snapshot, key ed25519.PublicKey) *core.Order {
	order := &core.Order{
		CreatedAt: snapshot.CreatedAt,
		UserID:    snapshot.OpponentID,
		AssetID:   snapshot.AssetID,
		TraceID:   snapshot.TraceID,
		BrokerID:  snapshot.UserID,
		Amount:    snapshot.Amount,
		Memo:      snapshot.Memo,
	}

	if order.UserID == "" || !order.Amount.IsPositive() {
		return order
	}

	action, err := DecodeTransactionAction(order.Memo)
	if err != nil {
		order.State = core.OrderStateRejected
		return order
	}

	var mtgAction mtg.Action

	switch ParseTransactionType(action.Type) {
	case TransactionTypeAdd:
		mtgAction = mtg.AddAction(order.UserID, action.Deposit, action.AssetID, 10*time.Minute, number.Decimal("0.005"))
	case TransactionTypeRemove:
		mtgAction = mtg.RemoveAction(order.UserID, order.TraceID)
	case TransactionTypeSwap:
		mtgAction = mtg.SwapAction(order.UserID, order.TraceID, action.AssetID, action.Routes, number.Decimal(action.Minimum))
	}

	memo, err := mtg.EncodeAction(mtgAction, key)
	if err != nil {
		order.State = core.OrderStateRejected
		return order
	}

	order.State = core.OrderStateApproved
	order.MtgMemo = memo
	return order
}
