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

func NewBenchmark(s *scenario.Scenario) (*Benchmark, error) {
	b, err := isucandar.NewBenchmark(
		isucandar.WithLoadTimeout(scenario.BenchTime + time.Second*10),
	)
	if err != nil {
		return nil, err
	}

	b.AddScenario(s)

	return &Benchmark{
		ib: b,
	}, nil
}

func (b *Benchmark) Run(ctx context.Context) {
	_ = b.ib.Start(ctx)
}
