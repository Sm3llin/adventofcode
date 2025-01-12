package main

import (
	"adventofcode"
	"adventofcode/toolbox/arrays"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/grid"
	"bytes"
	"fmt"
	"log/slog"
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

	// get the number at start
	//total, err := flood.Get(start.X, start.Y)
	//if err != nil {
	//	panic(err)
	//}

	// TODO: we could walk the best path and then do a step walk down the walls and check everywall in range for its
	//       neighbours (minusing 1 for distance)
	savingsTracker := map[int]int{}
	for _, step := range path {
		currentValue, err := flood.Get(step.X, step.Y)
		if err != nil {
			panic(err)
		}

		// each step we determine what type of shortcuts it can take
		// create a function that will take the grid, flood and start to then find the fastest shortcut
		values := SearchForShortcut(step, currentValue, cheat, g, flood)

		for _, value := range values {
			if value < 0 {
				continue
			}
			shortcutSavings := value

			if shortcutSavings > 0 && shortcutSavings >= countPico {
				savingsTracker[shortcutSavings]++
			}
		}
	}

	var t int
	for k, v := range savingsTracker {
		fmt.Printf("saving %d %d\n", k, v)

		if k >= countPico {
			t += v
		}
	}

	return t
}

type searchGrid struct {
	pos  grid.Position
	step int
}

func SearchForShortcut(start grid.Position, startScore, steps int, g grid.Grid[byte], flood grid.Grid[int]) []int {
	steps -= 1
	// walk the edges of the "maze" for the step distance
	finalWalls := arrays.Queue[searchGrid]{}
	mazeWalls := arrays.Queue[searchGrid]{}
	consideredWalls := arrays.Queue[grid.Position]{}

	for pos, adj := range g.Neighbours(start.X, start.Y, grid.ConnectedDirections) {
		if adj == '#' {
			slog.Info("seeding initial matches", "pos", pos, "step", 1)
			mazeWalls.Push(searchGrid{pos: pos, step: 1})
			finalWalls.Push(searchGrid{pos: pos, step: 1})
			consideredWalls.Push(pos)
		}
	}

	for prevPos := range mazeWalls.Iter() {
		if prevPos.step >= steps {
			continue
		}

		for pos, adj := range g.Neighbours(prevPos.pos.X, prevPos.pos.Y, grid.ConnectedDirections) {
			// check for a bounce back and don't add
			if adj == '#' && !consideredWalls.Exists(pos) {
				consideredWalls.Push(pos)
				finalWalls.Push(searchGrid{pos: pos, step: prevPos.step + 1})
				mazeWalls.Push(searchGrid{pos: pos, step: prevPos.step + 1})
			}
		}
	}

	// redecorate and print the considered blocks
	//clone := g.Clone()
	//
	//clone.Set(start.X, start.Y, 'S')
	//for wall := range consideredWalls.Iter() {
	//	clone.Set(wall.X, wall.Y, '%')
	//}
	//slog.Info(clone.RenderFunc(func(v byte) string {
	//	return fmt.Sprintf("%c", v)
	//}))

	// store the ending positions to find the best saving for that square
	found := map[string]int{}
	scores := []int{}
	for wall := range finalWalls.Iter() {
		for pos, adj := range g.Neighbours(wall.pos.X, wall.pos.Y, grid.ConnectedDirections) {
			if adj == '.' {
				score, err := flood.Get(pos.X, pos.Y)
				if err != nil {
					continue
				}
				nextScore := (startScore - score) - (wall.step + 1)

				// need to ignore when start and end is the same
				prevScore, ok := found[fmt.Sprintf("%d|%d", pos.X, pos.Y)]
				if ok && prevScore >= nextScore {
					slog.Info("ignoring score", "pos", pos, "score", prevScore, "nextScore", nextScore)
					continue
				}

				// need to add in the steps to see if it is worth it
				if nextScore > 0 {
					found[fmt.Sprintf("%d|%d", pos.X, pos.Y)] = nextScore
					scores = append(scores, nextScore)
					slog.Info("setting cheats", "pos", pos, "start", startScore, "score", score, "step", wall.step, "saving", nextScore)
				}
			}
		}
	}

	return scores
}
