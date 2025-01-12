package main

import (
	"adventofcode/toolbox/tests"
	"testing"
)

var (
	example1 = `###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`

	example2 = `######
#S##E#
#....#
######`
	example3 = `######
#S##E#
#.#..#
#...##
######`
)

// T7: 4, 2

type input struct {
	s         string
	cheatTime int
	inclusive int
}

func TestDay20(t *testing.T) {
	tests.TestTables(t, tests.TestTable[input, int]{
		{
			Name:   "exampleA",
			Input:  input{example1, 2, 0},
			Expect: 44,
		},
		{
			Name:   "exampleB",
			Input:  input{example2, 2, 0},
			Expect: 0,
		},
		{
			Name:   "exampleB",
			Input:  input{example2, 3, 0},
			Expect: 1,
		},
		{
			Name:   "exampleC",
			Input:  input{example3, 3, 0},
			Expect: 4,
		},
		{
			Name:   "exampleA",
			Input:  input{example1, 20, 50},
			Expect: 285,
		},
	}, func(test tests.Test[input, int], t *testing.T) {
		n := solve([]byte(test.GetInput().s), test.GetInput().cheatTime, test.GetInput().inclusive)

		if n != test.GetExpecting() {
			t.Errorf("did not get expected output, got %d, want %d", n, test.GetExpecting())
		}
	})
}
