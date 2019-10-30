package test_data

import (
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
	//"fmt"
	//"github.com/jiseruk/minesweeper/cmd/minesweeper/config"
	//"github.com/MartinHeinz/go-project-blueprint/cmd/blueprint/models"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	//"io/ioutil"
	//"strings"
)

func GetTestCaseFolder() string {
	return "/test_data"
}

func GetBoard(width int, height int, mines []models.Point) models.Game {
	board := models.Game{Width: width, Height: height, Mines: len(mines), Status: models.StatusActive, UnknownMines: len(mines),
		UnselectedPoints: width*height - len(mines), Model: models.Model{ID: 1}}
	board.Populate(mines...)
	return board
}
