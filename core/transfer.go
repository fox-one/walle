package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Transfer struct {
		ID        int64           `sql:"PRIMARY_KEY" json:"id,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
		BrokerID  string          `sql:"size:36" json:"broker_id,omitempty"`
		TraceID   string          `sql:"size:36" json:"trace_id,omitempty"`
		AssetID   string          `sql:"size:36" json:"asset_id,omitempty"`
		Memo      string          `sql:"size:200" json:"memo,omitempty"`
		Amount    decimal.Decimal `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		// 如果是空的，则是转账给多签节点
		OpponentID string `sql:"size:36" json:"opponent_id,omitempty"`
	}

	TransferStore interface {
		Create(ctx context.Context, transfer *Transfer) error
		Delete(ctx context.Context, transfers []*Transfer) error
		List(ctx context.Context, limit int) ([]*Transfer, error)
	}

	TransferService interface {
		Handle(ctx context.Context, transfer *Transfer) error
	}
)
