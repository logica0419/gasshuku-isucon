//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/benchmark"
	"github.com/logica0419/gasshuku-isucon/bench/config"
	"github.com/logica0419/gasshuku-isucon/bench/flow"
	"github.com/logica0419/gasshuku-isucon/bench/repository"
	"github.com/logica0419/gasshuku-isucon/bench/scenario"
)

func injectBenchmark(c chan worker.WorkerFunc) (*benchmark.Benchmark, error) {
	wire.Build(
		config.NewConfig,

		repository.NewRepository,
		wire.Bind(new(repository.MemberRepository), new(*repository.Repository)),

		action.NewActionController,
		wire.Bind(new(action.MemberActionController), new(*action.ActionController)),
		wire.Bind(new(action.InitializeActionController), new(*action.ActionController)),

		flow.NewFlowController,

		scenario.NewScenario,

		benchmark.NewBenchmark,
	)

	return nil, nil
}
