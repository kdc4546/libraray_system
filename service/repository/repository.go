package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"library-management/service/models"
	db "library-management/service/repository/db"
)

type BookRepo interface {
	Create(b *models.Book) (int64, error)
	GetByID(id int64) (*models.Book, error)
	GetAll() ([]models.Book, error)
	Search(qry string) ([]models.Book, error)
	Update(b *models.Book) error
	Delete(id int64) error
	ChangeAvailability(id int64, delta int) (bool, error)
}

type MemberRepo interface {
	Create(m *models.Member) (int64, error)
	GetByID(id int64) (*models.Member, error)
	GetAll() ([]models.Member, error)
	Update(m *models.Member) error
	Delete(id int64) error
}

type IssueRepo interface {
	Create(issue *models.Issue) (int64, error)
	GetActiveByBookAndMember(bookID, memberID int64) (*models.Issue, error)
	GetByMember(memberID int64) ([]models.Issue, error)
	GetByID(id int64) (*models.Issue, error)
	Return(issueID int64, returnedAt time.Time, fine float64) (bool, error)
}

type Repo struct {
	BookRepo   BookRepo
	MemberRepo MemberRepo
	IssueRepo  IssueRepo
}

func NewRepo(dbx *sqlx.DB) *Repo {
	return &Repo{
		BookRepo:   &bookRepository{db: dbx},
		MemberRepo: &memberRepository{db: dbx},
		IssueRepo:  &issueRepository{db: dbx},
	}
}

type bookRepository struct {
	db *sqlx.DB
}

func (r *bookRepository) Create(b *models.Book) (int64, error) {
	res, err := r.db.Exec(db.QCreateBook, b.Title, b.Author, b.Copies, b.Available)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *bookRepository) GetByID(id int64) (*models.Book, error) {
	var b models.Book
	if err := r.db.Get(&b, db.QGetBookByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (r *bookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	if err := r.db.Select(&books, db.QGetAllBooks); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) Search(qry string) ([]models.Book, error) {
	like := "%" + qry + "%"
	var books []models.Book
	if err := r.db.Select(&books, db.QSearchBooks, like, like); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) Update(b *models.Book) error {
	_, err := r.db.Exec(db.QUpdateBook, b.Title, b.Author, b.Copies, b.Available, b.ID)
	return err
}

func (r *bookRepository) Delete(id int64) error {
	_, err := r.db.Exec(db.QDeleteBook, id)
	return err
}

func (r *bookRepository) ChangeAvailability(id int64, delta int) (bool, error) {
	res, err := r.db.Exec(db.QChangeAvailability, delta, id, delta)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

type memberRepository struct {
	db *sqlx.DB
}

func (r *memberRepository) Create(m *models.Member) (int64, error) {
	res, err := r.db.Exec(db.QCreateMember, m.Name, m.Email, m.RollNo)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *memberRepository) GetByID(id int64) (*models.Member, error) {
	var m models.Member
	if err := r.db.Get(&m, db.QGetMemberByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *memberRepository) GetAll() ([]models.Member, error) {
	var members []models.Member
	if err := r.db.Select(&members, db.QGetAllMembers); err != nil {
		return nil, err
	}
	return members, nil
}

func (r *memberRepository) Update(m *models.Member) error {
	_, err := r.db.Exec(db.QUpdateMember, m.Name, m.Email, m.RollNo, m.ID)
	return err
}

func (r *memberRepository) Delete(id int64) error {
	_, err := r.db.Exec(db.QDeleteMember, id)
	return err
}

type issueRepository struct {
	db *sqlx.DB
}

func (r *issueRepository) Create(issue *models.Issue) (int64, error) {
	res, err := r.db.Exec(db.QCreateIssue, issue.BookID, issue.MemberID, issue.DueDate)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *issueRepository) GetActiveByBookAndMember(bookID, memberID int64) (*models.Issue, error) {
	var it models.Issue
	if err := r.db.Get(&it, db.QGetActiveIssueByBookAndMember, bookID, memberID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &it, nil
}

func (r *issueRepository) GetByMember(memberID int64) ([]models.Issue, error) {
	var issues []models.Issue
	if err := r.db.Select(&issues, db.QGetIssuesByMember, memberID); err != nil {
		return nil, err
	}
	return issues, nil
}

func (r *issueRepository) GetByID(id int64) (*models.Issue, error) {
	var it models.Issue
	if err := r.db.Get(&it, db.QGetIssueByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &it, nil
}

func (r *issueRepository) Return(issueID int64, returnedAt time.Time, fine float64) (bool, error) {
	res, err := r.db.Exec(db.QReturnIssue, returnedAt, fine, issueID)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
