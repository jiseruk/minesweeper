package apis

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/models"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/services"
	"github.com/jiseruk/minesweeper/cmd/minesweeper/test_data"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var BoardService *services.BoardService

func init() {
	BoardService = services.NewBoardService(&test_data.MockBoardDAO{
		Records: make(map[int]models.Board, 0),
	})
	//BoardService = services.NewBoardService(&daos.BoardDAOImpl{})
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
	BoardService.Create(&board)
	c.JSON(http.StatusCreated, board.Print())
}

// SelectPoint godoc
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
	board, err := BoardService.SelectPoint(id, point)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, board.Print())
}

// Get godoc
// @Summary Returns the board
// @Produce json
// @Success 200 {object} models.Board
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", Health)
		v1.GET("/boards/:id", Get)
		v1.PUT("/boards/:id", SelectPoint)
		v1.POST("/boards/", CreateBoard)
	}
	return r
}
