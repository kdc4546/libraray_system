package libhttp

import (
	"net/http"
	"os"
	"strconv"

	svc "library-management/service/handler"
	"library-management/service/models"
	"library-management/service/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func buildRepo(db *sqlx.DB) *repository.Repo {
	return repository.NewRepo(db)
}

func jsonError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

type adminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AdminLoginHandler(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	adminUser := getenvOrDefault("ADMIN_USER", "admin")
	adminPass := getenvOrDefault("ADMIN_PASS", "123")

	if req.Username == adminUser && req.Password == adminPass {
		c.JSON(http.StatusOK, gin.H{"message": "login successful"})
		return
	}

	jsonError(c, http.StatusUnauthorized, "invalid credentials")
}

func getenvOrDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func ListBooksHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)
		books, err := svc.ListBooks(r)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func SearchBooksHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)
		q := c.Query("q")
		books, err := svc.SearchBooks(r, q)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func GetBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)
		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		book, err := svc.GetBook(r, id)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if book == nil {
			jsonError(c, http.StatusNotFound, "book not found")
			return
		}
		c.JSON(http.StatusOK, book)
	}
}
func CreateBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)
		var b models.Book
		if err := c.ShouldBindJSON(&b); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}

		id, err := svc.CreateBook(r, &b)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func UpdateBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		var b models.Book
		if err := c.ShouldBindJSON(&b); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}

		if err := svc.UpdateBook(r, id, &b); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func DeleteBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		if err := svc.DeleteBook(r, id); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func ListMembersHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)
		members, err := svc.ListMembers(r)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, members)
	}
}

func GetMemberHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		member, err := svc.GetMember(r, id)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if member == nil {
			jsonError(c, http.StatusNotFound, "member not found")
			return
		}
		c.JSON(http.StatusOK, member)
	}
}

func CreateMemberHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		var m models.Member
		if err := c.ShouldBindJSON(&m); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}

		id, err := svc.CreateMember(r, &m)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func UpdateMemberHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		var m models.Member
		if err := c.ShouldBindJSON(&m); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}

		if err := svc.UpdateMember(r, id, &m); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func DeleteMemberHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		if err := svc.DeleteMember(r, id); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.Status(http.StatusNoContent)
	}
}

type issueRequest struct {
	BookID   int64 `json:"book_id" binding:"required"`
	MemberID int64 `json:"member_id" binding:"required"`
	DueDays  int   `json:"due_days"`
}

func IssueBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		var req issueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}

		id, err := svc.IssueBook(r, req.BookID, req.MemberID, req.DueDays)
		if err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func ReturnBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		issueID, _ := strconv.ParseInt(idStr, 10, 64)

		fine, err := svc.ReturnBook(r, issueID)
		if err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"fine": fine})
	}
}

func IssuesByMemberHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("member_id")
		memberID, _ := strconv.ParseInt(idStr, 10, 64)

		issues, err := svc.GetIssuesByMember(r, memberID)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, issues)
	}
}
