package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type GameStatus int

const (
	StatusActive GameStatus = iota
	StatusWin
	StatusLost
)

var StatusNames = [...]string{"ACTIVE", "WIN", "GAME_OVER"}
var StatusValues = map[string]GameStatus{
	StatusNames[0]: StatusActive,
	StatusNames[1]: StatusWin,
	StatusNames[2]: StatusLost,
}

func (s GameStatus) String() string {
	return StatusNames[s]
}
func (s GameStatus) MarshalJSON() ([]byte, error) {
	val := fmt.Sprintf(`"%s"`, s.String())
	return []byte(val), nil
}

type Model struct {
	ID        int        `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"-"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`
}

type Point struct {
	X              int   `json:"x"`
	Y              int   `json:"y"`
	Mine           *bool `json:"mine,omitempty"`
	Selected       bool  `json:"selected"`
	MineCandidate  bool  `json:"mine_candidate"`
	MineNeighbours int   `json:"mine_neighbours"`
}

type Board [][]Point
type Game struct {
	Model
	Width            int        `gorm:"column:width" json:"width"`
	Height           int        `gorm:"column:height" json:"height"`
	Mines            int        `gorm:"column:mines" json:"mines"`
	UnknownMines     int        `gorm:"unknown_mines" json:"-"`
	UnselectedPoints int        `gorm:"unselected_points" json:"-"`
	Matrix           [][]string `gorm:"-" json:"matrix"`
	Board            Board      `gorm:"-" json:"-"`
	BoardDB          string     `gorm:"board" json:"-"`
	Status           GameStatus `gorm:"status" json:"status"`
}

type GameError struct {
	StatusCode int    `json:"status_code"`
	ErrorMsg   string `json:"error"`
}

func (e *GameError) Error() string {
	return e.ErrorMsg
}

func (p *Point) click() {
	p.Selected = true
}
func (b *Game) String() string {
	b.Print()
	s := ""
	for x := range b.Matrix {
		for y := range b.Matrix[x] {
			s += b.Matrix[x][y]
			if y == b.Height-1 {
				s += "\n"
			} else {
				s += " "
			}
		}
	}
	return s
}

func (b *Game) Print() *Game {
	b.Matrix = make([][]string, b.Height)
	for i := range b.Matrix {
		b.Matrix[i] = make([]string, b.Width)
	}
	for i, points := range b.Board {
		for j := range points {
			if b.Status != StatusActive {
				if *b.Board[i][j].Mine {
					if !b.Board[i][j].MineCandidate {
						b.Matrix[j][i] = "X"
					} else {
						b.Matrix[j][i] = "*"
					}

				} else {
					b.Matrix[j][i] = strconv.Itoa(b.Board[i][j].MineNeighbours)
				}
			} else if b.Board[i][j].MineCandidate {
				b.Matrix[j][i] = "*"
			} else if !b.Board[i][j].Selected {
				b.Matrix[j][i] = "-"
			} else {
				b.Matrix[j][i] = strconv.Itoa(b.Board[i][j].MineNeighbours)
			}
			//looping over each element of array and assigning it a random variable
		}
	}
	return b
}

func (b *Game) SetMine(p Point) error {
	if b.Status != StatusActive {
		return &GameError{StatusCode: 400, ErrorMsg: "ALREADY_FINISHED"}
	} else if p.Mine == nil {
		return &GameError{StatusCode: 400, ErrorMsg: "INVALID"}
	}
	var point *Point
	point = &b.Board[p.X][p.Y]
	fmt.Printf("%#v", point)
	if point.Selected {
		return &GameError{StatusCode: 400, ErrorMsg: "ALREADY_SELECTED"}
	} else if point.MineCandidate {
		if *p.Mine {
			return &GameError{StatusCode: 400, ErrorMsg: "ALREADY_MINE_CANDIDATE"}
		} else {
			point.MineCandidate = false
			b.UnknownMines++
		}
	} else if !point.MineCandidate {
		if !*p.Mine {
			return &GameError{StatusCode: 400, ErrorMsg: "ALREADY_NOT_MINE_CANDIDATE"}
		} else {
			point.MineCandidate = true
			b.UnknownMines--
			if b.UnselectedPoints == 0 && b.UnknownMines == 0 {
				b.Status = StatusWin
			}
		}
	}

	b.setBoardForDB()
	return nil
}

