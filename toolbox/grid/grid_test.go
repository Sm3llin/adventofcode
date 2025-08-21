package grid

import (
	"adventofcode/toolbox/tests"
	"testing"
)

func TestNewGrid(t *testing.T) {
	tests.TestTables(t, tests.TestTable[[][]byte, Grid[byte]]{
		{
			Name:  "ensure grid creates",
			Input: [][]byte{{'a', 'b', 'c'}, {'d', 'e', 'f'}, {'g', 'h', 'i'}},
			Expect: Grid[byte]{
				Data: [][]byte{
					{'a', 'b', 'c'},
					{'d', 'e', 'f'},
					{'g', 'h', 'i'},
				},
				Width:  3,
				Height: 3,
			},
			GetResult: func(data [][]byte) (Grid[byte], error) {
				return NewGrid(data), nil
			},
		},
	})
}

func TestGrid_All(t *testing.T) {
	tests.TestTables(t, tests.TestTable[[][]byte, Grid[byte]]{
		{
			Name:  "ensure grid all is range",
			Input: [][]byte{{'a', 'b', 'c'}, {'d', 'e', 'f'}, {'g', 'h', 'i'}},
			GetResult: func(data [][]byte) (Grid[byte], error) {
				return NewGrid(data), nil
			},
		},
	}, func(test tests.Test[[][]byte, Grid[byte]], t *testing.T) {
		grid, err := test.GetResult(test.GetInput())

		if err != nil {
			t.Errorf("did not get expected error, got %v, want %v", err, nil)
		}

		var i Position
		var v byte
		for i, v = range grid.All() {
		}
		if i.X != 2 || i.Y != 2 {
			t.Errorf("did not get expected output, got %d, want %d", i, 9)
		}
		if v != 'i' {
			t.Errorf("did not get expected output, got %d, want %d", v, 'i')
		}
	})
}
