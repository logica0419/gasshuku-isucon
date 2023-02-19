package flow

import "github.com/logica0419/gasshuku-isucon/bench/action"

type AddWorkerRequest struct {
	Flow string
	With any
}

type FlowController struct {
	wc chan<- AddWorkerRequest

	ma action.MemberActionController
}

func NewFlowController(
	c chan<- AddWorkerRequest,
	ma action.MemberActionController,
) *FlowController {
	return &FlowController{
		ma: ma,
	}
}
