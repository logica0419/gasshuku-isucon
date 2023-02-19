package repository

import "github.com/logica0419/gasshuku-isucon/bench/model"

type MemberRepository interface {
	GetMemberByID(id string) (*model.MemberWithLending, error)
	AddMembers(members []*model.MemberWithLending)
}

var _ MemberRepository = &Repository{}

func (r *Repository) GetMemberByID(id string) (*model.MemberWithLending, error) {
	r.mLock.RLock()
	defer r.mLock.RUnlock()

	v, ok := r.memberMap[id]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func (r *Repository) AddMembers(members []*model.MemberWithLending) {
	r.mLock.Lock()
	defer r.mLock.Unlock()

	r.memberSlice = append(r.memberSlice, members...)
	for _, m := range members {
		r.memberMap[m.ID] = m
	}
}
