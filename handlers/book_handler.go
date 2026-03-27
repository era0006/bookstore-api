package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/era0006/bookstore-api/models"
	"github.com/gin-gonic/gin"
)

var books = []models.Book{
	{ID: 1, Title: "The Go Programming Language", AuthorID: 1, CategoryID: 1, Price: 45.99},
	{ID: 2, Title: "Clean Code", AuthorID: 2, CategoryID: 1, Price: 39.99},
}
var nextBookID = 3

var Authors []models.Author
var Categories []models.Category

func GetBooks(c *gin.Context) {
	filtered := make([]models.Book, len(books))
	copy(filtered, books)

	categoryFilter := c.Query("category")
	if categoryFilter != "" {
		temp := []models.Book{}
		for _, b := range filtered {
			for _, cat := range Categories {
				if cat.ID == b.CategoryID && strings.EqualFold(cat.Name, categoryFilter) {
					temp = append(temp, b)
					break
				}
			}
		}
		filtered = temp
	}

	authorFilter := c.Query("author")
	if authorFilter != "" {
		temp := []models.Book{}
		for _, b := range filtered {
			for _, a := range Authors {
				if a.ID == b.AuthorID && strings.EqualFold(a.Name, authorFilter) {
					temp = append(temp, b)
					break
				}
			}
		}
		filtered = temp
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(filtered) {
		start = len(filtered)
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	c.JSON(http.StatusOK, filtered[start:end])
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if book.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}
	if book.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be positive"})
		return
	}

	book.ID = nextBookID
	nextBookID++
	books = append(books, book)

	c.JSON(http.StatusCreated, book)
}

func GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for _, book := range books {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updated models.Book
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if updated.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}
	if updated.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Price must be positive"})
		return
	}

	for i, book := range books {
		if book.ID == id {
			updated.ID = id
			books[i] = updated
			c.JSON(http.StatusOK, updated)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}
