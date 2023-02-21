package main

import (
	"context"

	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/logger"
)

func main() {
	b, err := injectBenchmark(make(chan worker.WorkerFunc, 100))
	if err != nil {
		logger.Admin.Panic(err)
	}

	b.Run(context.Background())
}
