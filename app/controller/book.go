package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/service"
)

// SearchBook returns HandlerFunc that searches book
func SearchBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Query("isbn")
		publisher := c.Query("publisher")
		person := c.Query("person")
		title := c.Query("title")

		books, err := service.BookService.Search(isbn, publisher, person, title)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(200, gin.H{
			"data": books,
		})
	}
}
