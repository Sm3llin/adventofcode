package main

import (
	"adventofcode"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/grid"
	"bytes"
	"fmt"
	"log/slog"
	"math"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelWarn)
	adventofcode.Time(func() {
		fmt.Println("Starting...")
		n := solve(fs.LoadFile("2024/day_twenty/input.txt"), 2, 100)
		// 1.23ms
		fmt.Printf("Part 1: %d\n", n)
	})
	adventofcode.Time(func() {
		fmt.Println("Starting...")
		n := solve(fs.LoadFile("2024/day_twenty/input.txt"), 20, 100)
		// 1.23ms
		fmt.Printf("Part 2: %d\n", n)
	})
}

// Run a flood fill on a maze to determine the deltas between walls
func solve(data []byte, cheat int, countPico int) int {
	g := grid.NewGrid(bytes.Split(data, []byte("\n")))

	// find and replace start/end position
	start, _ := g.FindAndReplace(func(v byte) bool {
		return v == 'S'
	}, '.')
	end, _ := g.FindAndReplace(func(v byte) bool {
		return v == 'E'
	}, '.')

	m := grid.NewMaze(g, '#', func(from, to grid.Position) (allow bool, score int) {
		return true, 1
	})

	// extract a grid that will show how far a cell is from the end
	flood := m.FloodFill(end)
	path, success := m.Solve(start, end)

	if !success {
		panic("did not get expected output")
	}

	savingsTracker := map[int]int{}
	for _, step := range path {
		currentValue, err := flood.Get(step.X, step.Y)
		if err != nil {
			panic(err)
		}

		// each step we determine what type of shortcuts it can take
		// create a function that will take the grid, flood and start to then find the fastest shortcut
		values := SearchForShortcut(step, currentValue, cheat, g, flood)

		for _, shortcutSavings := range values {
			if shortcutSavings > 0 && shortcutSavings >= countPico {
				savingsTracker[shortcutSavings]++
			}
		}
	}

	var t int
	for k, v := range savingsTracker {
		if k >= countPico {
			fmt.Printf("saving %d %d\n", k, v)
			t += v
		}
	}

	for pos, currentValue := range g.Around(grid.Position{X: 1, Y: 3}, cheat) {
		if currentValue == '#' {
			g.Set(pos.X, pos.Y, '%')
		} else {
			g.Set(pos.X, pos.Y, '*')
		}
	}

	fmt.Println(m.RenderFunc(func(v byte) string {
		return fmt.Sprintf("%c", v)
	}))
	return t
}

type searchGrid struct {
	pos  grid.Position
	step int
}

func SearchForShortcut(start grid.Position, startScore, steps int, g grid.Grid[byte], flood grid.Grid[int]) []int {
	scores := []int{}
	for pos, option := range g.Around(start, steps) {
		if option != '.' {
			continue
		}

		distance := int(math.Abs(float64(start.X-pos.X)) + math.Abs(float64(start.Y-pos.Y)))
		score, _ := flood.Get(pos.X, pos.Y)

		savings := (startScore - score) - distance

		if distance > steps {
			continue
		}

		if savings > 0 {
			slog.Info("adding score", "savings", savings, "from", start, "pos", pos)
			scores = append(scores, savings)
		}
	}
	slog.Info("found scores", "pos", start, "scores", len(scores))

	return scores
}
