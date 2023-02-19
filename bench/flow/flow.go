package flow

import (
	"context"

	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/repository"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type flow func(ctx context.Context)

type FlowController struct {
	wc chan<- worker.WorkerFunc

	key string
	cr  *utils.Crypt

	ia action.InitializeActionController
	ma action.MemberActionController

	mr repository.MemberRepository
}

func NewFlowController(
	c chan<- worker.WorkerFunc,
	ia action.InitializeActionController,
	ma action.MemberActionController,
	mr repository.MemberRepository,
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
		mr:  mr,
	}, nil
}
