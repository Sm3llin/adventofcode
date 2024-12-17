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
