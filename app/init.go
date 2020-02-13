package app

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/controller"
	"github.com/librarios/go-librarios/app/service"
	"log"
	"os"
)

func StartServer(configFilename string) {
	config, err := LoadConfigFile(configFilename)
	if err != nil {
		log.Panicf("failed to load config file: %s. error: %v", configFilename, err)
	}
	log.Printf("Loaded: %s\n", configFilename)

	// init plugins
	service.InitPlugins(config.Plugins)

	// connect to DB
	dbConn, err = ConnectDB(config.DB)
	if err != nil {
		log.Panicf("failed to connect DB. err: %v", err)
	}

	// set gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// services
	bookService := service.NewBookService()

	// router
	controller.InitEndpoints(config.Port, bookService)
}