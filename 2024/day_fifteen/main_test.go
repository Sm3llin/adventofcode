package main

import (
	"fmt"
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
	inputC = `#######
#...#.#
#.....#
#..OO@#
#..O..#
#.....#
#######

<vv<<^^<<^^`
)

func TestNewFactory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		steps    int
		isBig    bool
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
		{
			name:  "is big",
			input: inputB,
			steps: 3,
			isBig: true,
			expected: `####################
##....[]....[]..[]##
##............[]..##
##[][]@.....[]..[]##
##....[]......[]..##
##[]##....[]......##
##[]....[]....[]..##
##..[][]..[]..[][]##
##........[]......##
####################`,
		},
		{
			name:  "is big",
			input: inputC,
			steps: -1,
			isBig: true,
			expected: `##############
##...[].##..##
##...@.[]...##
##....[]....##
##..........##
##..........##
##############`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory([]byte(tt.input), tt.isBig)

			robot := factory.GetRobot()

			if tt.steps == -1 {
				for robot.Step(factory) {
					// running
					fmt.Printf(factory.Render())
				}
			}
			for x := range tt.steps {
				fmt.Printf("step %d\n", x)
				robot.Step(factory)
			}

			if factory.Render() != tt.expected+"\n" {
				t.Errorf("Test %q failed: expected %v but got %v", tt.name, tt.expected, factory.Render())
			}
		})
	}
}

func TestRobot(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		isBig    bool
		expected int
	}{
		{
			name:     "basic",
			input:    inputA,
			isBig:    false,
			expected: 10092,
		},
		{
			name:     "big",
			input:    inputA,
			isBig:    true,
			expected: 9021,
		},
		{
			name:     "smaller example",
			input:    inputC,
			isBig:    true,
			expected: 618,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory([]byte(tt.input), tt.isBig)

			RunAutomation(factory)

			if factory.Score() != tt.expected {
				t.Errorf("Test %q failed: expected %d but got %d", tt.name, tt.expected, factory.Score())
			}
			fmt.Printf("%s", factory.Render())
		})
	}
}
