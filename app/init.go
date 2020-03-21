package app

import (
	"github.com/gin-gonic/gin"
	"github.com/librarios/go-librarios/app/config"
	"github.com/librarios/go-librarios/app/controller"
	"github.com/librarios/go-librarios/app/model"
	"github.com/librarios/go-librarios/app/plugin"
	"github.com/librarios/go-librarios/app/service"
	"log"
	"os"
)

func StartServer(configFilename string) {
	c, err := config.LoadConfigFile(configFilename)
	if err != nil {
		log.Panicf("failed to load config file: %s. error: %v", configFilename, err)
	}
	log.Printf("Loaded: %s\n", configFilename)

	// init plugins
	plugin.InitPlugins(c.Plugins)

	// connect to DB
	if err = config.InitDB(c.DB); err != nil {
		log.Panicf("failed to connect DB. err: %v", err)
	}
	// auto-migrate
	if c.DB["autoMigrate"] == true {
		model.AutoMigrate()
	}


	defer config.CloseDB()

	// set gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// services
	bookService := service.NewBookService()

	// router
	controller.InitEndpoints(c.Port, bookService)
}
