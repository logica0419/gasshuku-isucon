package flow

import (
	"context"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/model"
)

const (
	initialLibWorker = 7
	initialMemWorker = 9
)

// ワーカーの初期起動用ワーカー
func (c *Controller) StartUpFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for i := 0; i < initialLibWorker; i++ {
			w := c.baseLibraryFlow(step)
			c.wc <- w
		}

		mem, err := c.mr.GetInactiveMemberID(initialMemWorker)
		if err != nil {
			step.AddError(failure.NewError(model.ErrCritical, err))
			return
		}
		for _, id := range mem {
			w := c.baseMemberFlow(id, step)
			c.wc <- w
		}
	}
}
