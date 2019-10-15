package test_data

import (
	"errors"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
)

type MockBoardDAO struct {
	Records map[int]models.Board
}

func (m *MockBoardDAO) Create(board *models.Board) (*models.Board, error) {
	board.ID = len(m.Records) + 1
	board.Status = models.StatusActive
	m.Records[board.ID] = *board

	return board, nil
}

func (m *MockBoardDAO) Update(board *models.Board) (*models.Board, error) {
	board.ID = len(m.Records) + 1
	board.Status = models.StatusActive
	m.Records[board.ID] = *board

	return board, nil
}

func (m *MockBoardDAO) Get(id int) (*models.Board, error) {
	if board, ok := m.Records[id]; ok {
		return &board, nil
	}
	return nil, errors.New("Not found")
}
