package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

const memberFlowCycle = 500 * time.Millisecond

// 会員フロー
func (c *Controller) baseMemberFlow(memberID string, step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		timer := time.After(memberFlowCycle)

		runner := utils.WeightedSelect(
			[]utils.Choice[flow]{
				{Val: c.searchBooksFlow(step)},
				{Val: c.lendingsPostFlow(int(c.activeMemWorkerCount*2), step), Weight: 3},
			},
		)
		runner(ctx)

		select {
		case <-ctx.Done():
			return
		case <-timer:
			return
		default:
		}

		c.addMemInCycleCount()

		select {
		case <-ctx.Done():
			return
		case <-timer:
			return
		}
	}
}
