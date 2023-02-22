package repository

import (
	_ "embed"
	"errors"
	"sync"

	"github.com/logica0419/gasshuku-isucon/bench/model"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
)

//go:embed init_data.json
var initData []byte

type InitData struct {
	Members []*model.MemberWithLending `json:"members"`
	Books   []*model.BookWithLending   `json:"books"`
}

var (
	ErrNotFound         = errors.New("not found")
	ErrNotEnoughRecords = errors.New("not enough records")
)

type Repository struct {
	mLock            sync.RWMutex
	memberSlice      []*model.MemberWithLending
	memberMap        map[string]*model.MemberWithLending
	inactiveMemberID []string

	bLock     sync.RWMutex
	bookSlice []*model.BookWithLending
	bookMap   map[string]*model.BookWithLending

	lLock            sync.RWMutex
	lendingMemberMap map[string][]*model.LendingWithNames
	lendingMap       map[string]*model.LendingWithNames
}

func NewRepository() (*Repository, error) {
	r := &Repository{
		mLock:            sync.RWMutex{},
		memberSlice:      []*model.MemberWithLending{},
		memberMap:        map[string]*model.MemberWithLending{},
		inactiveMemberID: []string{},

		bLock:     sync.RWMutex{},
		bookSlice: []*model.BookWithLending{},
		bookMap:   map[string]*model.BookWithLending{},

		lLock:            sync.RWMutex{},
		lendingMemberMap: map[string][]*model.LendingWithNames{},
		lendingMap:       map[string]*model.LendingWithNames{},
	}

	var data InitData
	if err := utils.DecodeJson(initData, &data); err != nil {
		return nil, err
	}

	r.AddMembers(data.Members)
	r.AddBooks(data.Books)

	return r, nil
}
