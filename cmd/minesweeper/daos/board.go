package daos

import (
	"fmt"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/config"
	//"github.com/MartinHeinz/go-project-blueprint/cmd/blueprint/config"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
)

type BoardDAO interface {
	Create(board *models.Game) (*models.Game, error)
	Update(board *models.Game) (*models.Game, error)
	Get(id int) (*models.Game, error)
}

//BoardDAO Database persistence dao implementation
type BoardDAOImpl struct{}

// Creates a new board
func (dao *BoardDAOImpl) Create(board *models.Game) (*models.Game, error) {
	config.Config.DB.Create(&board)
	return board, nil
}

func (dao *BoardDAOImpl) Update(board *models.Game) (*models.Game, error) {
	if err := config.Config.DB.Save(&board).Error; err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return board, nil
}

// Get does the actual query to database, if user with specified id is not found error is returned
func (dao *BoardDAOImpl) Get(id int) (*models.Game, error) {
	var board models.Game
	config.Config.DB.Where([]int{id}).First(&board)
	board.GetBoardFromDB()
	return &board, nil
}
