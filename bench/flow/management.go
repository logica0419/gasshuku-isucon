package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/failure"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/logger"
	"github.com/logica0419/gasshuku-isucon/bench/model"
)

// ワーカーの初期起動用ワーカー
func (c *Controller) StartUpFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mem, err := c.mr.GetInactiveMemberID(10)
		if err != nil {
			step.AddError(failure.NewError(model.ErrCritical, err))
			return
		}
		for _, id := range mem {
			w := c.baseMemberFlow(id, step)
			c.wc <- w
			c.addActiveMemWorkerCount()
		}

		for i := 0; i < 10; i++ {
			w := c.baseLibraryFlow(step)
			c.wc <- w
			c.addActiveLibWorkerCount()
		}
	}
}

const checkerCycle = 10 * time.Millisecond

// ワーカーの追加ワーカー
func (c *Controller) ScaleUpFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		baseTicker := time.NewTicker(checkerCycle)
		memberTicker := time.NewTicker(memberFlowCycle)
		libraryTicker := time.NewTicker(libraryFlowCycle)

		for {
			select {
			case <-ctx.Done():
				return
			case <-memberTicker.C:
				c.resetMemInCycleCount()
			case <-libraryTicker.C:
				c.resetLibInCycleCount()

			case <-baseTicker.C:
				// 図書館職員フローが時間内に4/5終了かつ会員フローが時間内に1/5終了したら。図書館職員フローを追加
				if c.libInCycleCount > c.activeLibWorkerCount*4/5 && c.memInCycleCount > c.activeMemWorkerCount/5 {
					join := int(c.activeLibWorkerCount / 5)
					for i := 0; i < join; i++ {
						w := c.baseLibraryFlow(step)
						c.wc <- w
						c.addActiveLibWorkerCount()
					}
					logger.Contestant.Printf("%d人の図書館職員が採用されました: 計%d", join, c.activeLibWorkerCount)
				}

				if c.memInCycleCount > c.activeMemWorkerCount*4/5 && c.libInCycleCount > c.activeLibWorkerCount/5 {
					join := int(c.activeMemWorkerCount / 5)
					mem, err := c.mr.GetInactiveMemberID(join)
					if err != nil {
						c.memberPostFlow(join, step)(ctx)
						mem, err = c.mr.GetInactiveMemberID(join)
						if err != nil {
							step.AddError(failure.NewError(model.ErrCritical, err))
							return
						}
						logger.Contestant.Printf("%d人の会員が新規に登録しました", join)
					}
					for _, id := range mem {
						w := c.baseMemberFlow(id, step)
						c.wc <- w
						c.addActiveMemWorkerCount()
					}
					logger.Contestant.Printf("%d人の会員がアクティブになりました: 計%d人", join, c.activeMemWorkerCount)
				}
			}
		}
	}
}
