package libhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *gin.Engine, db *sqlx.DB) {

	r.POST("/admin/login", AdminLoginHandler)
	r.GET("/books", ListBooksHandler(db))
	r.GET("/books/search", SearchBooksHandler(db))
	r.GET("/books/:id", GetBookHandler(db))
	r.GET("/members/:id", GetMemberHandler(db))

	admin := r.Group("/admin")
	admin.Use(AdminAuthMiddleware())

	admin.POST("/books", CreateBookHandler(db))
	admin.PUT("/books/:id", UpdateBookHandler(db))
	admin.DELETE("/books/:id", DeleteBookHandler(db))
	admin.GET("/books", ListBooksHandler(db))

	admin.POST("/members", CreateMemberHandler(db))
	admin.PUT("/members/:id", UpdateMemberHandler(db))
	admin.DELETE("/members/:id", DeleteMemberHandler(db))
	admin.GET("/members", ListMembersHandler(db))

	admin.POST("/issues", IssueBookHandler(db))
	admin.POST("/issues/:id/return", ReturnBookHandler(db))
	admin.GET("/issues/member/:member_id", IssuesByMemberHandler(db))
}
