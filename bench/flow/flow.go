package flow

import (
	"context"

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
		wc:  wc,
		sc:  sc,
		key: key,
		cr:  cr,
		ia:  ia,
		ma:  ma,
		ba:  ba,
		la:  la,
		mr:  mr,
		br:  br,
		lr:  lr,
	}, nil
}
