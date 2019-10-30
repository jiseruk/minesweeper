package services

import (
	"fmt"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/daos"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
)

type BoardService struct {
	dao daos.BoardDAO
}

// NewBoardService creates a new UserService with the given user DAO.
func NewBoardService(dao daos.BoardDAO) *BoardService {
	return &BoardService{dao}
}

// Creates Board DAO, here can be additional logic for processing data retrieved by DAOs
func (s *BoardService) Create(board *models.Game) (*models.Game, error) {
	board.Populate()
	board.Status = models.StatusActive
	board.UnknownMines = board.Mines
	board.UnselectedPoints = board.Width*board.Height - board.Mines
	return s.dao.Create(board)
}

//Selects a point in the board
func (s *BoardService) SelectPoint(id int, point models.Point) (*models.Game, error) {
	point.X--
	point.Y--
	var err error
	if board, dberr := s.dao.Get(id); dberr == nil {
		fmt.Printf("%v", point)
		if point.Mine != nil {
			err = board.SetMine(point)
		} else {
			err = board.Select(point)
		}
		if _, dberr := s.dao.Update(board); dberr != nil {
			board.GetBoardFromDB()
			return nil, dberr
		}
		return board, err
	} else {
		return nil, dberr
	}
}

//Set a point as a mine candidate in the board
func (s *BoardService) SelectMine(id int, point models.Point) (*models.Game, error) {
	if board, err := s.dao.Get(id); err == nil {
		board.SetMine(point)
		if _, err := s.dao.Update(board); err != nil {
			return nil, err
		}
		return board, nil
	} else {
		return nil, err
	}
}

//Get the board
func (s *BoardService) Get(id int) (*models.Game, error) {
	if board, err := s.dao.Get(id); err == nil {
		return board, nil
	} else {
		return nil, err
	}
}
