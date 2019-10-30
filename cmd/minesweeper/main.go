package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/apis"

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
	// Creates a router without any middleware by default
	r := apis.GetRouter()
	/*
		// load application configurations
		if err := config.LoadConfig("./config"); err != nil {
			panic(fmt.Errorf("invalid application configuration: %s", err))
		}
		config.Config.DB, config.Config.DBErr = gorm.Open(config.Config.Dialect, config.Config.DSN)
		if config.Config.DBErr != nil {
			panic(config.Config.DBErr)
		}

		config.Config.DB.AutoMigrate(&models.Board{}) // This is needed for generation of schema for postgres image.

		defer config.Config.DB.Close()

		fmt.Println(fmt.Sprintf("Successfully connected to :%v", config.Config.DSN))*/
	defer config.Config.DB.Close()

	r.Run(fmt.Sprintf(":%v", config.Config.ServerPort))
	//http.HandleFunc("/", handler.Handler)
	//http.ListenAndServe(":8080", nil)
}
