package flow

import (
	"context"
	"sync/atomic"

	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/repository"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type flow func(ctx context.Context)

type Controller struct {
	wc chan<- worker.WorkerFunc

	key string
	cr  *utils.Crypt

	memInCycleCount uint32
	libInCycleCount uint32

	activeMemWorkerCount uint32
	activeLibWorkerCount uint32

	ia action.InitializeController
	ma action.MemberController

	mr repository.MemberRepository
}

func NewController(
	c chan worker.WorkerFunc,
	ia action.InitializeController,
	ma action.MemberController,
	mr repository.MemberRepository,
) (*Controller, error) {
	key := utils.RandStringWithSign(16)
	cr, err := utils.NewCrypt(key)
	if err != nil {
		return nil, err
	}

	return &Controller{
		wc:              c,
		key:             key,
		cr:              cr,
		libInCycleCount: 0,
		ia:              ia,
		ma:              ma,
		mr:              mr,
	}, nil
}

func (c *Controller) addLibInCycleCount() {
	atomic.AddUint32(&c.libInCycleCount, 1)
}

func (c *Controller) resetLibInCycleCount() {
	atomic.StoreUint32(&c.libInCycleCount, 0)
}

func (c *Controller) addMemInCycleCount() {
	atomic.AddUint32(&c.memInCycleCount, 1)
}

func (c *Controller) resetMemInCycleCount() {
	atomic.StoreUint32(&c.memInCycleCount, 0)
}

func (c *Controller) addActiveMemWorkerCount() {
	atomic.AddUint32(&c.activeMemWorkerCount, 1)
}

func (c *Controller) addActiveLibWorkerCount() {
	atomic.AddUint32(&c.activeLibWorkerCount, 1)
}
