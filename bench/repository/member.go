package repository

import (
	"math/rand"

	"github.com/logica0419/gasshuku-isucon/bench/model"
)

type MemberRepository interface {
	GetInactiveMemberID(num int) ([]string, error)

	GetNotLendingMember() (*model.MemberWithLending, error)

	GetMemberTotal() int
	GetMemberByID(id string) (*model.MemberWithLending, error)
	GetRandomMember() *model.MemberWithLending
	AddMembers(members []*model.MemberWithLending)
}

var _ MemberRepository = &Repository{}

func (r *Repository) GetInactiveMemberID(num int) ([]string, error) {
	r.mLock.Lock()
	defer r.mLock.Unlock()

	if len(r.inactiveMemberID) < num {
		return nil, ErrNotEnoughRecords
	}

	mem := r.inactiveMemberID[:num]
	r.inactiveMemberID = r.inactiveMemberID[num:]
	return mem, nil
}

func (r *Repository) GetNotLendingMember() (*model.MemberWithLending, error) {
	r.mLock.RLock()
	defer r.mLock.RUnlock()

	startIndex := rand.Intn(len(r.memberSlice))
	slice := append(r.memberSlice[startIndex:], r.memberSlice[:startIndex]...)

	for _, m := range slice {
		if !m.Lending {
			return m, nil
		}
	}
	return nil, ErrNotFound
}

func (r *Repository) GetMemberTotal() int {
	r.mLock.RLock()
	defer r.mLock.RUnlock()

	return len(r.memberSlice)
}

func (r *Repository) GetMemberByID(id string) (*model.MemberWithLending, error) {
	r.mLock.RLock()
	defer r.mLock.RUnlock()

	v, ok := r.memberMap[id]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func (r *Repository) GetRandomMember() *model.MemberWithLending {
	r.mLock.RLock()
	defer r.mLock.RUnlock()

	return r.memberSlice[rand.Intn(len(r.memberSlice))]
}

func (r *Repository) AddMembers(members []*model.MemberWithLending) {
	r.mLock.Lock()
	defer r.mLock.Unlock()

	r.memberSlice = append(r.memberSlice, members...)
	for _, m := range members {
		r.memberMap[m.ID] = m
		r.inactiveMemberID = append(r.inactiveMemberID, m.ID)
	}
}
