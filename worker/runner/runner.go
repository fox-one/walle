package runner

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/pkg/property"
	"github.com/fox-one/pkg/uuid"
	"github.com/fox-one/walle/core"
)

const (
	checkpoint = "4swap_mtg_agent_orders_checkpoint"
	limit      = 500
)

func New(
	orders core.OrderStore,
	orderz core.OrderService,
	transfers core.TransferStore,
	property property.Store,
) *Runner {
	return &Runner{
		orders:    orders,
		orderz:    orderz,
		transfers: transfers,
		property:  property,
	}
}

type Runner struct {
	orders    core.OrderStore
	orderz    core.OrderService
	transfers core.TransferStore
	property  property.Store
}

func (w *Runner) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "runner")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 200 * time.Millisecond
			} else {
				dur = 500 * time.Millisecond
			}
		}
	}
}

func (w *Runner) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	v, err := w.property.Get(ctx, checkpoint)
	if err != nil {
		log.WithError(err).Errorln("property.Get", checkpoint)
		return err
	}

	offset := v.Time()
	if offset.IsZero() {
		offset = time.Now()
	}

	orders, err := w.orderz.Pull(ctx, offset, limit)
	if err != nil {
		log.WithError(err).Errorln("orderz.Pull")
		return err
	}

	if len(orders) == 0 {
		return errors.New("end of list")
	}

	for _, order := range orders {
		offset = order.CreatedAt

		t := &core.Transfer{
			BrokerID: order.BrokerID,
			TraceID:  order.TraceID,
			AssetID:  order.AssetID,
			Memo:     order.MtgMemo,
			Amount:   order.Amount,
		}

		switch order.State {
		case core.OrderStateRejected:
			t.OpponentID = order.UserID
			t.Memo = "order rejected"
		case core.OrderStateApproved:
		default:
			continue
		}

		t.TraceID = uuid.Modify(t.TraceID, "forwarding")
		if err := w.orders.Create(ctx, order); err != nil {
			log.WithError(err).Errorln("orders.Create")
			return err
		}

		if err := w.transfers.Create(ctx, t); err != nil {
			log.WithError(err).Errorln("transfers.Create")
			return err
		}
	}

	if err := w.property.Save(ctx, checkpoint, offset); err != nil {
		log.WithError(err).Errorln("property.Save", checkpoint)
		return err
	}

	return nil
}
