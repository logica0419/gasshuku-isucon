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
			for {
				ticker := time.NewTicker(printPeriod)

				select {
				case <-ticker.C:
					_ = grader.CalcResult(step.Result(), false)
				case <-ctx.Done():
					return nil
				}
			}
		}
	})
}
