// Copyright 2020 Anouar Fadili.

package main

import (
	"fmt"
	"math"
	"math/rand"
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
	x, y int
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

func (b *box) isMarked() bool {
	return b.status == marked
}

func (b *box) isOpen() bool {
	return b.status == open
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
	rows, cols     int
	boxes          [][]*box
	numMines       int
	numCovered     int
	hints          int
	remainingMines int
}

func newField(rows, cols int) *field {
	f := new(field) // new pointer of type field
	f.numCovered = rows * cols
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

// print all boxes (for debug or in the end of a game)
func (f *field) printAll() {
	for i := range f.boxes {
		for _, b := range f.boxes[i] {
			fmt.Printf("%2s ", b.show())
		}
		fmt.Printf("\n")
	}
}

func (f *field) addWalls() {
	// use the same wall box pointer
	// because wall boxes don't change state
	wallBox := &box{wall, open}
	for i := 0; i < f.rows; i++ {
		f.boxes[i][0] = wallBox
		f.boxes[i][f.cols-1] = wallBox
	}
	for j := 0; j < f.cols; j++ {
		f.boxes[0][j] = wallBox
		f.boxes[f.rows-1][j] = wallBox
	}
}

func (f *field) randomPoint() point {
	x := rand.Intn(f.rows-2) + 1
	y := rand.Intn(f.cols-2) + 1
	return point{x, y}
}

func (f *field) initMines(num int) {
	f.numMines = num
	f.remainingMines = num
	for i := 0; i < num; {
		p := f.randomPoint()
		if !f.getBox(p).isMine() {
			f.getBox(p).val = mine
			i++
		}
	}
	f.hints = calculateHints(f.numCovered, f.numMines)
}

func calculateHints(numBoxes, numMines int) int {
	percentageMines := float64(numMines) / float64(numBoxes)
	percentageHints := 0.0092 * math.Exp(7.8*percentageMines)
	return int(math.Round(float64(numBoxes) * percentageHints))
}

func (f *field) calculateAdjacentMines() {
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

func (f *field) useHint(p point) {
	f.hints--
	if !f.uncoverBox(p) {
		// the hint uncover a mine!
		f.numCovered++
		f.remainingMines--
	}
}

// return false if uncover a mine
func (f *field) uncoverBox(p point) bool {
	curBox := f.getBox(p)
	if curBox.isWall() || !curBox.isHidden() {
		// "We can't uncover the box in position (%v, %v)", p.x, p.y
		return true
	}
	curBox.status = open
	f.numCovered--
	if curBox.isMine() {
		// "Ops! You uncovered a mine!"
		return false
	}
	if curBox.val == clear {
		// no mine surround this box
		for i := p.x - 1; i <= p.x+1; i++ {
			for j := p.y - 1; j <= p.y+1; j++ {
				f.uncoverBox(point{i, j})
			}
		}
	}
	return true
}

func (f *field) uncoverAdjacentBoxes(p point) bool {
	for i := p.x - 1; i <= p.x+1; i++ {
		for j := p.y - 1; j <= p.y+1; j++ {
			curPoint := point{i, j}
			if p != curPoint {
				if !f.uncoverBox(curPoint) {
					// We uncover a mine: game over
					return false
				}
			}
		}
	}
	return true
}

func (f *field) gameEnds() bool {
	// if we use a hint that uncover
	// a mine safely then:
	return f.numMines >= f.numCovered
}

func (f *field) getBox(p point) *box {
	return f.boxes[p.x][p.y]
}

func (f *field) toggleMarkMineWith(p point, status state) bool {
	curBox := f.getBox(p)
	if curBox.isWall() || curBox.isOpen() {
		return true
	}
	if curBox.isHidden() {
		if status != marked || f.remainingMines > 0 {
			// change status if it's suspect
			// or if it is marked and remaining mines > 0
			curBox.status = status
			if status == marked {
				f.remainingMines--
			}
		}
	} else {
		if curBox.isMarked() {
			f.remainingMines++
		}
		curBox.status = hidden
	}
	return true
}

func inRange(val int, min, max int) bool {
	return val >= min && val <= max
}

func (f *field) validPoint(p point) bool {
	return inRange(p.x, 1, f.rows-2) &&
		inRange(p.y, 1, f.cols-2)
}

func scanInput(nameVar string, min, max int) int {
	val := min - 1
	for !inRange(val, min, max) {
		for {
			fmt.Printf("Enter number of %s  [%d-%d]: ", nameVar, min, max)
			if _, err := fmt.Scanf("%d\n", &val); err == nil {
				break
			}
		}
	}
	return val
}

func input(printStr, format string, a ...interface{}) {
	for {
		fmt.Print(printStr)
		if _, err := fmt.Scanf(format, a...); err == nil {
			break
		}
	}
}

func (f *field) runAction(p point, cmd int) bool {
	if !f.validPoint(p) {
		return true
	}
	switch cmd {
	case 0:
		return f.uncoverBox(p)
	case 1:
		return f.toggleMarkMineWith(p, marked)
	case 2:
		return f.toggleMarkMineWith(p, suspect)
	case 3:
		return f.uncoverAdjacentBoxes(p)
	case 4:
		f.useHint(p)
		return true
	default:

	}
	return true
}

func main() {
	rand.Seed(time.Now().UnixNano())
	rows := scanInput("rows", 4, 15)
	cols := scanInput("cols", 5, 20)
	// number of mines between 10% and 50% of total boxes
	mines := scanInput("mines", rows*cols/10, rows*cols/2)
	f := newField(rows, cols)
	f.make()
	f.initMines(mines)
	f.calculateAdjacentMines()
	f.addWalls()
	var p point
	var cmd int
	for !f.gameEnds() {
		f.print()
		fmt.Printf("You have %d 🔮 hints left\n", f.hints)
		fmt.Printf("You have %d \U0001F9E8 mines left\n", f.remainingMines)
		input("Enter x and y: ", "%d %d\n", &p.x, &p.y)
		input(
			`Enter action number:
	0 to uncover
	1 to mark
	2 to suspect
	3 to uncover all adjacent boxes
	4 to use hint (safe uncover)
>> `,
			"%d\n", &cmd)
		if !f.runAction(p, cmd) {
			// lose game
			f.printAll()
			fmt.Println("💥 Ops! Game Over...")
			return
		}
	}
	f.printAll()
	fmt.Println("🤓 Great! You Win!")
}
