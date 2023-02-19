package scenario

import (
	"github.com/logica0419/gasshuku-isucon/bench/flow"
)

type Scenario struct {
	wc <-chan flow.AddWorkerRequest
	fc *flow.FlowController
}

func NewScenario(wc <-chan flow.AddWorkerRequest, fc *flow.FlowController) *Scenario {
	return &Scenario{
		wc: wc,
		fc: fc,
	}
}
