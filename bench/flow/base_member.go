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

		member, err := c.mr.GetMemberByID(memberID)
		if err != nil {
			timer := time.After(memberFlowCycle)
			select {
			case <-ctx.Done():
			case <-timer:
				return
			}
		}

		timer := time.After(memberFlowCycle)

		choices := []utils.Choice[flow]{
			{Val: c.searchBooksFlow(step)},
			{Val: c.patchMemberFlow(memberID, step)},
		}
		if member.Lending {
			choices = append(choices, utils.Choice[flow]{
				Val:    c.returnLendingsFlow(memberID, step),
				Weight: 2,
			})
		} else {
			choices = append(choices, utils.Choice[flow]{
				Val:    c.postLendingFlow(memberID, int(c.activeMemWorkerCount*2), step),
				Weight: 2,
			})
		}

		runner := utils.WeightedSelect(choices)
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
