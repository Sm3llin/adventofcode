package main

import (
	"adventofcode/toolbox/grid"
	"adventofcode/toolbox/tests"
	"fmt"
	"testing"
)

func TestSimulatedMaze(t *testing.T) {
	tests.TestTables(t, tests.TestTable[string, grid.Grid[string]]{
		{
			Name: "example",
			Input: `5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`,
		},
	}, func(test tests.Test[string, grid.Grid[string]], t *testing.T) {
		positions := convertInput([]byte(test.GetInput()))

		g := grid.NewGridValue(".", 7, 7)

		for i := range 12 {
			p := positions[i]
			g.Set(p.X, p.Y, "#")
		}
		fmt.Println(g.Render())

		maze := grid.NewMaze(g, "#", func(from, to grid.Position) (allow bool, score int) {
			return true, 1
		})

		path, success := maze.Solve(grid.Position{}, grid.Position{X: 6, Y: 6})

		if !success {
			t.Errorf("did not get expected output, got %v, want %v", path, nil)
		}

		if len(path)-1 != 22 {
			t.Errorf("did not get expected output, got %v, want %v", len(path), 22)
		}

		fmt.Printf("%v\n", path)
		fmt.Printf("%d\n", len(path))
	})
}
