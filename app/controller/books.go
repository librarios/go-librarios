package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/service"
	"net/http"
)

// addBookHandlerFn add book
func addBookHandlerFn(s service.IBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body service.AddBookCommand
		if err := c.BindJSON(&body); err != nil {
			_ = c.Error(err)
			return
		}

		if book, err := s.AddBook(body); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"data": book,
			})
		}
	}
}

// updateBookHandlerFn update book
func updateBookHandlerFn(s service.IBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")
		body := make(gin.H)
		if err := c.BindJSON(&body); err != nil {
			_ = c.Error(err)
			return
		}

		if book, err := s.UpdateBook(isbn, body); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data": book,
			})
		}
	}
}

// updateOwnedBookHandlerFn update book
func updateOwnedBookHandlerFn(s service.IBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")
		body := make(gin.H)
		if err := c.BindJSON(&body); err != nil {
			_ = c.Error(err)
			return
		}

		if book, err := s.UpdateOwnedBook(isbn, body); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data": book,
			})
		}
	}
}
