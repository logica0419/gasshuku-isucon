package repository

import "github.com/logica0419/gasshuku-isucon/bench/model"

type LendingRepository interface {
	GetLendingByID(id string) (*model.LendingWithNames, error)
	GetLendingsByMemberID(id string) ([]*model.LendingWithNames, error)
	AddLendings(lendings []*model.LendingWithNames)
	DeleteLendings(memberID string)
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

func (r *Repository) AddLendings(lendings []*model.LendingWithNames) {
	r.lLock.RLock()
	if _, ok := r.memberMap[lendings[0].MemberID]; !ok {
		r.lLock.RUnlock()
		return
	}
	r.lLock.RUnlock()

	r.lLock.Lock()
	r.mLock.Lock()
	r.bLock.Lock()
	defer r.lLock.Unlock()
	defer r.mLock.Unlock()
	defer r.bLock.Unlock()

	for _, l := range lendings {
		if _, ok := r.lendingMap[l.MemberID]; !ok {
			r.lendingMemberMap[l.MemberID] = make([]*model.LendingWithNames, 0)
		}
		r.lendingMemberMap[l.MemberID] = append(r.lendingMemberMap[l.MemberID], l)
		r.lendingMap[l.ID] = l

		r.memberMap[l.MemberID].Lending = true
		r.bookMap[l.BookID].Lending = true
		for _, m := range r.memberSlice {
			if m.ID == l.MemberID {
				m.Lending = true
				break
			}
		}
		for _, b := range r.bookSlice {
			if b.ID == l.BookID {
				b.Lending = true
				break
			}
		}
	}
}

func (r *Repository) DeleteLendings(memberID string) {
	r.lLock.RLock()
	if _, ok := r.memberMap[memberID]; !ok {
		r.lLock.RUnlock()
		return
	}
	r.lLock.RUnlock()

	r.lLock.Lock()
	r.mLock.Lock()
	r.bLock.Lock()
	defer r.lLock.Unlock()
	defer r.mLock.Unlock()
	defer r.bLock.Unlock()

	slice, ok := r.lendingMemberMap[memberID]
	if !ok {
		return
	}

	r.memberMap[memberID].Lending = false
	for _, m := range r.memberSlice {
		if m.ID == memberID {
			m.Lending = false
			break
		}
	}

	for _, l := range slice {
		delete(r.lendingMap, l.ID)
		r.bookMap[l.BookID].Lending = false
		for _, b := range r.bookSlice {
			if b.ID == l.BookID {
				b.Lending = false
				break
			}
		}
	}

	delete(r.lendingMemberMap, memberID)
}
