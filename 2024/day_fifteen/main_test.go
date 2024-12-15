package main

import (
	"testing"
)

var (
	inputA = `##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########

<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^
vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v
><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<
<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^
^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><
^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^
>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^
<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>
^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>
v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^`
	inputB = `##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########

^<<`
)

func TestNewFactory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		steps    int
	}{
		{
			name:  "basic",
			input: inputA,
			expected: `##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########`,
		},
		{
			name:  "movement",
			input: inputB,
			steps: 1,
			expected: `##########
#..O..O.O#
#......O.#
#.OO@.O.O#
#..O...O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########`,
		},
		{
			name:  "push",
			input: inputB,
			steps: 2,
			expected: `##########
#..O..O.O#
#......O.#
#OO@..O.O#
#..O...O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########`,
		},
		{
			name:  "cannot push",
			input: inputB,
			steps: 3,
			expected: `##########
#..O..O.O#
#......O.#
#OO@..O.O#
#..O...O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory([]byte(tt.input))

			for range tt.steps {
				factory.GetRobot().Step(factory)
			}

			if factory.Render() != tt.expected+"\n" {
				t.Errorf("Test %q failed: expected %v but got %v", tt.name, tt.expected, factory.Render())
			}
		})
	}
}
