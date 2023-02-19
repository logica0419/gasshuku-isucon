package main

import (
	"context"
	"log"

	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/benchmark"
	"github.com/logica0419/gasshuku-isucon/bench/flow"
	"github.com/logica0419/gasshuku-isucon/bench/scenario"
)

func main() {
	a, err := action.NewActionController(3, 3, "http://localhost:8080")
	if err != nil {
		log.Panic(err)
	}

	c := make(chan worker.WorkerFunc, 100)

	f, err := flow.NewFlowController(c, a, a)
	if err != nil {
		log.Panic(err)
	}

	s := scenario.NewScenario(c, f)

	b, err := benchmark.NewBenchmark(s)
	if err != nil {
		log.Panic(err)
	}

	ctx := context.Background()
	b.Run(ctx)
}
