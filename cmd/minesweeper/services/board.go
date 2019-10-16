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
func (s *BoardService) Create(board *models.Board) (*models.Board, error) {
	board.Populate()
	board.Status = models.StatusActive
	return s.dao.Create(board)
}

//Selects a point in the board
func (s *BoardService) SelectPoint(id int, point models.Point) (*models.Board, error) {
	if board, err := s.dao.Get(id); err == nil {
		fmt.Printf("%v", point)
		if point.Mine != nil {
			board.SetMine(point)
		} else {
			board.Select(point)
		}
		if _, err := s.dao.Update(board); err != nil {
			return nil, err
		}
		return board, nil
	} else {
		return nil, err
	}
}

//Set a point as a mine candidate in the board
func (s *BoardService) SelectMine(id int, point models.Point) (*models.Board, error) {
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
func (s *BoardService) Get(id int) (*models.Board, error) {
	if board, err := s.dao.Get(id); err == nil {
		return board, nil
	} else {
		return nil, err
	}
}
