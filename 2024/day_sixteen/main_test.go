package main

import (
	"fmt"
	"testing"
)

var (
	exampleA = `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`

	exampleB = `#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################`
)

func TestSolveMaze(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"example", exampleA, 7036},
		{"example", exampleB, 11048},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maze, start, end := NewMaze([]byte(tt.in))
			score := maze.Solve(start, end)

			fmt.Printf("%s\n", maze)
			if score != tt.want {
				t.Errorf("Test %q failed: expected %v but got %v", tt.name, tt.want, score)
			}
		})
	}
}
