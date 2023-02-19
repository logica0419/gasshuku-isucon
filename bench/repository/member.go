package repository

import "github.com/logica0419/gasshuku-isucon/bench/model"

type MemberRepository interface {
	GetMemberByID(id string) *model.MemberWithLending
	AddMembers(members []*model.MemberWithLending)
}

var _ MemberRepository = &Repository{}

func (r *Repository) GetMemberByID(id string) *model.MemberWithLending {
	r.mLock.RLock()
	defer r.mLock.RUnlock()

	return r.memberMap[id]
}

func (r *Repository) AddMembers(members []*model.MemberWithLending) {
	r.mLock.Lock()
	defer r.mLock.Unlock()

	r.memberSlice = append(r.memberSlice, members...)
	for _, m := range members {
		r.memberMap[m.ID] = m
	}
}
