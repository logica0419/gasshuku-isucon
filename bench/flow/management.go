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

		for i := 0; i < 5; i++ {
			w := c.baseLibraryFlow(step)
			c.wc <- w
			c.addActiveLibWorkerCount()
		}
	}
}

const checkerCycle = 10 * time.Millisecond

// ワーカーの追加 / 停止ワーカー
func (c *Controller) ScalingFlow(step *isucandar.BenchmarkStep) worker.WorkerFunc {
	prevTimeoutCount := 0

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
				// タイムアウトが発生していたら、スケールダウン
				timeoutCount := 0
				for _, err := range step.Result().Errors.All() {
					if model.IsErrTimeout(err) {
						timeoutCount++
					}
				}
				if timeoutCount > prevTimeoutCount {
					prevTimeoutCount = timeoutCount
					switch c.addedWorkerHistory[len(c.addedWorkerHistory)-1] {
					case "lib":
						c.decActiveLibWorkerCount()
						logger.Admin.Print("タイムアウトが発生したため、図書館職員ワーカーを1つ停止しました")
					case "mem":
						c.decActiveMemWorkerCount()
						logger.Admin.Print("タイムアウトが発生したため、会員ワーカーを1つ停止しました")
					}
					c.sc <- struct{}{}
					break
				}

				// 図書館職員フローが時間内に9/10終了かつ会員フローが時間内に1/5終了したら、図書館職員フローを追加
				if c.libInCycleCount > c.activeLibWorkerCount*9/10 && c.memInCycleCount > c.activeMemWorkerCount/5 {
					join := int(c.activeLibWorkerCount / 5)
					for i := 0; i < join; i++ {
						w := c.baseLibraryFlow(step)
						c.wc <- w
						c.addActiveLibWorkerCount()
					}
					c.resetLibInCycleCount()
					logger.Contestant.Printf("追加で%d個の図書館職員ワーカーが開始されました", join)
				}

				// 会員フローが時間内に9/10終了かつ図書館職員フローが時間内に1/5終了したら、会員フローを追加
				if c.memInCycleCount > c.activeMemWorkerCount*9/5 && c.libInCycleCount > c.activeLibWorkerCount/5 {
					join := int(c.activeMemWorkerCount / 5)
					mem, err := c.mr.GetInactiveMemberID(join)
					if err != nil {
						c.postMemberFlow(join, step)(ctx)
						mem, err = c.mr.GetInactiveMemberID(join)
						if err != nil {
							step.AddError(failure.NewError(model.ErrCritical, err))
							return
						}
					}
					for _, id := range mem {
						w := c.baseMemberFlow(id, step)
						c.wc <- w
						c.addActiveMemWorkerCount()
					}
					c.resetMemInCycleCount()
					logger.Admin.Printf("追加で%d個の会員ワーカーが開始されました", join)
				}
			}
		}
	}
}
