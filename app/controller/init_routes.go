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

// InitEndpoints initializes gin-gonic router engine and registers http endpoints.
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

// setCORS accepts all requests from any source.
func setCORS(r *gin.Engine) {
	r.Use(cors.Default())
}

// addEntpoints registers http endpoints.
func addEndpoints(
	r *gin.Engine,
	cacheStore persistence.CacheStore,
	bookService service.IBookService,
) {
	r.GET("/book/search", cache.CachePage(cacheStore, 5*time.Minute, SearchBookHandlerFn(bookService)))
	r.POST("/books", AddBookHandlerFn(bookService))
	r.PATCH("/books/:isbn", UpdateBookHandlerFn(bookService))
	r.PATCH("/books/:isbn/owned", UpdateOwnedBookHandlerFn(bookService))
}
