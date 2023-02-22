package scenario

import (
	"context"
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
)

const BenchTime = 60 * time.Second

func (s *Scenario) Load(ctx context.Context, step *isucandar.BenchmarkStep) error {
	ctx, cancel := context.WithTimeout(ctx, BenchTime)
	defer cancel()

	stw, err := worker.NewWorker(s.fc.StartUpFlow(step), worker.WithLoopCount(1))
	if err != nil {
		return err
	}
	go stw.Process(ctx)

	scw, err := worker.NewWorker(s.fc.ScaleUpFlow(step), worker.WithInfinityLoop(), worker.WithMaxParallelism(1))
	if err != nil {
		return err
	}
	go scw.Process(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case wf := <-s.wc:
			nw, err := worker.NewWorker(wf, worker.WithInfinityLoop(), worker.WithMaxParallelism(1))
			if err != nil {
				return err
			}
			go nw.Process(ctx)
		}
	}
}
