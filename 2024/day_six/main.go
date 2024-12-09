package main

import (
	"adventofcode"
	"bytes"
	"fmt"
	"os"
	"slices"
)

type Direction []int

var (
	UP    = Direction{0, -1}
	DOWN  = Direction{0, 1}
	LEFT  = Direction{-1, 0}
	RIGHT = Direction{1, 0}

	GuardUp    = '^'
	GuardDown  = 'v'
	GuardLeft  = '<'
	GuardRight = '>'
	Empty      = '.'
	Wall       = '#'

	Guard = []rune{GuardUp, GuardDown, GuardLeft, GuardRight}

	GuardDirection = map[rune]Direction{
		GuardUp:    UP,
		GuardDown:  DOWN,
		GuardLeft:  LEFT,
		GuardRight: RIGHT,
	}
)

func (d Direction) Equal(b Direction) bool {
	if len(d) != len(b) {
		return false
	}

	for i := range d {
		if d[i] != b[i] {
			return false
		}
	}
	return true
}

func NextDirection(direction Direction) Direction {
	if direction.Equal(UP) {
		return RIGHT
	} else if direction.Equal(LEFT) {
		return UP
	} else if direction.Equal(DOWN) {
		return LEFT
	} else if direction.Equal(RIGHT) {
		return DOWN
	}
	return direction
}

func (d Direction) Rune() rune {
	if d.Equal(UP) {
		return GuardUp
	} else if d.Equal(LEFT) {
		return GuardLeft
	} else if d.Equal(DOWN) {
		return GuardDown
	} else if d.Equal(RIGHT) {
		return GuardRight
	}
	return Empty
}

type Map struct {
	Width         int
	Height        int
	Grid          [][]byte
	WalkDirection [][]byte
}

func NewMap(data []byte) Map {
	var grid [][]byte
	for _, line := range bytes.Split(data, []byte("\n")) {
		grid = append(grid, line)
	}

	// Clone the grid to ensure modifications do not affect the original grid
	originalGrid := make([][]byte, len(grid))
	for i := range grid {
		originalGrid[i] = make([]byte, len(grid[i]))

		for j := range grid[i] {
			originalGrid[i][j] = byte(Empty)
		}
	}

	return Map{
		Width:         len(grid[0]),
		Height:        len(grid),
		Grid:          grid,
		WalkDirection: originalGrid,
	}
}

func (m Map) Get(x, y int) rune {
	return rune(m.Grid[y][x])
}

func (m Map) CheckBounds(x, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
}

func (m Map) Solve() (int, int) {
	guardX, guardY := -1, -1
	var guardDirection Direction

	// locate the guard
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if slices.Contains(Guard, m.Get(x, y)) {
				guardX, guardY = x, y
				guardDirection = GuardDirection[m.Get(x, y)]
				break
			}
		}
		if guardX != -1 && guardY != -1 {
			break
		}
	}

	if guardX == -1 && guardY == -1 {
		panic(fmt.Sprintf("Guard: x=%d, y=%d", guardX, guardY))
	}

	blockageChance := 0
	valid := 0
	for m.CheckBounds(guardX, guardY) {
		location := m.Get(guardX, guardY)

		// is this a unique walking location
		if location != 'X' {
			valid++
			m.Grid[guardY][guardX] = 'X'
		}

		// ...
		// .#.
		// .^.
		if m.CheckBounds(guardX+guardDirection[0], guardY+guardDirection[1]) {
			nextLocation := m.Get(guardX+guardDirection[0], guardY+guardDirection[1])
			for nextLocation == Wall {
				guardDirection = NextDirection(guardDirection)
				nextLocation = m.Get(guardX+guardDirection[0], guardY+guardDirection[1])
			}

			blockDirection := guardDirection
			blockX, blockY := guardX+blockDirection[0], guardY+blockDirection[1]

			// cannot place block in traveled location
			if m.Get(blockX, blockY) != 'X' {
				// start at the guards next position
				for x, y := guardX, guardY; true; x, y = x+blockDirection[0], y+blockDirection[1] {
					if !m.CheckBounds(x, y) {
						break
					}
					if m.WalkDirection[y][x] == byte(blockDirection.Rune()) {
						//fmt.Printf("Loop detected at x=%d,y=%d, gx=%d, gy=%d\n", x, y, guardX, guardY)
						blockageChance++
						break
					}

					// leaves a shadow
					m.WalkDirection[y][x] = byte(blockDirection.Rune())
					var shouldBreak bool
					for step := 0; m.CheckBounds(x+blockDirection[0], y+blockDirection[1]) && (m.Grid[y+blockDirection[1]][x+blockDirection[0]] == byte(Wall) || (y+blockDirection[1] == blockY && x+blockDirection[0] == blockX)); step++ {
						blockDirection = NextDirection(blockDirection)

						if step >= 5 {
							// stuck in a small box
							blockageChance++
							shouldBreak = true
							break
						}
					}
					if shouldBreak {
						break
					}
				}

				// reset walk direction
				for y := 0; y < m.Height; y++ {
					for x := 0; x < m.Width; x++ {
						m.WalkDirection[y][x] = byte(Empty)
					}
				}
			}
		}

		// walk a step
		guardX += guardDirection[0]
		guardY += guardDirection[1]
	}

	return valid, blockageChance
}

func main() {
	adventofcode.Time(func() {
		data, err := os.ReadFile("2024/day_six/input.txt")

		if err != nil {
			panic(err)
		}

		m := NewMap(data)

		fmt.Println(m.Solve())
	})
}
