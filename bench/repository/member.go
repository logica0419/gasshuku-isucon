package repository

import (
	"math/rand"

	"github.com/logica0419/gasshuku-isucon/bench/model"
)

type MemberRepository interface {
	GetInactiveMemberID(num int) ([]string, error)

	GetMemberTotal() int
	GetMemberByID(id string) (*model.MemberWithLending, error)
	GetRandomMember() *model.MemberWithLending
	AddMembers(members []*model.MemberWithLending)
	DeleteMember(id string)
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

func (r *Repository) DeleteMember(id string) {
	r.mLock.Lock()
	defer r.mLock.Unlock()

	delete(r.memberMap, id)
	for i, m := range r.memberSlice {
		if m.ID == id {
			r.memberSlice = append(r.memberSlice[:i], r.memberSlice[i+1:]...)
			return
		}
	}
}
