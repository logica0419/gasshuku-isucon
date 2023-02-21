package benchmark

import (
	"context"
	"os"
	"time"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/config"
	"github.com/logica0419/gasshuku-isucon/bench/grader"
	"github.com/logica0419/gasshuku-isucon/bench/scenario"
)

type Benchmark struct {
	ib         *isucandar.Benchmark
	exitStatus bool
}

func newBenchmark(c *config.Config, s *scenario.Scenario) (*Benchmark, error) {
	ib, err := isucandar.NewBenchmark(
		isucandar.WithLoadTimeout(scenario.BenchTime + time.Second*1),
	)
	if err != nil {
		return nil, err
	}

	ib.AddScenario(s)

	b := &Benchmark{
		ib:         ib,
		exitStatus: c.ExitStatusOnFail,
	}

	registerErrorHandler(b)
	registerScorePrinter(b)

	return b, nil
}

func (b *Benchmark) Run(ctx context.Context) {
	res := b.ib.Start(ctx)

	if !grader.CalcResult(res, true) && b.exitStatus {
		os.Exit(1)
	}
}
