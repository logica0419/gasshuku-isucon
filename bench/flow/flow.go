package flow

import (
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type AddWorkerRequest struct {
	Flow string
	With any
}

type FlowController struct {
	wc chan<- AddWorkerRequest

	key string
	cr  *utils.Crypt

	ia action.InitializeActionController
	ma action.MemberActionController
}

func NewFlowController(
	c chan<- AddWorkerRequest,
	ia action.InitializeActionController,
	ma action.MemberActionController,
) (*FlowController, error) {
	key := utils.RandStringWithSign(16)
	cr, err := utils.NewCrypt(key)
	if err != nil {
		return nil, err
	}

	return &FlowController{
		wc: c,
		key: key,
		cr: cr,
		ia: ia,
		ma: ma,
	}, nil
}
