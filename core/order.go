package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type OrderState int

const (
	_ OrderState = iota
	OrderStateRejected
	OrderStateApproved
)

type (
	Order struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		State     OrderState      `json:"state,omitempty"`
		UserID    string          `sql:"size:36" json:"user_id,omitempty"`
		AssetID   string          `sql:"size:36" json:"asset_id,omitempty"`
		TraceID   string          `sql:"size:36" json:"trace_id,omitempty"`
		BrokerID  string          `sql:"size:36" json:"broker_id,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Memo      string          `sql:"size:140" json:"memo,omitempty"`
		MtgMemo   string          `sql:"size:200" json:"mtg_memo,omitempty"`
	}

	OrderStore interface {
		Create(ctx context.Context, order *Order) error
	}

	OrderService interface {
		Pull(ctx context.Context, checkpoint time.Time, limit int) ([]*Order, error)
	}
)
