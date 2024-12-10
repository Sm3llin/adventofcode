package main

import (
	"adventofcode"
	"bytes"
	"fmt"
	"os"
)

type Position []int
type Positions []Position

func NewPosition(x, y int) Position {
	return Position{x, y}
}

type Map struct {
	locations map[rune]Positions

	Height int
	Width  int
}

func (m *Map) CheckBounds(x, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
}

func (m Map) CountResonantAntinodes() int {
	antinodeTable := make([][]int, m.Height)
	for i, _ := range antinodeTable {
		antinodeTable[i] = make([]int, m.Width)
	}

	for r, positions := range m.locations {
		fmt.Sprintf("r=%c, positions=%v\n", r, positions)

		deltas := positions.FindDeltas()
		for _, delta := range deltas {
			x, y := delta.FindDiff()
			//fmt.Printf("delta=%v, diff=%d, %d\n", delta, x, y)

			i := 0
			for m.CheckBounds(delta.a[0]-(x*i), delta.a[1]-(y*i)) {
				antinodeTable[delta.a[1]-(y*i)][delta.a[0]-(x*i)]++
				i++
			}

			i = 0
			for m.CheckBounds(delta.b[0]+(x*i), delta.b[1]+(y*i)) {
				antinodeTable[delta.b[1]+(y*i)][delta.b[0]+(x*i)]++
				i++
			}
		}
	}

	var uniqueAntinodes int
	for _, row := range antinodeTable {
		//fmt.Printf("%v\n", row)
		for _, count := range row {
			if count >= 1 {
				uniqueAntinodes++
			}
		}
	}
	return uniqueAntinodes
}

func (m Map) CountUniqueAntinodes() int {
	antinodeTable := make([][]int, m.Height)
	for i, _ := range antinodeTable {
		antinodeTable[i] = make([]int, m.Width)
	}

	for r, positions := range m.locations {
		fmt.Sprintf("r=%c, positions=%v\n", r, positions)

		deltas := positions.FindDeltas()
		for _, delta := range deltas {
			x, y := delta.FindDiff()
			//fmt.Printf("delta=%v, diff=%d, %d\n", delta, x, y)

			if m.CheckBounds(delta.a[0]-x, delta.a[1]-y) {
				antinodeTable[delta.a[1]-y][delta.a[0]-x]++
			}
			if m.CheckBounds(delta.b[0]+x, delta.b[1]+y) {
				antinodeTable[delta.b[1]+y][delta.b[0]+x]++
			}
		}
	}

	var uniqueAntinodes int
	for _, row := range antinodeTable {
		//fmt.Printf("%v\n", row)
		for _, count := range row {
			if count >= 1 {
				uniqueAntinodes++
			}
		}
	}
	return uniqueAntinodes
}

type Delta struct {
	a Position
	b Position
}

func (d Delta) FindDiff() (int, int) {
	return d.b[0] - d.a[0], d.b[1] - d.a[1]
}

func (p Positions) FindDeltas() []Delta {
	deltas := []Delta{}
	for a := 0; a < len(p)-1; a++ {
		for b := a + 1; b < len(p); b++ {
			deltas = append(deltas, Delta{a: p[a], b: p[b]})
		}
	}
	return deltas
}

func LoadLocations(data []byte) Map {
	positions := map[rune]Positions{}
	var x, y int
	var line []byte
	for y, line = range bytes.Split(data, []byte("\n")) {
		var r rune
		for x, r = range []rune(string(line)) {
			if r != '.' {
				if _, ok := positions[r]; !ok {
					positions[r] = Positions{}
				}
				positions[r] = append(positions[r], NewPosition(x, y))
			}
		}
	}
	return Map{
		locations: positions,

		Height: y + 1,
		Width:  x + 1,
	}
}

func main() {
	adventofcode.Time(func() {
		data, err := os.ReadFile("2024/day_eight/input.txt")
		if err != nil {
			panic(err)
		}

		m := LoadLocations(data)

		fmt.Printf("%d\n", m.CountUniqueAntinodes())
		fmt.Printf("%d\n", m.CountResonantAntinodes())
	})
}
