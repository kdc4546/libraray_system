package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"library-management/service/libhttp"
	"library-management/service/repository/db"
)

func main() {
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("failed connect db:", err)
	}
	log.Println("connected to db")
	if err := db.EnsureSchema(database); err != nil {
		log.Fatal("ensure schema:", err)
	}
	r := gin.Default()

	libhttp.RegisterRoutes(r, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server running :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
