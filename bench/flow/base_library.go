package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

const libraryFlowCycle = 1000 * time.Millisecond

// 図書館職員フロー
func (c *Controller) baseLibraryFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		timer := time.After(libraryFlowCycle)

		runner := utils.WeightedSelect(
			[]utils.Choice[flow]{
				{Val: c.getMembersFlow("", step)},
				{Val: c.getMembersFlow(c.mr.GetRandomMember().ID, step)},
				{Val: c.searchBooksFlow(step), Weight: 3},
				{Val: c.postBooksFlow(int(c.activeLibWorkerCount)*2, step), Weight: 2},
				{Val: c.getLendingsFlow(step), Weight: 4},
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
