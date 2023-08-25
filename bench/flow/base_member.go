package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/logger"
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
			{Val: c.searchBooksFlow(step), Weight: 20},
			{Val: c.patchMemberFlow(memberID, step), Weight: 10},
			{Val: func(ctx context.Context) {
				member, _ = c.mr.GetMemberByID(memberID)
				if member.Lending {
					c.returnLendingsFlow(memberID, step)(ctx)
				} else {
					c.postLendingFlow(memberID, 2, step)(ctx)
				}
			}, Weight: 30},
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

		logger.Admin.Print("finish member cycle")

		select {
		case <-ctx.Done():
			return
		case <-timer:
			return
		}
	}
}
