package transfer

import (
	"context"

	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/core"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Transfer{})

		if err := tx.AutoMigrate(core.Transfer{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_transfers_trace", "trace_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.TransferStore {
	return &transferStore{db: db}
}

type transferStore struct {
	db *db.DB
}

func (s *transferStore) Create(ctx context.Context, transfer *core.Transfer) error {
	return s.db.Update().Where("trace_id = ?", transfer.TraceID).FirstOrCreate(transfer).Error
}

func (s *transferStore) Delete(ctx context.Context, transfers []*core.Transfer) error {
	ids := make([]int64, len(transfers))
	for idx, t := range transfers {
		ids[idx] = t.ID
	}

	if len(ids) == 0 {
		return nil
	}

	if err := s.db.Update().Where("id IN (?)", ids).Delete(core.Transfer{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *transferStore) List(ctx context.Context, limit int) ([]*core.Transfer, error) {
	var transfers []*core.Transfer
	if err := s.db.View().Limit(limit).Find(&transfers).Error; err != nil {
		return nil, err
	}

	return transfers, nil
}
