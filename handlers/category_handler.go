package handlers

import (
	"net/http"

	"github.com/era0006/bookstore-api/models"
	"github.com/gin-gonic/gin"
)

var categories = []models.Category{
	{ID: 1, Name: "Programming"},
	{ID: 2, Name: "Software Design"},
}
var nextCategoryID = 3

func GetCategories(c *gin.Context) {
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	category.ID = nextCategoryID
	nextCategoryID++
	categories = append(categories, category)
	c.JSON(http.StatusCreated, category)
}
