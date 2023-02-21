package benchmark

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
)

const printPeriod = 1 * time.Second

func registerScorePrinter(b *Benchmark) {
	b.ib.Load(func(ctx context.Context, step *isucandar.BenchmarkStep) error {
		for {
			ticker := time.NewTicker(printPeriod)

			select {
			case <-ticker.C:
				pass := grader.CalcResult(step.Result(), false)
				if !pass {
					step.Cancel()
				}
			case <-ctx.Done():
				return nil
			}
		}
	})
}
