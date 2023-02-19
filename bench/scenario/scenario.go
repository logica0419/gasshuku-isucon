package scenario

import (
	"github.com/isucon/isucandar"
	"github.com/logica0419/gasshuku-isucon/bench/flow"
)

type Scenario struct {
	wc <-chan flow.AddWorkerRequest
	fc *flow.FlowController
}

var _ isucandar.PrepareScenario = &Scenario{}

func NewScenario(wc <-chan flow.AddWorkerRequest, fc *flow.FlowController) *Scenario {
	return &Scenario{
		wc: wc,
		fc: fc,
	}
}
