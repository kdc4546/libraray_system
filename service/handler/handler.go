package handler

import (
	"errors"
	"time"

	"library-management/service/models"
	"library-management/service/repository"
)

const finePerDay = 10.0

func ListBooks(r *repository.Repo) ([]models.Book, error) {
	return r.BookRepo.GetAll()
}

func SearchBooks(r *repository.Repo, query string) ([]models.Book, error) {
	return r.BookRepo.Search(query)
}

func GetBook(r *repository.Repo, id int64) (*models.Book, error) {
	return r.BookRepo.GetByID(id)
}

func CreateBook(r *repository.Repo, b *models.Book) (int64, error) {
	if b.Available == 0 && b.Copies > 0 {
		b.Available = b.Copies
	}
	return r.BookRepo.Create(b)
}

func UpdateBook(r *repository.Repo, id int64, input *models.Book) error {
	existing, err := r.BookRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("book not found")
	}

	existing.Title = input.Title
	existing.Author = input.Author
	existing.Copies = input.Copies
	existing.Available = input.Available

	return r.BookRepo.Update(existing)
}

func DeleteBook(r *repository.Repo, id int64) error {
	return r.BookRepo.Delete(id)
}

func ListMembers(r *repository.Repo) ([]models.Member, error) {
	return r.MemberRepo.GetAll()
}

func GetMember(r *repository.Repo, id int64) (*models.Member, error) {
	return r.MemberRepo.GetByID(id)
}

func CreateMember(r *repository.Repo, m *models.Member) (int64, error) {
	return r.MemberRepo.Create(m)
}

func UpdateMember(r *repository.Repo, id int64, input *models.Member) error {
	existing, err := r.MemberRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("member not found")
	}

	existing.Name = input.Name
	existing.Email = input.Email
	existing.RollNo = input.RollNo

	return r.MemberRepo.Update(existing)
}

func DeleteMember(r *repository.Repo, id int64) error {
	return r.MemberRepo.Delete(id)
}

func IssueBook(r *repository.Repo, bookID, memberID int64, dueDays int) (int64, error) {

	book, err := r.BookRepo.GetByID(bookID)
	if err != nil {
		return 0, err
	}
	if book == nil {
		return 0, errors.New("book not found")
	}

	member, err := r.MemberRepo.GetByID(memberID)
	if err != nil {
		return 0, err
	}
	if member == nil {
		return 0, errors.New("member not found")
	}

	active, err := r.IssueRepo.GetActiveByBookAndMember(bookID, memberID)
	if err != nil {
		return 0, err
	}
	if active != nil {
		return 0, errors.New("this member already has this book issued")
	}

	ok, err := r.BookRepo.ChangeAvailability(bookID, -1)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("no available copies")
	}

	var dueDateStr *string
	if dueDays > 0 {
		d := time.Now().AddDate(0, 0, dueDays).Format("2006-01-02")
		dueDateStr = &d
	}

	issue := &models.Issue{
		BookID:   bookID,
		MemberID: memberID,
		DueDate:  dueDateStr,
	}

	issueID, err := r.IssueRepo.Create(issue)
	if err != nil {

		_, _ = r.BookRepo.ChangeAvailability(bookID, 1)
		return 0, err
	}

	return issueID, nil
}

func ReturnBook(r *repository.Repo, issueID int64) (float64, error) {
	issue, err := r.IssueRepo.GetByID(issueID)
	if err != nil {
		return 0, err
	}
	if issue == nil {
		return 0, errors.New("issue record not found")
	}

	var fine float64
	if issue.DueDate != nil && *issue.DueDate != "" {
		if due, err := time.Parse("2006-01-02", *issue.DueDate); err == nil && time.Now().After(due) {
			diff := time.Since(due)
			days := int(diff.Hours() / 24)
			if days < 1 {
				days = 1
			}
			fine = float64(days) * finePerDay
		}
	}

	now := time.Now()

	updated, err := r.IssueRepo.Return(issueID, now, fine)
	if err != nil {
		return 0, err
	}
	if !updated {
		return 0, errors.New("already returned")
	}

	_, _ = r.BookRepo.ChangeAvailability(issue.BookID, 1)

	return fine, nil
}

func GetIssuesByMember(r *repository.Repo, memberID int64) ([]models.Issue, error) {
	return r.IssueRepo.GetByMember(memberID)
}
