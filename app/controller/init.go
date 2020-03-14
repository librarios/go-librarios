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
	addEndpoints(r, bookService)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Panicf("failed to start server on %d port. error: %v", port, err)
	}
}

func addEndpoints(r *gin.Engine, bookService service.IBookService) {
	r.GET("/book/search", searchBookHandlerFn(bookService))
	r.POST("/books", addBookHandlerFn(bookService))
	r.PATCH("/books/:isbn", updateBookHandlerFn(bookService))
	r.PATCH("/books/:isbn/owned", updateOwnedBookHandlerFn(bookService))
}
