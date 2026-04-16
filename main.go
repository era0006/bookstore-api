package main

import (
	"fmt"
	"log"

	"github.com/era0006/bookstore-api/database"
	"github.com/era0006/bookstore-api/handlers"
	"github.com/era0006/bookstore-api/models"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	err := database.DB.AutoMigrate(
		&models.Book{},
		&models.Author{},
		&models.Category{},
		&models.User{},
		&models.Favorite{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("✅ Database migrated!")

	r := gin.Default()

	r.GET("/books", handlers.GetBooks)
	r.GET("/books/:id", handlers.GetBookByID)
	r.GET("/authors", handlers.GetAuthors)
	r.GET("/categories", handlers.GetCategories)

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.POST("/books", handlers.CreateBook)
		protected.PUT("/books/:id", handlers.UpdateBook)
		protected.DELETE("/books/:id", handlers.DeleteBook)
		protected.POST("/authors", handlers.CreateAuthor)
		protected.POST("/categories", handlers.CreateCategory)

		protected.GET("/favorites", handlers.GetFavorites)
		protected.PUT("/favorites/books/:bookId", handlers.AddToFavorites)
		protected.DELETE("/favorites/books/:bookId", handlers.RemoveFromFavorites)
	}

	fmt.Println("🚀 Server starting on http://localhost:8080")
	r.Run(":8080")
}
