package libhttp

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	handler "library-management/service/handler"
	"library-management/service/models"
	"library-management/service/repository"
)

func buildRepo(db *sqlx.DB) *repository.Repo {
	return repository.NewRepository(db)
}

func jsonError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func AdminLoginHandler(c *gin.Context) {
	var payload struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonError(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := handler.AdminLogin(payload.Password)
	if err != nil {
		jsonError(c, http.StatusUnauthorized, "invalid password")
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := os.Getenv("ADMIN_TOKEN")
		if w == "" {
			w = "admintoken"
		}

		got := c.GetHeader("X-Admin-Token")
		if got != w {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
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

		id, err := handler.CreateBook(r, &b)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id, "message": "book created"})
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

		b.ID = id
		if err := handler.UpdateBook(r, &b); err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	}
}

func DeleteBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		if err := handler.DeleteBook(r, id); err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}

func GetBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		b, err := handler.GetBookByID(r, id)
		if err != nil {
			jsonError(c, http.StatusNotFound, "not found")
			return
		}
		c.JSON(http.StatusOK, b)
	}
}

func ListBooksHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		books, err := handler.ListBooks(r)
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

		books, err := handler.SearchBooks(r, q)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, books)
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

		id, err := handler.CreateMember(r, &m)
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

		m.ID = id

		if err := handler.UpdateMember(r, &m); err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	}
}

func DeleteMemberHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		if err := handler.DeleteMember(r, id); err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}

func GetMemberHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		m, err := handler.GetMemberByID(r, id)
		if err != nil {
			jsonError(c, http.StatusNotFound, "not found")
			return
		}
		c.JSON(http.StatusOK, m)
	}
}

func ListMembersHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		members, err := handler.ListMembers(r)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, members)
	}
}

func IssueBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		var p struct {
			BookID   int64 `json:"book_id" binding:"required"`
			MemberID int64 `json:"member_id" binding:"required"`
			DueDays  int   `json:"due_days"`
		}

		if err := c.ShouldBindJSON(&p); err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}

		issueID, err := handler.IssueBook(r, p.BookID, p.MemberID, p.DueDays)
		if err != nil {
			jsonError(c, http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{"issue_id": issueID})
	}
}

func ReturnBookHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := buildRepo(db)

		idStr := c.Param("id")
		issueID, _ := strconv.ParseInt(idStr, 10, 64)

		fine, err := handler.ReturnBook(r, issueID)
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

		issues, err := handler.GetIssuesByMember(r, memberID)
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, issues)
	}
}
