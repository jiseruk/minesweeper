package apis

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/daos"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var BoardService *services.BoardService

func init() {
	/*BoardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: make(map[int]models.Game, 0),
	})*/
	BoardService = services.NewBoardService(&daos.BoardDAOImpl{})
}

// CreateBoard godoc
// @Summary Creates board based on given json data
// @Produce json
// @Success 201 {object} models.Game
// @Router /boards/ [post]
func CreateBoard(c *gin.Context) {

	var board models.Game
	if err := c.ShouldBindJSON(&board); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	BoardService.Create(&board)
	c.JSON(http.StatusCreated, board.Print())
}

// SelectPoint godoc
// @Summary Selects a point in the board
// @Produce json
// @Success 200 {object} models.Game
// @Router /boards/{id} [put]
func SelectPoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var point models.Point
	if err := c.ShouldBindJSON(&point); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	board, err := BoardService.SelectPoint(id, point)
	if err != nil {
		if gameErr, ok := err.(*models.GameError); ok {
			c.JSON(gameErr.StatusCode, &gameErr)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, board.Print())
}

// Get godoc
// @Summary Returns the board
// @Produce json
// @Success 200 {object} models.Game
// @Router /boards/{id} [get]
func Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	board, err := BoardService.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, board.Print())
}

// Get godoc
// @Summary Returns the board
// @Produce json
// @Success 200 {object} models.Game
// @Router /boards/{id}/board [get]
func GetBoard(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	board, err := BoardService.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, board.String())
}

func Health(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
func GetRouter() *gin.Engine {
	// Creates a router without any middleware by default
	r := gin.New()
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1.GET("/ping", Health)
	v1.GET("/boards/:id", Get)
	v1.GET("/boards/:id/board", GetBoard)
	v1.PUT("/boards/:id", SelectPoint)
	v1.POST("/boards/", CreateBoard)
	return r
}
