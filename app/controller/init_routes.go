package controller

import (
	"fmt"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// InitEndpoints initializes gin-gonic router engine and registers http endpoints.
func InitEndpoints(
	port int,
) {
	r := gin.Default()

	// middleware
	addMiddlewares(r)

	cacheStore := persistence.NewInMemoryStore(60 * time.Second)

	setCORS(r)
	addEndpoints(r, cacheStore)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Panicf("failed to start server on %d port. error: %v", port, err)
	}
}

// setCORS accepts all requests from any source.
func setCORS(r *gin.Engine) {
	r.Use(cors.Default())
}

// addMiddlewares registers middlewares
func addMiddlewares(r *gin.Engine) {
}

// addEntpoints registers http endpoints.
func addEndpoints(
	r *gin.Engine,
	cacheStore persistence.CacheStore,
) {
	r.GET("/book/search", cache.CachePage(cacheStore, 5*time.Minute, SearchBook()))
	r.POST("/books/own", AddOwnedBook())
	r.PATCH("/books/own/:id", UpdateOwnedBook())
	r.PATCH("/books/book/:id", UpdateBook())

	r.Group("/graphql").POST("", GraphqlHandler())
	r.Group("/graphiql").GET("", GraphiqlHandler("/graphql"))
}
