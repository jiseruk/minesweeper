package models

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

const StatusActive = "ACTIVE"
const StatusWin = "WIN"
const StatusLost = "LOST"

type Model struct {
	ID        int        `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type Point struct {
	X             int  `json:"x"`
	Y             int  `json:"y"`
	Mine          bool `json:"mine"`
	selected      bool `json:"selected"`
	mineCandidate bool `json:"mine_candidate"`
	mineNeigboors int  `json:"mine_neigboors"`
}

type Board struct {
	Model
	Width   int        `gorm:"column:width" json:"width"`
	Height  int        `gorm:"column:height" json:"height"`
	Mines   int        `gorm:"column:mines" json:"mines"`
	Matrix  [][]string `gorm:"-" json:"matrix"`
	Board   [][]Point  `gorm:"-" json:"-"`
	BoardDB string     `gorm:"board" json:"-"`
	Status  string     `json:"status"`
}

func (p *Point) click() {
	p.selected = true
}

func (p *Point) unselect() {
	p.selected = false
}

func (p *Point) mark() {
	p.mineCandidate = true
}

func (p *Point) unmark() {
	p.mineCandidate = false
}

func (b *Board) Print() *Board {
	//out := Board{Width: b.Width, Height: b.Height, Mines: b.Mines, Model: Model{ID: b.ID}, Status: b.Status}
	b.Matrix = make([][]string, b.Height)
	for i, points := range b.Board {
		b.Matrix[i] = make([]string, b.Width)
		for j := range points {
			if b.Status == StatusLost {
				if b.Board[i][j].Mine {
					b.Matrix[i][j] = "*"
				} else {
					b.Matrix[i][j] = strconv.Itoa(b.Board[i][j].mineNeigboors)
				}
			} else if !points[j].selected {
				b.Matrix[i][j] = "-"
			} else if points[j].Mine && points[j].mineCandidate {
				b.Matrix[i][j] = "*"
			} else if points[j].Mine && !points[j].mineCandidate {
				b.Matrix[i][j] = "X"
			} else {
				if points[j].mineNeigboors == 0 {
					b.Matrix[i][j] = " "
				} else {
					b.Matrix[i][j] = strconv.Itoa(points[j].mineNeigboors)
				}
			}
			//looping over each element of array and assigning it a random variable
		}
	}
	return b
}

func (b *Board) Select(p Point) bool {
	if b.Status != StatusActive {
		return false
	}
	var point *Point
	point = &b.Board[p.X][p.Y]
	if point.selected || point.mineCandidate {
		return false
	} else if point.Mine {
		b.Status = StatusLost
		return false
	} else {
		point.selected = true
		defer b.setBoardForDB()
		for _, pt := range b.neigboors(point) {
			if !pt.selected && !pt.mineCandidate && !pt.Mine {
				if pt.mineNeigboors == 0 {
					b.Select(*pt)
				} else {
					b.Board[pt.X][pt.Y].selected = true
					defer b.setBoardForDB()
				}
			}
		}
		return true
	}
}

func (b *Board) Populate() *Board {
	b.Board = make([][]Point, b.Height)
	for i, points := range b.Board {
		points = make([]Point, b.Width)
		for j := range points {
			points[j] = Point{X: i, Y: j, Mine: false, selected: false, mineCandidate: false}
			//looping over each element of array and assigning it a random variable
		}
		b.Board[i] = points
	}
	i := 0
	for i < b.Mines {
		x := rand.Intn(b.Width)
		y := rand.Intn(b.Height)
		if !b.Board[x][y].Mine {
			b.Board[x][y] = Point{X: x, Y: y, Mine: true, selected: false, mineCandidate: false}
			i++
		}
	}

	for _, points := range b.Board {
		for j := range points {
			if !points[j].Mine {
				points[j].mineNeigboors = len(filter(b.neigboors(&points[j]), mineNeigboors))
			}
		}
	}
	b.setBoardForDB()
	return b
}

func (b *Board) neigboors(point *Point) []*Point {
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

func (b *Board) setBoardForDB() {
	jsonBoard, _ := json.Marshal(b.Board)
	b.BoardDB = string(jsonBoard)
}

func filter(points []*Point, test func(*Point) bool) (ret []*Point) {
	for _, p := range points {
		if test(p) {
			ret = append(ret, p)
		}
	}
	return
}

var mineNeigboors = func(p *Point) bool { return p.Mine }
