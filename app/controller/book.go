package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/service"
	"net/http"
)

// searchBookHandlerFn returns HandlerFunc that searches book
func searchBookHandlerFn(s service.IBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Query("isbn")
		publisher := c.Query("publisher")
		person := c.Query("person")
		title := c.Query("title")

		books, err := s.Search(isbn, publisher, person, title)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(200, gin.H{
			"data": books,
		})
	}
}

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
