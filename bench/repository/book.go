package repository

import (
	"math/rand"

	"github.com/logica0419/gasshuku-isucon/bench/model"
)

type BookRepository interface {
	GetBookByID(id string) (*model.BookWithLending, error)
	GetRandomBook() *model.BookWithLending
	AddBooks(books []*model.BookWithLending)
}

var _ BookRepository = &Repository{}

func (r *Repository) GetBookByID(id string) (*model.BookWithLending, error) {
	r.bLock.RLock()
	defer r.bLock.RUnlock()

	v, ok := r.bookMap[id]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func (r *Repository) GetRandomBook() *model.BookWithLending {
	r.bLock.RLock()
	defer r.bLock.RUnlock()

	return r.bookSlice[rand.Intn(len(r.bookSlice))]
}

func (r *Repository) AddBooks(books []*model.BookWithLending) {
	r.bLock.Lock()
	defer r.bLock.Unlock()

	r.bookSlice = append(r.bookSlice, books...)
	for _, book := range books {
		r.bookMap[book.ID] = book
	}
}
