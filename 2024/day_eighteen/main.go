package main

import (
	"adventofcode"
	"adventofcode/toolbox/conversion"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/grid"
	"bytes"
	"fmt"
	"slices"
)

func convertInput(d []byte) []grid.Position {
	return conversion.To(bytes.Split(d, []byte("\n")), func(l []byte) grid.Position {
		s := bytes.Split(l, []byte(","))
		x, _ := conversion.ToInt(string(s[0]))
		y, _ := conversion.ToInt(string(s[1]))
		return grid.Position{
			X: x,
			Y: y,
		}
	})
}

func CreateGrid(in []byte, width, height int) (grid.Grid[string], []grid.Position) {
	positions := convertInput(in)

	g := grid.NewGridValue(".", width, height)

	return g, positions
}

func main() {
	adventofcode.Time(func() {
		g, positions := CreateGrid(fs.LoadFile("2024/day_eighteen/input.txt"), 71, 71)

		for i := range 1024 {
			p := positions[i]
			g.Set(p.X, p.Y, "#")
		}

		maze := grid.NewMaze(g, "#", func(from, to grid.Position, d grid.Direction) (allow bool, score int) {
			return true, 1
		})

		path, success := maze.Solve(grid.Position{}, grid.Position{X: 70, Y: 70})

		for _, p := range path {
			g.Set(p.X, p.Y, "X")
		}
		fmt.Println(g.Render())

		if !success {
			fmt.Println("did not get expected output")
		} else {
			fmt.Printf("%v\n", path)
			fmt.Printf("%d\n", len(path)-1)
		}
	})

	// Run until byte is in path
	// Then simulate the array again
	// might need new maze
	adventofcode.Time(func() {
		g, positions := CreateGrid(fs.LoadFile("2024/day_eighteen/input.txt"), 71, 71)

		for i := range 1024 {
			p := positions[i]
			g.Set(p.X, p.Y, "#")
		}

		var path []grid.Position
		for i := range len(positions) {
			p := positions[i]
			g.Set(p.X, p.Y, "#")

			if path != nil && !slices.Contains(path, p) {
				continue
			}

			maze := grid.NewMaze(g, "#", func(from, to grid.Position, d grid.Direction) (allow bool, score int) {
				return true, 1
			})
			var success bool
			path, success = maze.Solve(grid.Position{}, grid.Position{X: 70, Y: 70})

			if !success {

				fmt.Printf("position x=%d y=%d is the last byte\n", p.X, p.Y)
				break
			}
		}

	})
}
