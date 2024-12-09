package main

import (
	"math"
	"os"
	"strings"
)

type Puzzle [][]rune

func (p Puzzle) FindX(needle string) int {
	needles := []rune(needle)

	if len(needles)%2 == 0 {
		panic("Need an odd number of characters")
	}

	centre := int(math.Ceil(float64(len(needles) / 2)))
	start := needles[centre]

	needleWithSuffix := []rune(" " + needle)

	var found int
	// find the first character x,y
	for y := range p {
		for x := range p[y] {
			if start == p[y][x] {
				if (p.Search(x-centre-1, y-centre-1, 1, 1, needleWithSuffix) == 1 || p.Search(x+centre+1, y+centre+1, -1, -1, needleWithSuffix) == 1) && (p.Search(x-centre-1, y+centre+1, 1, -1, needleWithSuffix) == 1 || p.Search(x+centre+1, y-centre-1, -1, 1, needleWithSuffix) == 1) {
					found++
				}
			}
		}
	}

	return found
}

func (p Puzzle) Find(needle string) int {
	needles := []rune(needle)
	start := needles[0]

	var found int
	// find the first character x,y
	for y := range p {
		for x := range p[y] {
			if start == p[y][x] {
				// horizontal
				found += p.Search(x, y, 1, 0, needles)
				found += p.Search(x, y, -1, 0, needles)
				// vertical
				found += p.Search(x, y, 0, 1, needles)
				found += p.Search(x, y, 0, -1, needles)
				// diagonal
				found += p.Search(x, y, 1, 1, needles)
				found += p.Search(x, y, -1, -1, needles)
				found += p.Search(x, y, 1, -1, needles)
				found += p.Search(x, y, -1, 1, needles)
			}
		}
	}

	return found
}

func (p Puzzle) Size() (int, int) {
	for _, line := range p {
		return len(line), len(p)
	}
	return 0, 0
}

func (p Puzzle) Search(x, y, xStep, yStep int, needles []rune) int {
	boundX, boundY := p.Size()

	for i := 1; i < len(needles); i++ {
		x += xStep
		y += yStep

		// ensure within bounds
		if x < 0 || x >= boundX || y < 0 || y >= boundY {
			return 0
		}

		if p[y][x] != needles[i] {
			return 0
		}
	}
	return 1
}

func NewPuzzle(in string) Puzzle {
	puzzle := [][]rune{}
	for _, line := range strings.Split(in, "\n") {
		puzzlePieces := []rune(line)

		if len(puzzlePieces) != 0 {
			puzzle = append(puzzle, puzzlePieces)
		}
	}
	return puzzle
}

func main() {
	data, err := os.ReadFile("2024/day_four/input.txt")
	if err != nil {
		panic(err)
	}
	puzzle := NewPuzzle(string(data))

	println(puzzle.Find("XMAS"))
	println(puzzle.FindX("MAS"))
}
