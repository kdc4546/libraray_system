package db

const (
	QCreateBook         = `INSERT INTO books (title, author, copies, available) VALUES (?, ?, ?, ?)`
	QGetBookByID        = `SELECT id, title, author, copies, available, created_at, updated_at FROM books WHERE id = ? LIMIT 1`
	QGetAllBooks        = `SELECT id, title, author, copies, available, created_at, updated_at FROM books ORDER BY id DESC`
	QSearchBooks        = `SELECT id, title, author, copies, available, created_at, updated_at FROM books WHERE title LIKE ? OR author LIKE ? ORDER BY id DESC`
	QUpdateBook         = `UPDATE books SET title = ?, author = ?, copies = ?, available = ? WHERE id = ?`
	QDeleteBook         = `DELETE FROM books WHERE id = ?`
	QChangeAvailability = `UPDATE books SET available = available + ? WHERE id = ?`

	QCreateMember  = `INSERT INTO members (name, email, roll_no) VALUES (?, ?, ?)`
	QGetMemberByID = `SELECT id, name, email, roll_no, created_at, updated_at FROM members WHERE id = ? LIMIT 1`
	QGetAllMembers = `SELECT id, name, email, roll_no, created_at, updated_at FROM members ORDER BY id DESC`
	QUpdateMember  = `UPDATE members SET name = ?, email = ?, roll_no = ? WHERE id = ?`
	QDeleteMember  = `DELETE FROM members WHERE id = ?`

	QCreateIssue                   = `INSERT INTO issues (book_id, member_id, due_date) VALUES (?, ?, ?)`
	QGetActiveIssueByBookAndMember = `SELECT id, book_id, member_id, issued_at, due_date, returned_at, fine_paid FROM issues WHERE book_id = ? AND member_id = ? AND returned_at IS NULL LIMIT 1`
	QGetIssuesByMember             = `SELECT id, book_id, member_id, issued_at, due_date, returned_at, fine_paid FROM issues WHERE member_id = ? ORDER BY issued_at DESC`
	QGetIssueByID                  = `SELECT id, book_id, member_id, issued_at, due_date, returned_at, fine_paid FROM issues WHERE id = ? LIMIT 1`
	QReturnIssue                   = `UPDATE issues SET returned_at = ?, fine_paid = ? WHERE id = ?`
)
