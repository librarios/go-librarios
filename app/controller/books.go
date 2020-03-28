package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/service"
	"log"
	"net/http"
)

// AddOwnedBook add owned book
func AddOwnedBook(s service.IBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body service.AddOwnedBook
		if err := c.BindJSON(&body); err != nil {
			_ = c.Error(err)
			return
		}

		if book, err := s.AddOwnedBook(body); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"data": book,
			})
		}
	}
}

// UpdateOwnedBook update owned book
func UpdateOwnedBook(s service.IBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")
		body := service.UpdateOwnedBook{}
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

// UpdateBook update book
func UpdateBook(s service.IBookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		isbn := c.Param("isbn")
		body := service.UpdateBook{}
		if err := c.BindJSON(&body); err != nil {
			_ = c.Error(err)
			return
		}
		log.Printf("body: %#v", body)

		if book, err := s.UpdateBook(isbn, body); err != nil {
			_ = c.Error(err)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"data": book,
			})
		}
	}
}
