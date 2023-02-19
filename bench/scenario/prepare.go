package scenario

import (
	"context"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
)

func (s *Scenario) Prepare(ctx context.Context, step *isucandar.BenchmarkStep) error {
	initWorker, err := worker.NewWorker(s.fc.InitializeFlow(step), worker.WithLoopCount(1))
	if err != nil {
		return err
	}

	initWorker.Process(ctx)
	return nil
}
