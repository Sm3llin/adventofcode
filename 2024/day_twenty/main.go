package main

import (
	"adventofcode"
	"adventofcode/toolbox/arrays"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/grid"
	"bytes"
	"fmt"
	"slices"
)

func main() {
	adventofcode.Time(func() {
		n := solve(fs.LoadFile("2024/day_twenty/input.txt"), 2, 100)
		// 1.23ms
		fmt.Printf("Part 1: %d\n", n)
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
	total, err := flood.Get(start.X, start.Y)
	if err != nil {
		panic(err)
	}

	// TODO: we could walk the best path and then do a step walk down the walls and check everywall in range for its
	//       neighbours (minusing 1 for distance)
	savingsTracker := map[int]int{}
	for _, step := range path {
		// each step we determine what type of shortcuts it can take
		// create a function that will take the grid, flood and start to then find the fastest shortcut
		value := SearchForShortcut(step, countPico, g, flood)

		savingsTracker[total-value]++
	}

	// potentially switch to solve the maze and selecting all cell within reach
	//x=0

	savingsTracker = map[int]int{}
	var z int
	_ = total
	// need to find the lowest path within 2 blocks, part 2 will be multiple steps
	for p, cell := range g.All() {
		if cell == '#' {
			continue
		}

		bestSaving := -1
		// select a grid within cheat by cheat
		for x := -cheat; x <= cheat; x++ {
			for y := -cheat; y <= cheat; y++ {
				if (x == 0 && y == 0) || !g.CheckBounds(p.X+x, p.Y+y) {
					continue
				}

				// we need to track the destination cells that require a wall hop
				deltaX, deltaY := p.Delta(grid.Position{
					X: x,
					Y: y,
				})

				if deltaX+deltaY > cheat {
					continue
				}

				cheatValue, _ := flood.Get(p.X+x, p.Y+y)
				if cheatValue > total {
					continue
				}

				if cheatValue > bestSaving {
					bestSaving = cheatValue
				}
			}
		}

		if bestSaving >= countPico {
			savingsTracker[bestSaving]++
			z++
		}
		value, _ := flood.Get(p.X, p.Y)

		type learn struct {
			p grid.Position
			v int
		}
		stepped := []grid.Position{}
		cellQueue := arrays.NewQueue([]learn{
			{p, 0},
		})

		bestSaving = -1
		// walk steps for each cell and in the number of cheats adding each to the queue to be explored
		for learner := range cellQueue.Iter() {
			stepped = append(stepped, learner.p)

			if learner.v > cheat {
				continue
			}

			// if I am a wall ignore me for checking score
			currentCell, err := g.Get(learner.p.X, learner.p.Y)
			if err != nil {
				continue
			} else if currentCell != '#' {
				cheatValue, err := flood.Get(learner.p.X, learner.p.Y)
				if err == nil {
					if cheatValue > value {
						continue
					}
					// problem is I'm including every possible step even on the way to shortcut
					// analysis for best cheat
					saving := value - (cheatValue + learner.v)

					if saving > 0 && saving > bestSaving {
						bestSaving = saving
					}
				}
			}

			if learner.v >= cheat {
				continue
			}

			for _, neighbour := range grid.ConnectedDirections {
				pos := grid.Position{
					X: learner.p.X + neighbour[0],
					Y: learner.p.Y + neighbour[1],
				}
				if slices.Contains(stepped, pos) {
					continue
				}

				cellQueue.Push(learn{
					p: pos,
					v: learner.v + 1,
				})
			}
		}

		fmt.Printf("best saving %d\n", bestSaving)
		if bestSaving >= countPico {
			savingsTracker[bestSaving]++
			z++
		}
		continue

		// if next position in a direction is wall check next hop for lower score
		//for _, neighbour := range grid.ConnectedDirections {
		//	c, err := g.Get(p.X+neighbour[0], p.Y+neighbour[1])
		//	if err != nil {
		//		continue
		//	}
		//	if c == '#' {
		//		// check the next tab
		//		c, err = g.Get(p.X+neighbour[0]*2, p.Y+neighbour[1]*2)
		//		if err != nil {
		//			continue
		//		}
		//		if c != '#' {
		//			cheatValue, _ := flood.Get(p.X+neighbour[0]*2, p.Y+neighbour[1]*2)
		//			if cheatValue > value {
		//				continue
		//			}
		//			//84 - (((84-0)-(84-66)) - 2)
		//			saving := value - (cheatValue + 2)
		//
		//			if saving >= 100 {
		//				z++
		//			}
		//		}
		//	}
		//}
	}

	for k, v := range savingsTracker {
		fmt.Printf("saving %d %d\n", k, v)
	}

	return z
}
