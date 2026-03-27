package handlers

import (
	"net/http"

	"github.com/era0006/bookstore-api/models"
	"github.com/gin-gonic/gin"
)

var authors = []models.Author{
	{ID: 1, Name: "Alan A.A. Donovan"},
	{ID: 2, Name: "Robert C. Martin"},
}
var nextAuthorID = 3

func GetAuthors(c *gin.Context) {
	c.JSON(http.StatusOK, authors)
}

func CreateAuthor(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if author.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	author.ID = nextAuthorID
	nextAuthorID++
	authors = append(authors, author)
	c.JSON(http.StatusCreated, author)
}
