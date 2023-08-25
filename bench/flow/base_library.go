package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/logger"
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

		choices := []utils.Choice[flow]{
			{Val: c.getMembersFlow("", step), Weight: 10},
			{Val: c.getMembersFlow(c.mr.GetRandomMember().ID, step), Weight: 10},
			{Val: c.searchBooksFlow(step), Weight: 10},
			{Val: c.postBooksFlow(2, step), Weight: 30},
			{Val: c.getLendingsFlow(step), Weight: 30},
		}

		for {
			runner, err := utils.WeightedSelect(choices, true)
			if err != nil {
				break
			}
			runner(ctx)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		logger.Admin.Print("finish library cycle")

		select {
		case <-ctx.Done():
			return
		case <-timer:
			return
		}
	}
}
