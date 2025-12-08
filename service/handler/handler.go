package handler

import (
	"errors"
	"os"
	"time"

	"library-management/service/models"
	"library-management/service/repository"
)

const finePerDay = 1

func AdminLogin(password string) (string, error) {
	expected := os.Getenv("ADMIN_PASSWORD")
	if expected == "" {
		expected = "123"
	}
	if password != expected {
		return "", errors.New("invalid password")
	}

	token := os.Getenv("ADMIN_TOKEN")
	if token == "" {
		token = "admintoken"
	}
	return token, nil
}

func CreateBook(r *repository.Repo, b *models.Book) (int64, error) {
	if b.Copies <= 0 {
		b.Copies = 1
	}
	if b.Available <= 0 {
		b.Available = b.Copies
	}
	return r.BookRepo.Create(b)
}

func UpdateBook(r *repository.Repo, b *models.Book) error {
	if b.Available > b.Copies {
		b.Available = b.Copies
	}
	return r.BookRepo.Update(b)
}

func DeleteBook(r *repository.Repo, id int64) error {
	return r.BookRepo.Delete(id)
}

func GetBookByID(r *repository.Repo, id int64) (*models.Book, error) {
	return r.BookRepo.GetByID(id)
}

func ListBooks(r *repository.Repo) ([]models.Book, error) {
	return r.BookRepo.GetAll()
}

func SearchBooks(r *repository.Repo, q string) ([]models.Book, error) {
	return r.BookRepo.Search(q)
}

func CreateMember(r *repository.Repo, m *models.Member) (int64, error) {
	return r.MemberRepo.Create(m)
}

func UpdateMember(r *repository.Repo, m *models.Member) error {
	return r.MemberRepo.Update(m)
}

func DeleteMember(r *repository.Repo, id int64) error {
	return r.MemberRepo.Delete(id)
}

func GetMemberByID(r *repository.Repo, id int64) (*models.Member, error) {
	return r.MemberRepo.GetByID(id)
}

func ListMembers(r *repository.Repo) ([]models.Member, error) {
	return r.MemberRepo.GetAll()
}

func IssueBook(r *repository.Repo, bookID, memberID int64, dueDays int) (int64, error) {

	book, err := r.BookRepo.GetByID(bookID)
	if err != nil {
		return 0, err
	}
	if book.Available <= 0 {
		return 0, errors.New("no available copies")
	}

	active, err := r.IssueRepo.GetActiveByBookAndMember(bookID, memberID)
	if err != nil {
		return 0, err
	}
	if active != nil {
		return 0, errors.New("this member already has this book issued")
	}

	var dueDateStr *string
	if dueDays > 0 {
		dt := time.Now().AddDate(0, 0, dueDays).Format("2006-01-02")
		dueDateStr = &dt
	}
	iss := &models.Issue{
		BookID:   bookID,
		MemberID: memberID,
		DueDate:  dueDateStr,
	}
	id, err := r.IssueRepo.Create(iss)
	if err != nil {
		return 0, err
	}
	if err := r.BookRepo.ChangeAvailability(bookID, -1); err != nil {
		return id, err
	}
	return id, nil
}

func ReturnBook(r *repository.Repo, issueID int64) (float64, error) {
	issue, err := r.IssueRepo.GetByID(issueID)
	if err != nil {
		return 0, err
	}
	if issue == nil {
		return 0, errors.New("issue record not found")
	}
	if issue.ReturnedAt != nil {
		return 0, errors.New("already returned")
	}

	var fine float64
	if issue.DueDate != nil && *issue.DueDate != "" {
		due, err := time.Parse("2006-01-02", *issue.DueDate)
		if err == nil && time.Now().After(due) {
			diff := time.Since(due)
			days := int(diff.Hours()/24 + 0.999)
			fine = float64(days) * finePerDay
		}
	}
	now := time.Now()
	if err := r.IssueRepo.Return(issueID, now, fine); err != nil {
		return 0, err
	}
	_ = r.BookRepo.ChangeAvailability(issue.BookID, 1)
	return fine, nil
}

func GetIssuesByMember(r *repository.Repo, memberID int64) ([]models.Issue, error) {
	return r.IssueRepo.GetByMember(memberID)
}
