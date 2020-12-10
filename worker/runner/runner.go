package runner

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
)

type Runner struct {
}

func (r *Runner) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "runner")
	ctx = logger.WithContext(ctx, log)

	dur := time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(dur):
			// do somthing
		}
	}
}
