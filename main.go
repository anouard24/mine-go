// Copyright 2020 Anouar Fadili.

package main

import (
	"fmt"
	"strconv"
	"time"
)

// values from 0 to 8 indicate how much
// mines surround this box
type value int

// state of the box
type state int

const (
	// const box values
	clear value = 0
	wall  value = -1 //  wall box
	mine  value = 9  //  mine box

	// const box state
	hidden  state = 0
	open    state = 1
	suspect state = 5
	marked  state = 7
)

type point struct {
	x, y int64
}

type box struct {
	val    value
	status state
}

func (b *box) isMine() bool {
	return b.val == mine
}

func (b *box) isWall() bool {
	return b.val == wall
}

func (b *box) isHidden() bool {
	return b.status == hidden
}

func (b *box) show() string {
	switch b.val {
	case wall:
		return "+"
	case mine:
		return "@"
	case clear:
		return " "
	default:
		return strconv.Itoa(int(b.val))
	}
}

func (b *box) str() string {
	switch b.status {
	case marked:
		return "X"
	case suspect:
		return "S"
	case hidden:
		return "."
	default:
		return b.show()
	}
}

type field struct {
	rows, cols int
	boxes      [][]*box
}

func newField(rows, cols int) *field {
	f := new(field) // new pointer of type field
	// add 2 rows for top and bottom walls
	f.rows = rows + 2
	// add 2 cols for right and left walls
	f.cols = cols + 2
	return f
}

func (f *field) make() {
	f.boxes = make([][]*box, f.rows)
	for i := range f.boxes {
		f.boxes[i] = make([]*box, f.cols)
		for j := range f.boxes[i] {
			f.boxes[i][j] = new(box)
		}
	}
}

func (f *field) print() {
	for i := range f.boxes {
		for _, b := range f.boxes[i] {
			fmt.Printf("%2s ", b.str())
		}
		fmt.Printf("\n")
	}
}

func (f *field) addWalls() {
	// use the same wall box pointer
	// because wall boxes don't change state
	wallbox := &box{wall, open}
	for i := 0; i < f.rows; i++ {
		f.boxes[i][0] = wallbox
		f.boxes[i][f.cols-1] = wallbox
	}
	for j := 0; j < f.cols; j++ {
		f.boxes[0][j] = wallbox
		f.boxes[f.rows-1][j] = wallbox
	}
}

func (f *field) randomPoint() point {
	x := (time.Now().UnixNano() % int64(f.rows-2)) + 1
	y := (time.Now().UnixNano() % int64(f.cols-2)) + 1
	return point{x, y}
}

func (f *field) initMines(num int) {
	for i := 0; i < num; {
		p := f.randomPoint()
		if !f.boxes[p.x][p.y].isMine() {
			f.boxes[p.x][p.y].val = mine
			i++
		}
	}
}

func (f *field) calculateAdjacentsMines() {
	for i := 1; i < f.rows-1; i++ {
		for j := 1; j < f.cols-1; j++ {
			if !f.boxes[i][j].isMine() {
				var val value = 0
				for k := i - 1; k <= i+1; k++ {
					for l := j - 1; l <= j+1; l++ {
						if f.boxes[k][l].isMine() {
							val++
						}
					}
				}
				f.boxes[i][j].val = val
			}
		}
	}
}

func main() {
	f := newField(4, 4)
	f.make()
	f.initMines(4)
	f.calculateAdjacentsMines()
	f.addWalls()
	f.print()
}
