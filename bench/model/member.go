package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
	"github.com/logica0419/gasshuku-isucon/bench/utils"
	"github.com/mattn/go-gimei"
)

type Member struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Address     string    `json:"address" db:"address"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Banned      bool      `json:"banned" db:"banned"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type MemberWithLending struct {
	Member
	Lending bool `json:"lending"`
}

func NewMember() *MemberWithLending {
	return &MemberWithLending{
		Member: Member{
			ID:          utils.GenerateID(),
			Name:        NewMemberName(),
			Address:     NewMemberAddress(),
			PhoneNumber: NewMemberPhoneNumber(),
			Banned:      false,
			CreatedAt:   time.Now(),
		},
		Lending: false,
	}
}

func NewMemberName() string {
	if rand.Intn(10) == 0 {
		return fake.FullName()
	}
	return gimei.NewName().Kanji()
}

func NewMemberAddress() string {
	return fmt.Sprintf("%s%d-%d-%d", gimei.NewAddress().Kanji(),
		rand.Intn(10), rand.Intn(100), rand.Intn(1000))
}

func NewMemberPhoneNumber() string {
	if rand.Intn(10) == 0 {
		return fake.Phone()
	}
	return fmt.Sprintf("0%d-%d-%d", rand.Intn(10), rand.Intn(10000), rand.Intn(10000))
}
