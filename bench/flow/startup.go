package flow

import (
	"context"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
)

// ワーカーの初期起動用ワーカー
func (c *FlowController) StartUpFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mem := c.mr.GetInactiveMemberID(10)
		for _, id := range mem {
			w := c.baseMemberFlow(id, step)
			c.wc <- w
		}

		for i := 0; i < 10; i++ {
			w := c.baseLibraryFlow(step)
			c.wc <- w
		}
	}
}
