package repository

import "github.com/logica0419/gasshuku-isucon/bench/model"

type LendingRepository interface {
	GetLendingByID(id string) (*model.LendingWithNames, error)
	GetLendingsByMemberID(id string) ([]*model.LendingWithNames, error)
	AddLendings(lendings []*model.LendingWithNames) error
}

var _ LendingRepository = &Repository{}

func (r *Repository) GetLendingByID(id string) (*model.LendingWithNames, error) {
	r.lLock.RLock()
	defer r.lLock.RUnlock()

	v, ok := r.lendingMap[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (r *Repository) GetLendingsByMemberID(id string) ([]*model.LendingWithNames, error) {
	r.lLock.RLock()
	defer r.lLock.RUnlock()

	v, ok := r.lendingMemberMap[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (r *Repository) AddLendings(lendings []*model.LendingWithNames) error {
	r.lLock.Lock()
	defer r.lLock.Unlock()

	for _, l := range lendings {
		if _, ok := r.lendingMap[l.MemberID]; !ok {
			r.lendingMemberMap[l.MemberID] = make([]*model.LendingWithNames, 0)
		}
		r.lendingMemberMap[l.MemberID] = append(r.lendingMemberMap[l.MemberID], l)
		r.lendingMap[l.ID] = l
	}
	return nil
}
