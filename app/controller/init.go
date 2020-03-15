package controller

import (
	"fmt"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/service"
	"log"
	"time"
)

func InitEndpoints(
	port int,
	bookService service.IBookService,
) {
	r := gin.Default()

	cacheStore := persistence.NewInMemoryStore(60 * time.Second)

	setCORS(r)
	addEndpoints(r, cacheStore, bookService)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Panicf("failed to start server on %d port. error: %v", port, err)
	}
}

func setCORS(r *gin.Engine) {
	r.Use(cors.Default())
}

func addEndpoints(
	r *gin.Engine,
	cacheStore persistence.CacheStore,
	bookService service.IBookService,
) {
	r.GET("/book/search", cache.CachePage(cacheStore, 5 * time.Minute, searchBookHandlerFn(bookService)))
	r.POST("/books", addBookHandlerFn(bookService))
	r.PATCH("/books/:isbn", updateBookHandlerFn(bookService))
	r.PATCH("/books/:isbn/owned", updateOwnedBookHandlerFn(bookService))
}
