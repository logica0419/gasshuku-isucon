package benchmark

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
)

const printPeriod = 3 * time.Second

func scorePrinter(ctx context.Context, step *isucandar.BenchmarkStep) error {
	for {
		ticker := time.NewTicker(printPeriod)

		select {
		case <-ticker.C:
			_ = grader.CulcResult(step.Result(), false)
		case <-ctx.Done():
			return nil
		}
	}
}
