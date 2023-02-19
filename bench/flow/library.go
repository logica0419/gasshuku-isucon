package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

const libraryFlowCycle = 50 * time.Millisecond

// 図書館職員フロー
func (c *FlowController) LibraryFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		timer := time.After(libraryFlowCycle)

		runner := utils.WeightedSelect(
			[]utils.Choice[flow]{
				{Val: c.membersGetFlow(step)},
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

		c.addLibInCycleCount()

		select {
		case <-ctx.Done():
			return
		case <-timer:
			return
		}
	}
}
