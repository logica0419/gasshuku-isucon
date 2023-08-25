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

	cancelFuncs := []context.CancelFunc{}

	for {
		select {
		case <-ctx.Done():
			return nil

		case wf := <-s.wc:
			nw, err := worker.NewWorker(wf, worker.WithInfinityLoop(), worker.WithMaxParallelism(1))
			if err != nil {
				return err
			}
			ctx, cancel := context.WithCancel(ctx)
			cancelFuncs = append(cancelFuncs, cancel)
			go nw.Process(ctx)

		case <-s.sc:
			if len(cancelFuncs) == 0 {
				continue
			}
			cancelFuncs[len(cancelFuncs)-1]()
			cancelFuncs = cancelFuncs[:len(cancelFuncs)-1]
		}
	}
}
