package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/service"
	"net/http"
	"strconv"
)

// AddOwnedBook add owned book
func AddOwnedBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body service.AddOwnedBook
		if err := c.BindJSON(&body); err != nil {
			_ = c.Error(err)
			return
		}

		if book, err := service.BookService.AddOwnedBook(body); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusCreated, gin.H{"data": book})
		}
	}
}

// UpdateOwnedBook update owned book
func UpdateOwnedBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			_ = c.Error(err)
			return
		}
		update := make(gin.H)
		if err := c.BindJSON(&update); err != nil {
			_ = c.Error(err)
			return
		}

		if book, err := service.BookService.UpdateOwnedBook(uint(id), update); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": book})
		}
	}
}

// UpdateBook update book
func UpdateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			_ = c.Error(err)
			return
		}
		update := make(gin.H)
		if err := c.BindJSON(&update); err != nil {
			_ = c.Error(err)
			return
		}

		if book, err := service.BookService.UpdateBook(uint(id), update); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusOK, gin.H{"data": book})
		}
	}
}
