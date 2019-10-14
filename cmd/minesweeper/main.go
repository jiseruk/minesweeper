package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/apis"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/jiseruk/minesweeper/cmd/minesweeper/docs"
)

// @title Minesweeper Swagger API
// @version 1.0
// @description Swagger API for Minesweeper API.

// @contact.name Javier Iseruk
// @contact.email javier.iseruk@gmail.com

// @BasePath /api/v1
func main() {
	// load application configurations

	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	// Creates a router without any middleware by default
	r := gin.New()

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.PUT("/boards/:id", apis.SelectPoint)
		v1.POST("/boards/", apis.CreateBoard)
	}

	config.Config.DB, config.Config.DBErr = gorm.Open("postgres", config.Config.DSN)
	if config.Config.DBErr != nil {
		panic(config.Config.DBErr)
	}

	config.Config.DB.AutoMigrate(&models.Board{}) // This is needed for generation of schema for postgres image.

	defer config.Config.DB.Close()

	fmt.Println(fmt.Sprintf("Successfully connected to :%v", config.Config.DSN))
	r.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
}
