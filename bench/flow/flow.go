package flow

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/isucon/isucandar/worker"
	"github.com/logica0419/gasshuku-isucon/bench/action"
	"github.com/logica0419/gasshuku-isucon/bench/repository"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

type flow func(ctx context.Context)

type Controller struct {
	wc chan<- worker.WorkerFunc
	sc chan<- struct{}

	key string
	cr  *utils.Crypt

	memInCycleCount uint32
	libInCycleCount uint32

	activeMemWorkerCount uint32
	activeLibWorkerCount uint32
	addedWorkerHistory   []string
	historyLock          sync.Mutex

	ia action.InitializeController
	ma action.MemberController
	ba action.BookController
	la action.LendingController

	mr repository.MemberRepository
	br repository.BookRepository
	lr repository.LendingRepository
}

func NewController(
	wc chan worker.WorkerFunc,
	sc chan struct{},
	ia action.InitializeController,
	ma action.MemberController,
	ba action.BookController,
	la action.LendingController,
	mr repository.MemberRepository,
	br repository.BookRepository,
	lr repository.LendingRepository,
) (*Controller, error) {
	key := utils.RandStringWithSign(16)
	cr, err := utils.NewCrypt(key)
	if err != nil {
		return nil, err
	}

	return &Controller{
		wc:                 wc,
		sc:                 sc,
		key:                key,
		cr:                 cr,
		addedWorkerHistory: []string{},
		historyLock:        sync.Mutex{},
		ia:                 ia,
		ma:                 ma,
		ba:                 ba,
		la:                 la,
		mr:                 mr,
		br:                 br,
		lr:                 lr,
	}, nil
}

func (c *Controller) addActiveMemWorkerCount() {
	c.historyLock.Lock()
	defer c.historyLock.Unlock()
	c.addedWorkerHistory = append(c.addedWorkerHistory, "mem")
	atomic.AddUint32(&c.activeMemWorkerCount, 1)
}

func (c *Controller) decActiveMemWorkerCount() {
	c.historyLock.Lock()
	defer c.historyLock.Unlock()
	c.addedWorkerHistory = c.addedWorkerHistory[:len(c.addedWorkerHistory)-1]
	atomic.AddUint32(&c.activeMemWorkerCount, ^uint32(0))
}

func (c *Controller) addActiveLibWorkerCount() {
	c.historyLock.Lock()
	defer c.historyLock.Unlock()
	c.addedWorkerHistory = append(c.addedWorkerHistory, "lib")
	atomic.AddUint32(&c.activeLibWorkerCount, 1)
}

func (c *Controller) decActiveLibWorkerCount() {
	c.historyLock.Lock()
	defer c.historyLock.Unlock()
	c.addedWorkerHistory = c.addedWorkerHistory[:len(c.addedWorkerHistory)-1]
	atomic.AddUint32(&c.activeLibWorkerCount, ^uint32(0))
}
