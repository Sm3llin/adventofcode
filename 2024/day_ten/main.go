package main

import (
	"adventofcode"
	"adventofcode/toolbox/conversion"
	"adventofcode/toolbox/fs"
	"bytes"
	"fmt"
)

type Direction []int

var (
	UP    = Direction{0, -1}
	DOWN  = Direction{0, 1}
	LEFT  = Direction{-1, 0}
	RIGHT = Direction{1, 0}
)

type Map struct {
	Tiles [][]Tile

	Width  int
	Height int
}

func (t Tile) String() string {
	// enables seeing increment of hits per travel chart
	return fmt.Sprintf("%d", t.Score)
}

func (m Map) Score() (int, int) {
	// Find the 10s
	tens := make([][]int, 0)
	ones := make([][]int, 0)

	for y, _ := range m.Tiles {
		for x, tile := range m.Tiles[y] {
			if tile.Value == 10 {
				tens = append(tens, []int{x, y})
			} else if tile.Value == 1 {
				ones = append(ones, []int{x, y})
			}
		}
	}

	for _, t := range tens {
		cheatSheet := make([][]int, m.Height)
		for i, _ := range cheatSheet {
			cheatSheet[i] = make([]int, m.Width)
		}

		m.March(t[0], t[1], UP, cheatSheet)
		m.March(t[0], t[1], DOWN, cheatSheet)
		m.March(t[0], t[1], LEFT, cheatSheet)
		m.March(t[0], t[1], RIGHT, cheatSheet)

		//fmt.Printf("Cheat Sheet for %v\n", t)
		//for j, line := range cheatSheet {
		//	fmt.Printf("%v %s\n", line, m.Tiles[j])
		//}
		//fmt.Printf("\n")
	}

	var score int
	var rating int
	for _, o := range ones {
		score += m.Tiles[o[1]][o[0]].Score
		rating += m.Tiles[o[1]][o[0]].Rating
	}

	return score, rating
}

func (m Map) March(x, y int, direction Direction, sheet [][]int) {
	value := m.Tiles[y][x].Value

	stepX, stepY := x+direction[0], y+direction[1]
	if !m.Contains(stepX, stepY) || value == 0 {
		return
	}
	onlyLeaveRating := sheet[stepY][stepX] == 1

	stepValue := m.Tiles[stepY][stepX].Value
	if stepValue != value-1 {
		return
	} else if stepValue == 1 {
		if !onlyLeaveRating {
			m.Tiles[stepY][stepX].Score += 1
		}
		m.Tiles[stepY][stepX].Rating += 1
		sheet[stepY][stepX] = 1
		return
	}

	sheet[stepY][stepX] = 1

	m.March(stepX, stepY, UP, sheet)
	m.March(stepX, stepY, DOWN, sheet)
	m.March(stepX, stepY, LEFT, sheet)
	m.March(stepX, stepY, RIGHT, sheet)
}

func (m Map) Contains(x, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
}

type Tile struct {
	// Value is plus 1 for the map
	Value, Score, Rating int
}

func LoadMap(data []byte) Map {
	lines := bytes.Split(data, []byte("\n"))
	m := Map{
		Tiles: make([][]Tile, len(lines)),
	}

	var x, y int
	var line []byte
	for y, line = range lines {
		m.Tiles[y] = make([]Tile, len(line))
		var b byte
		for x, b = range line {
			var v int
			if b == '.' {
				v = 0
			} else {
				v, _ = conversion.ToInt(b)
				v++
			}

			m.Tiles[y][x].Value = v
		}
	}

	m.Height = y + 1
	m.Width = x + 1

	return m
}

func main() {
	adventofcode.Time(func() {
		data := fs.LoadFile("2024/day_ten/input.txt")
		m := LoadMap(data)

		score, rating := m.Score()
		fmt.Printf("Part 1: %d, Part 2: %d\n", score, rating)
	})
}
