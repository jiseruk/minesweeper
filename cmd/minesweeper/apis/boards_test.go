package apis

import (
	"net/http"
	"testing"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/test_data"
)

func TestBoard(t *testing.T) {
	path := test_data.GetTestCaseFolder()
	runTestsSuite(t, []apiTestCase{
		{"t1 - create a board", "POST", "/boards/", "/boards/", `{"width":7, "height":7, "mines":5}`, CreateBoard, http.StatusCreated, path + "/board_tc1.json"},
		{"t2 - select point", "PUT", "/boards/:id", "/boards/2", `{"x":1, "y":2}`, SelectPoint, http.StatusOK, path + "/board_tc1.json"},
	})
}
