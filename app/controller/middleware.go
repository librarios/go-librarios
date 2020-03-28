package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
)

// RequestLoggerMiddleware returns middleware that logs request body.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Printf("body: %s", string(body))
		log.Printf("header: %s", c.Request.Header)
		c.Next()
	}
}
