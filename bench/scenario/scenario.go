package scenario

import (
	"time"

	"github.com/isucon/isucandar"
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/flow"
)

const BenchTime = 10 * time.Second

type Scenario struct {
	wc <-chan worker.WorkerFunc
	fc *flow.FlowController
}

var (
	_ isucandar.PrepareScenario = &Scenario{}
	_ isucandar.LoadScenario    = &Scenario{}
)

func NewScenario(wc <-chan worker.WorkerFunc, fc *flow.FlowController) *Scenario {
	return &Scenario{
		wc: wc,
		fc: fc,
	}
}
