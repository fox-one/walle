package order

import (
	"context"

	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/core"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Order{})

		if err := tx.AutoMigrate(core.Order{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_orders_trace", "trace_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.OrderStore {
	return &orderStore{db: db}
}

type orderStore struct {
	db *db.DB
}

func (s *orderStore) Create(ctx context.Context, order *core.Order) error {
	return s.db.Update().Where("trace_id = ?", order.TraceID).FirstOrCreate(order).Error
}
