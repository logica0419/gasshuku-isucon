package repository

import (
	"sync"

	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

// go:embed init_data.json
var initData []byte

type InitData struct {
	Members []*model.MemberWithLending `json:"members"`
	Books   []*model.BookWithLending   `json:"books"`
}

type Repository struct {
	mLock       sync.RWMutex
	memberSlice []*model.MemberWithLending
	memberMap   map[string]*model.MemberWithLending
}

func NewRepository() (*Repository, error) {
	r := &Repository{
		mLock:       sync.RWMutex{},
		memberSlice: []*model.MemberWithLending{},
		memberMap:   map[string]*model.MemberWithLending{},
	}

	var data InitData
	if err := utils.ByteToStruct(initData, &data); err != nil {
		return nil, err
	}

	return r, nil
}
