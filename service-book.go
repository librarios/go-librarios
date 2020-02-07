package main

import "github.com/gin-gonic/gin"

type BookPlugin interface {
	FindByISBN(string) (*Book, error)
	FindByTitle(string) ([]*Book, error)
}

func FindByISBN(c *gin.Context) {
	isbn := c.Param("isbn")

	for _, plugin := range pluginManager.GetPluginsByType(PluginTypeBook) {
		book, err := plugin.(BookPlugin).FindByISBN(isbn)
		if err != nil {
			c.Error(err)
			return
		}

		if book == nil {
			c.JSON(200, gin.H{
				"message": "no data",
				"data": nil,
			})
			return
		}

		c.JSON(200, gin.H{
			"data": book,
		})
		return
	}
}

func FindByTitle(c *gin.Context) {
	title := c.Query("title")

	for _, plugin := range pluginManager.GetPluginsByType(PluginTypeBook) {
		book, err := plugin.(BookPlugin).FindByTitle(title)
		if err != nil {
			c.Error(err)
			return
		}

		if book == nil {
			c.JSON(200, gin.H{
				"message": "no data",
				"data": nil,
			})
			return
		}

		c.JSON(200, gin.H{
			"data": book,
		})
		return
	}
}

