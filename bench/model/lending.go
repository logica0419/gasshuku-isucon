package model

import "time"

type Lending struct {
	ID        string    `json:"id" db:"id"`
	MemberID  string    `json:"member_id" db:"member_id"`
	BookID    string    `json:"book_id" db:"book_id"`
	Due       time.Time `json:"due" db:"due"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LendingWithNames struct {
	Lending
	MemberName string `json:"member_name"`
	BookTitle  string `json:"book_title"`
}
