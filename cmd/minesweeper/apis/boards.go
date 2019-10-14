package apis

import (
	"net/http"
	"strconv"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/daos"

	"github.com/gin-gonic/gin"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/services"
)

var boardService *services.BoardService

func init() {
	/*boardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: make(map[int]models.Board, 0),
	})*/
	boardService = services.NewBoardService(&daos.BoardDAOImpl{})
}

// CreateBoard godoc
// @Summary Creates board based on given json data
// @Produce json
// @Success 201 {object} models.Board
// @Router /boards/ [post]
func CreateBoard(c *gin.Context) {

	var board models.Board
	if err := c.ShouldBindJSON(&board); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	boardService.Create(&board)
	c.JSON(http.StatusCreated, board.Print())
}

// CreateBoard godoc
// @Summary Selects a point in the board
// @Produce json
// @Success 200 {object} models.Board
// @Router /boards/{id} [put]
func SelectPoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var point models.Point
	if err := c.ShouldBindJSON(&point); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	point.X--
	point.Y--
	board, err := boardService.SelectPoint(id, point)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, board.Print())
}
