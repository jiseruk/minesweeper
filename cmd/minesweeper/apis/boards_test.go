package apis

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/services"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/test_data"
)

func TestBoard(t *testing.T) {
	path := test_data.GetTestCaseFolder()
	runTestsSuite(t, []apiTestCase{
		{"t1 - create a board", "POST", "/api/v1/boards/", "/api/v1/boards/", `{"width":7, "height":7, "mines":5}`, CreateBoard, http.StatusCreated, path + "/board_tc1.json"},
		{"t2 - select point", "PUT", "/api/v1/boards/:id", "/api/v1/boards/2", `{"x":1, "y":2}`, SelectPoint, http.StatusOK, path + "/board_tc1.json"},
		{"t3 - Get a board", "GET", "/api/v1/boards/:id", "/api/v1/boards/2", ``, Get, http.StatusOK, path + "/board_tc1.json"},
	})
}

func TestSelectMineNeighboorPoint(t *testing.T) {
	BoardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: map[int]models.Game{1: test_data.GetBoard(3, 3, []models.Point{models.Point{X: 2, Y: 2}})},
	})
	res := testAPICall(http.MethodPut, `/api/v1/boards/1`, `{"x":1, "y":1}`)
	assert.JSONEq(t, `{"id":1, "width":3, "height":3, "mines":1, "matrix":[["1","-","-"],["-","-","-"],["-","-","-"]], "status":"ACTIVE"}`, res.Body.String())
}

func TestSelectPointWithMine(t *testing.T) {
	BoardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: map[int]models.Game{1: test_data.GetBoard(3, 3, []models.Point{models.Point{X: 2, Y: 2}})},
	})
	res := testAPICall(http.MethodPut, `/api/v1/boards/1`, `{"x":2, "y":2}`)
	assert.JSONEq(t, `{"id":1, "width":3, "height":3, "mines":1, "matrix":[["1","1","1"],["1","X","1"],["1","1","1"]], "status":"GAME_OVER"}`, res.Body.String())
}

func TestSelectNonMineNeighboorPoint(t *testing.T) {
	BoardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: map[int]models.Game{1: test_data.GetBoard(4, 3, []models.Point{models.Point{X: 4, Y: 1}})},
	})
	res := testAPICall(http.MethodPut, `/api/v1/boards/1`, `{"x":1, "y":1}`)
	assert.JSONEq(t, `{"id":1, "width":4, "height":3, "mines":1, "matrix":[[" "," ","1","-"],[" "," ","1","1"],[" "," "," "," "]], "status":"ACTIVE"}`, res.Body.String())
	res = testAPICall(http.MethodPut, `/api/v1/boards/1`, `{"x":4, "y":1}`)
	assert.JSONEq(t, `{"id":1, "width":4, "height":3, "mines":1, "matrix":[[" "," ","1","X"],[" "," ","1","1"],[" "," "," "," "]], "status":"GAME_OVER"}`, res.Body.String())
}

func TestSelectMinePointAndWin(t *testing.T) {
	BoardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: map[int]models.Game{1: test_data.GetBoard(4, 3, []models.Point{models.Point{X: 4, Y: 1}})},
	})
	BoardService.SelectPoint(1, models.Point{X: 1, Y: 1})
	res := testAPICall(http.MethodPut, `/api/v1/boards/1`, `{"x":4, "y":1, "mine":true}`)
	assert.JSONEq(t, `{"id":1, "width":4, "height":3, "mines":1, "matrix":[[" "," ","1","*"],[" "," ","1","1"],[" "," "," "," "]], "status":"WIN"}`, res.Body.String())
}

func TestSelectAlreadySelectedPoint(t *testing.T) {
	BoardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: map[int]models.Game{1: test_data.GetBoard(3, 3, []models.Point{models.Point{X: 3, Y: 1}})},
	})
	BoardService.SelectPoint(1, models.Point{X: 2, Y: 1})
	res := testAPICall(http.MethodPut, `/api/v1/boards/1`, `{"x":2, "y":1}`)
	assert.JSONEq(t, `{"error": "ALREADY_SELECTED", "status_code": 400}`, res.Body.String())
	assert.Equal(t, 400, res.Code)
}
