package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/era0006/bookstore-api/database"
	"github.com/era0006/bookstore-api/models"
	"github.com/gin-gonic/gin"
)

func AddToFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if result := database.DB.First(&book, bookID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var existing models.Favorite
	result := database.DB.Where("user_id = ? AND book_id = ?", userID, bookID).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Book already in favorites"})
		return
	}

	favorite := models.Favorite{
		UserID:    userID.(uint),
		BookID:    uint(bookID),
		CreatedAt: time.Now(),
	}

	database.DB.Create(&favorite)
	c.JSON(http.StatusCreated, gin.H{"message": "Book added to favorites"})
}

func RemoveFromFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	bookID, err := strconv.Atoi(c.Param("bookId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	result := database.DB.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not in favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book removed from favorites"})
}

func GetFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var favorites []models.Favorite
	database.DB.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&favorites)

	var books []models.Book
	for _, fav := range favorites {
		var book models.Book
		if result := database.DB.First(&book, fav.BookID); result.Error == nil {
			books = append(books, book)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"books":     books,
		"page":      page,
		"page_size": pageSize,
		"total":     len(books),
	})
}
