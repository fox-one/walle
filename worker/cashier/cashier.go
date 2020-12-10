package cashier

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/walle/core"
	"golang.org/x/sync/errgroup"
)

func New(
	brokers core.BrokerStore,
	transfers core.TransferStore,
	transferz func(broker *core.Broker) core.TransferService,
) *Cashier {
	return &Cashier{
		brokers:   brokers,
		transfers: transfers,
		transferz: transferz,
	}
}

type Cashier struct {
	brokers   core.BrokerStore
	transfers core.TransferStore
	transferz func(broker *core.Broker) core.TransferService
}

func (w *Cashier) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "cashier")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			if err := w.run(ctx); err == nil {
				dur = 100 * time.Millisecond
			} else {
				dur = time.Second
			}
		}
	}
}

func (w *Cashier) run(ctx context.Context) error {
	log := logger.FromContext(ctx)

	const Limit = 500
	transfers, err := w.transfers.List(ctx, Limit)
	if err != nil {
		log.WithError(err).Errorln("list transfers")
		return err
	}

	if len(transfers) == 0 {
		return errors.New("end of list")
	}

	transfers = groupTransfers(transfers, 10)
	handled := make([]int, len(transfers))

	var g errgroup.Group
	for idx, t := range transfers {
		idx, t := idx, t
		g.Go(func() error {
			broker, err := w.brokers.Find(ctx, t.BrokerID)
			if err != nil {
				log.WithError(err).Errorln("brokers.Find", t.BrokerID)
				return err
			}

			if err := w.transferz(broker).Handle(ctx, t); err != nil {
				log.WithError(err).Errorln("transferz.Handle")
				return err
			}

			handled[idx] = 1
			return nil
		})
	}

	_ = g.Wait()

	var idx int
	for i, transfer := range transfers {
		if handled[i] == 0 {
			continue
		}

		transfers[idx] = transfer
		idx++
	}

	if done := transfers[:idx]; len(done) > 0 {
		if err := w.transfers.Delete(ctx, done); err != nil {
			log.WithError(err).Errorln("wallets.DeleteTransfers")
			return err
		}
	}

	return nil
}

func groupTransfers(transfers []*core.Transfer, max int) []*core.Transfer {
	var (
		idx     int
		brokers = make(map[string]bool, max)
	)

	for _, t := range transfers {
		if brokers[t.BrokerID] {
			continue
		}

		brokers[t.BrokerID] = true
		transfers[idx] = t
		idx++

		if idx >= max {
			break
		}
	}

	return transfers[:idx]
}
