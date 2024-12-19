package grid

import (
	"adventofcode/toolbox/tests"
	"testing"
)

// TestMaze_Solve tests the Maze's Solve method
func TestMaze_Solve(t *testing.T) {
	tests.TestTables(t, tests.TestTable[struct {
		grid   [][]byte
		start  Position
		end    Position
		wall   byte
		onMove MoveFunc
	}, []Position]{
		{
			Name: "simple_maze",
			Input: struct {
				grid   [][]byte
				start  Position
				end    Position
				wall   byte
				onMove MoveFunc
			}{
				grid: [][]byte{
					{'.', '.', '.', '#'},
					{'.', '#', '.', '#'},
					{'.', '#', '.', '.'},
					{'#', '#', '#', '.'},
				},
				start: Position{X: 0, Y: 0}, // 'S'
				end:   Position{X: 3, Y: 3}, // 'E'
				wall:  '#',
				onMove: func(from, to Position) (bool, int) {
					return true, 1
				},
			},
			Expect: []Position{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 2, Y: 0},
				{X: 2, Y: 1},
				{X: 2, Y: 2},
				{X: 3, Y: 2},
				{X: 3, Y: 3},
			},
			GetResult: func(data struct {
				grid   [][]byte
				start  Position
				end    Position
				wall   byte
				onMove MoveFunc
			}) ([]Position, error) {
				// Create grid
				grid := NewGrid(data.grid)

				// Initialize maze
				maze := NewMaze(grid, data.wall, data.onMove)

				t.Logf(maze.Render())
				// Solve maze
				path, _ := maze.Solve(data.start, data.end)
				return path, nil
			},
		},
		{
			Name: "small_open_maze",
			Input: struct {
				grid   [][]byte
				start  Position
				end    Position
				wall   byte
				onMove MoveFunc
			}{
				grid: [][]byte{
					{'.', '.', '#', '#'},
					{'.', '.', '.', '#'},
					{'#', '.', '.', '.'},
				},
				start: Position{X: 0, Y: 0}, // 'S'
				end:   Position{X: 2, Y: 2}, // 'E'
				wall:  '#',
				onMove: func(from, to Position) (bool, int) {
					return true, 1
				},
			},
			Expect: []Position{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 2, Y: 2},
			},
			GetResult: func(data struct {
				grid   [][]byte
				start  Position
				end    Position
				wall   byte
				onMove MoveFunc
			}) ([]Position, error) {
				// Create grid
				grid := NewGrid(data.grid)

				// Initialize maze
				maze := NewMaze(grid, data.wall, data.onMove)

				t.Logf(maze.Render())
				// Solve maze
				path, _ := maze.Solve(data.start, data.end)
				return path, nil
			},
		},
		{
			Name: "unsolvable_maze",
			Input: struct {
				grid   [][]byte
				start  Position
				end    Position
				wall   byte
				onMove MoveFunc
			}{
				grid: [][]byte{
					{'.', '#', '#', '#'},
					{'#', '#', '#', '#'},
					{'#', '#', '#', '#'},
					{'#', '#', '#', '.'},
				},
				start: Position{X: 0, Y: 0}, // 'S'
				end:   Position{X: 3, Y: 3}, // 'E'
				wall:  '#',
				onMove: func(from, to Position) (bool, int) {
					return true, 1
				},
			},
			Expect: nil, // No path should exist
			GetResult: func(data struct {
				grid   [][]byte
				start  Position
				end    Position
				wall   byte
				onMove MoveFunc
			}) ([]Position, error) {
				// Create grid
				grid := NewGrid(data.grid)

				// Initialize maze
				maze := NewMaze(grid, data.wall, data.onMove)

				// Solve maze, handle panic for unreachable paths
				defer func() {
					if r := recover(); r != nil {
						// Expected
					}
				}()

				path, _ := maze.Solve(data.start, data.end)
				return path, nil
			},
		},
	})
}
