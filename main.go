package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	filename := "config/librarios.yaml"
	config, err := LoadConfigFile(filename)
	if err != nil {
		log.Panicf("failed to load config file: %s. error: %v", filename, err)
	}
	log.Printf("Loaded: %s\n", filename)

	// register plugins
	pluginManager.Register(kakaoDef)

	// init plugins
	for name, props := range config.Plugins {
		if plugin, ok := pluginManager.GetPluginByName(name); ok {
			plugin.SetProperties(props)
		}
	}

	// set gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// router
	r := gin.Default()

	r.GET("/book/search", SearchBook)
	if err := r.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.Panicf("failed to start server on %d port. error: %v", config.Port, err)
	}
}
