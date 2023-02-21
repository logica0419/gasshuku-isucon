package scenario

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
)

const BenchTime = 10 * time.Second

func (s *Scenario) Load(ctx context.Context, step *isucandar.BenchmarkStep) error {
	ctx, cancel := context.WithTimeout(ctx, BenchTime)
	defer cancel()

	w, err := worker.NewWorker(s.fc.StartUpFlow(step), worker.WithLoopCount(1))
	if err != nil {
		return err
	}
	go w.Process(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case wf := <-s.wc:
			w, err = worker.NewWorker(wf, worker.WithInfinityLoop(), worker.WithMaxParallelism(1))
			if err != nil {
				return err
			}
			go w.Process(ctx)

		}
	}
}
