package flow

import "github.com/logica0419/gasshuku-isucon/bench/action"

type FlowController struct {
	ma action.MemberActionController
}

func NewFlowController(ma action.MemberActionController) *FlowController {
	return &FlowController{
		ma: ma,
	}
}