func (b *Game) Select(p Point) error {
	if b.Status != StatusActive {
		return &GameError{StatusCode: 400, ErrorMsg: "ALREADY_FINISHED"}
	}
	var point *Point
	point = &b.Board[p.X][p.Y]
	if point.Selected || point.MineCandidate {
		return &GameError{StatusCode: 400, ErrorMsg: "ALREADY_SELECTED"}
	} else if *point.Mine {
		b.Status = StatusLost
		b.setBoardForDB()
		return nil
	} else {
		point.click()
		b.UnselectedPoints--
		if b.UnselectedPoints == 0 && b.UnknownMines == 0 {
			b.Status = StatusWin
		}
		if point.MineNeighbours > 0 {
			b.setBoardForDB()
			return nil
		}
		for _, pt := range b.neighbours(point) {
			if !pt.Selected && !pt.MineCandidate && !*pt.Mine {
				if pt.MineNeighbours == 0 {
					b.Select(*pt)
				} else {
					b.Board[pt.X][pt.Y].click()
					b.UnselectedPoints--
					if b.UnselectedPoints == 0 && b.UnknownMines == 0 {
						b.Status = StatusWin
					}
				}
			}
		}
		b.setBoardForDB()
		return nil
	}
}

func (b *Game) Populate(mines ...Point) *Game {
	b.Board = make([][]Point, b.Width)
	for i, points := range b.Board {
		points = make([]Point, b.Height)
		for j := range points {
			points[j] = Point{X: i, Y: j, Mine: BoolValue(false), Selected: false, MineCandidate: false}
		}
		b.Board[i] = points
	}

	for _, point := range mines {
		b.Board[point.X-1][point.Y-1].Mine = BoolValue(true)
	}
	i := 0
	for len(mines) == 0 && i < b.Mines {
		x := rand.Intn(b.Width)
		y := rand.Intn(b.Height)
		if !*b.Board[x][y].Mine {
			b.Board[x][y] = Point{X: x, Y: y, Mine: BoolValue(true), Selected: false, MineCandidate: false}
			i++
		}
	}

	for _, points := range b.Board {
		for j := range points {
			if !*points[j].Mine {
				points[j].MineNeighbours = len(filter(b.neighbours(&points[j]), mineNeigboors))
			}
		}
	}
	b.setBoardForDB()
	return b
}

func (b *Game) neighbours(point *Point) []*Point {
	points := make([]*Point, 0)
	if point.X > 0 && point.Y > 0 {
		points = append(points, &b.Board[point.X-1][point.Y-1])
	}
	if point.Y > 0 {
		points = append(points, &b.Board[point.X][point.Y-1])
	}
	if point.Y > 0 && point.X < b.Width-1 {
		points = append(points, &b.Board[point.X+1][point.Y-1])
	}
	if point.X > 0 {
		points = append(points, &b.Board[point.X-1][point.Y])
	}
	if point.X < b.Width-1 {
		points = append(points, &b.Board[point.X+1][point.Y])
	}
	if point.X > 0 && point.Y < b.Height-1 {
		points = append(points, &b.Board[point.X-1][point.Y+1])
	}
	if point.Y < b.Height-1 {
		points = append(points, &b.Board[point.X][point.Y+1])
	}
	if point.X < b.Width-1 && point.Y < b.Height-1 {
		points = append(points, &b.Board[point.X+1][point.Y+1])
	}
	return points
}

func (b *Game) setBoardForDB() {
	jsonBoard, _ := json.Marshal(b.Board)
	b.BoardDB = string(jsonBoard)
}

func (b *Game) GetBoardFromDB() {
	json.Unmarshal([]byte(b.BoardDB), &b.Board)
}

func filter(points []*Point, test func(*Point) bool) (ret []*Point) {
	for _, p := range points {
		if test(p) {
			ret = append(ret, p)
		}
	}
	return
}

var mineNeigboors = func(p *Point) bool { return *p.Mine }

func BoolValue(b bool) *bool {
	return &b
}
