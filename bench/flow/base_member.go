package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/model"
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

		var (
			member *model.MemberWithLending
			err    error
		)
		for i := 0; i < 10; i++ {
			member, err = c.mr.GetMemberByID(memberID)
			if err == nil {
				break
			}
			mem, err := c.mr.GetInactiveMemberID(1)
			if err != nil {
				return
			}
			memberID = mem[0]
		}

		timer := time.After(memberFlowCycle)

		choices := []utils.Choice[flow]{
			{Val: c.searchBooksFlow(step), Weight: 3},
			{Val: c.patchMemberFlow(memberID, step)},
		}
		if member.Lending {
			choices = append(choices, utils.Choice[flow]{
				Val: c.returnLendingsFlow(memberID, step),
			})
		} else {
			choices = append(choices, utils.Choice[flow]{
				Val: c.postLendingFlow(memberID, 2, step), Weight: 3,
			})
		}

		for {
			// memberがBANされてしまったら、強制的にフローを終了する
			member, err = c.mr.GetMemberByID(memberID)
			if err != nil {
				return
			}

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

		select {
		case <-ctx.Done():
			return
		case <-timer:
			return
		}
	}
}
