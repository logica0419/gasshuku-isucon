// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package benchmark

import (
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/config"
	"github.com/logica0419/gasshuku-isucon/bench/flow"
	"github.com/logica0419/gasshuku-isucon/bench/repository"
	"github.com/logica0419/gasshuku-isucon/bench/scenario"
)

// Injectors from wire.go:

func NewBenchmark(c chan worker.WorkerFunc) (*Benchmark, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	controller, err := action.NewController(configConfig)
	if err != nil {
		return nil, err
	}
	repositoryRepository, err := repository.NewRepository()
	if err != nil {
		return nil, err
	}
	flowController, err := flow.NewController(c, controller, controller, controller, repositoryRepository, repositoryRepository)
	if err != nil {
		return nil, err
	}
	scenarioScenario := scenario.NewScenario(c, flowController)
	benchmark, err := newBenchmark(configConfig, scenarioScenario)
	if err != nil {
		return nil, err
	}
	return benchmark, nil
}