package services

import "github.com/jiseruk/minesweeper/cmd/minesweeper/models"

type MockBoardDAO struct {
	records []models.Game
}

// Mock Create function that replaces real Board DAO
func (m *MockBoardDAO) Create(board *models.Game) (*models.Game, error) {
	board.ID = len(m.records)
	m.records = append(m.records, *board)

	return board, nil
}
