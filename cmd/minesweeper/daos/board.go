package daos

import (
	"fmt"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/config"
	//"github.com/MartinHeinz/go-project-blueprint/cmd/blueprint/config"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
)

type BoardDAO interface {
	Create(board *models.Board) (*models.Board, error)
	Update(board *models.Board) (*models.Board, error)
	Get(id int) (*models.Board, error)
}

//BoardDAO Database persistence dao implementation
type BoardDAOImpl struct{}

// Creates a new board
func (dao *BoardDAOImpl) Create(board *models.Board) (*models.Board, error) {
	board.Populate()
	config.Config.DB.Create(&board)
	return board, nil
}

func (dao *BoardDAOImpl) Update(board *models.Board) (*models.Board, error) {
	board.Populate()
	fmt.Printf("ID: %d", board.ID)
	s := config.Config.DB.Update(board)
	fmt.Print(s.Error)
	return board, nil
}

// Get does the actual query to database, if user with specified id is not found error is returned
func (dao *BoardDAOImpl) Get(id int) (*models.Board, error) {
	var board models.Board
	config.Config.DB.Where([]int{id}).First(&board)
	board.Populate()
	return &board, nil
}
