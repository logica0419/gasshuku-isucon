package flow

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/logger"
)

// ワーカーの初期起動用ワーカー
func (c *FlowController) StartUpFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	return func(ctx context.Context, _ int) {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mem := c.mr.GetInactiveMemberID(30)
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
func (c *FlowController) ScaleUpFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
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
				if c.libInCycleCount > c.activeLibWorkerCount*4/5 && c.memInCycleCount > c.activeMemWorkerCount/5 {
					join := int(c.activeLibWorkerCount / 5)
					for i := 0; i < join; i++ {
						w := c.baseLibraryFlow(step)
						c.wc <- w
						c.addActiveLibWorkerCount()
					}
					logger.Admin.Printf("%d人の図書館職員が採用されました: 計%d", join, c.activeLibWorkerCount)
				}

				if c.memInCycleCount > c.activeMemWorkerCount*4/5 && c.libInCycleCount > c.activeLibWorkerCount/5 {
					join := int(c.activeMemWorkerCount / 5)
					mem := c.mr.GetInactiveMemberID(join)
					for _, id := range mem {
						w := c.baseMemberFlow(id, step)
						c.wc <- w
						c.addActiveMemWorkerCount()
					}
					logger.Admin.Printf("%d人の会員が新規に登録しました: 計%d人", join, c.activeMemWorkerCount)
				}
			}
		}
	}
}