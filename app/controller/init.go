package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/service"
	"log"
)

func InitEndpoints(
	port int,
	bookService service.IBookService,
) {
	r := gin.Default()

	r.GET("/book/search", searchBookHandlerFn(bookService))
	r.POST("/book", addBookHandlerFn(bookService))

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Panicf("failed to start server on %d port. error: %v", port, err)
	}
}
