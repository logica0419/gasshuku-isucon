package flow

import (
	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type FlowController struct {
	wc chan<- worker.WorkerFunc

	key string
	cr  *utils.Crypt

	ia action.InitializeActionController
	ma action.MemberActionController
}

func NewFlowController(
	c chan<- worker.WorkerFunc,
	ia action.InitializeActionController,
	ma action.MemberActionController,
) (*FlowController, error) {
	key := utils.RandStringWithSign(16)
	cr, err := utils.NewCrypt(key)
	if err != nil {
		return nil, err
	}

	return &FlowController{
		wc:  c,
		key: key,
		cr:  cr,
		ia:  ia,
		ma:  ma,
	}, nil
}
