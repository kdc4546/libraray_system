package models

import "time"

type Book struct {
	ID        int64     `db:"id" json:"id"`
	Title     string    `db:"title" json:"title" binding:"required"`
	Author    string    `db:"author" json:"author" binding:"required"`
	Copies    int       `db:"copies" json:"copies"`
	Available int       `db:"available" json:"available"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Member struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" binding:"required"`
	Email     string    `db:"email" json:"email"`
	RollNo    string    `db:"roll_no" json:"roll_no"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Issue struct {
	ID         int64      `db:"id" json:"id"`
	BookID     int64      `db:"book_id" json:"book_id"`
	MemberID   int64      `db:"member_id" json:"member_id"`
	IssuedAt   time.Time  `db:"issued_at" json:"issued_at"`
	DueDate    *string    `db:"due_date" json:"due_date"`
	ReturnedAt *time.Time `db:"returned_at" json:"returned_at"`
	FinePaid   float64    `db:"fine_paid" json:"fine_paid"`
}
