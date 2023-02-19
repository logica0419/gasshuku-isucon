package benchmark

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/scenario"
)

type Benchmark struct {
	ib *isucandar.Benchmark
}

func newBenchmark(s *scenario.Scenario) (*Benchmark, error) {
	ib, err := isucandar.NewBenchmark(
		isucandar.WithLoadTimeout(scenario.BenchTime + time.Second*10),
	)
	if err != nil {
		return nil, err
	}

	ib.AddScenario(s)

	b := &Benchmark{
		ib: ib,
	}

	registerErrorHandler(b)

	return b, nil
}

func (b *Benchmark) Run(ctx context.Context) {
	_ = b.ib.Start(ctx)
}
