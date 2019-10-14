package apis

import (
	"bytes"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/services"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/test_data"

	//"github.com/jiseruk/minsweeper/cmd/minsweeper/config"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type apiTestCase struct {
	tag              string
	method           string
	urlToServe       string
	urlToHit         string
	body             string
	function         gin.HandlerFunc
	status           int
	responseFilePath string
}

// Creates new router in testing mode
func testingRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	//config.Config.DB = test_data.ResetDB()

	return router
}

// Used to run single API test case. It makes HTTP request and returns its response
func testAPICall(router *gin.Engine, method string, urlToServe string, urlToHit string, function gin.HandlerFunc, body string) *httptest.ResponseRecorder {
	router.Handle(method, urlToServe, function)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, urlToHit, bytes.NewBufferString(body))
	router.ServeHTTP(res, req)
	return res
}

// Used to run suite (list) of test cases. It checks JSON response is same as expected data in test case file.
// All test expected test case responses are stored in `test_data/test_case_data` folder in format `<suite_name>_t<number>.json`
func runTestsSuite(t *testing.T, tests []apiTestCase) {
	boardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: map[int]models.Board{1: models.Board{Model: models.Model{ID: 1}, Height: 5, Width: 5, Mines: 5, Status: models.StatusActive}},
	})
	for _, test := range tests {
		router := testingRouter()
		res := testAPICall(router, test.method, test.urlToServe, test.urlToHit, test.function, test.body)
		assert.Equal(t, test.status, res.Code, test.tag)
		if test.responseFilePath != "" {
			response, _ := ioutil.ReadFile(".." + test.responseFilePath)
			assert.JSONEq(t, string(response), res.Body.String(), test.tag)
		}
	}
}
