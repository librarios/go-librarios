package service

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/model"
)

type BookPlugin interface {
	FindByISBN(string) ([]*model.Book, error)
	FindByPerson(string) ([]*model.Book, error)
	FindByPublisher(string) ([]*model.Book, error)
	FindByTitle(string) ([]*model.Book, error)
}

func SearchBook(c *gin.Context) {
	isbn := c.Query("isbn")
	publisher := c.Query("publisher")
	person := c.Query("person")
	title := c.Query("title")

	var fn func(string) ([]*model.Book, error) = nil

	for _, plugin := range pluginManager.GetPluginsByType(PluginTypeBook) {
		bookPlugin := plugin.(BookPlugin)
		if isbn != "" {
			fn = bookPlugin.FindByISBN
		} else if publisher != "" {
			fn = bookPlugin.FindByPublisher
		} else if person != "" {
			fn = bookPlugin.FindByPerson
		} else if title != "" {
			fn = bookPlugin.FindByTitle
		} else {
			continue
		}

		books, err := fn(title)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(200, gin.H{
			"data": books,
		})
		return
	}
}
