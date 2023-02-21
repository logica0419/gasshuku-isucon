package benchmark

import (
	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/logger"
	"github.com/logica0419/gasshuku-isucon/bench/model"
)

func registerErrorHandler(b *Benchmark) {
	b.ib.OnError(func(err error, step *isucandar.BenchmarkStep) {
		if model.IsErrCanceled(err) {
			return
		}

		if model.IsErrCritical(err) {
			logger.Contestant.Printf("critical error - %v", err)
			logger.Contestant.Print("--------- stop benchmarking ---------")
			step.Cancel()
			return
		}

		logger.Contestant.Printf("error - %v", err)
	})
}
