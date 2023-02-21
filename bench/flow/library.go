package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

const libraryFlowCycle = 500 * time.Millisecond

// 図書館職員フロー
func (c *FlowController) libraryFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		timer := time.After(libraryFlowCycle)

		runner := utils.WeightedSelect(
			[]utils.Choice[flow]{
				{Val: c.membersGetFlow("", step)},
				{Val: c.membersGetFlow(utils.RandString(26), step), Weight: 2},
				{Val: c.membersGetFlow(c.mr.GetRandomMember().ID, step), Weight: 2},
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
