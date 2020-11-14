// Copyright 2020 Anouar Fadili.

package main

import (
	"fmt"
	"math/rand"
	"strconv"
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

func (f *field) initMines(num int) {
	for i := 0; i < num; {
		x := rand.Intn(f.rows-2) + 2
		y := rand.Intn(f.cols-2) + 2
		if !f.boxes[x][y].isMine() {
			f.boxes[x][y].val = mine
			i++
		}
	}
}

func main() {
	f := newField(3, 3)
	f.make()
	f.addWalls()
	f.print()
}
