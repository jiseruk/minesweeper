package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/apis"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/config"

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
	r := apis.GetRouter()

	config.Config.DB, config.Config.DBErr = gorm.Open("postgres", config.Config.DSN)
	if config.Config.DBErr != nil {
		panic(config.Config.DBErr)
	}

	config.Config.DB.AutoMigrate(&models.Board{}) // This is needed for generation of schema for postgres image.

	defer config.Config.DB.Close()

	fmt.Println(fmt.Sprintf("Successfully connected to :%v", config.Config.DSN))
	r.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
}
