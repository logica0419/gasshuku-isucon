package scenario

import (
	"context"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
)

func (s *Scenario) Load(ctx context.Context, step *isucandar.BenchmarkStep) error {
	ctx, cancel := context.WithTimeout(ctx, benchTime)
	defer cancel()

	worker, err := worker.NewWorker(s.fc.MemberGetFlow(step), worker.WithInfinityLoop(), worker.WithMaxParallelism(1))
	if err != nil {
		return err
	}
	worker.Process(ctx)

	return nil
}
