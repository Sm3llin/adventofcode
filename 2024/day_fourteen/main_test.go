package main

import "testing"

func TestRobotSafety(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  int
	}{
		{
			"example",
			[]byte(`p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3
`),
			12,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			robots := LoadRobots(test.input)
			grid := Grid{Height: 7, Width: 11}

			WaitFor(robots, 100, grid)

			value := SafetyScore(grid, robots, true)
			if value != test.want {
				t.Errorf("expected: %d, got: %d", test.want, value)
			}
		})
	}
}

func TestRobotPosition(t *testing.T) {
	tests := []struct {
		name          string
		grid          Grid
		robot         Robot
		seconds       int
		expectedRobot Robot
	}{
		{
			"position after 1",
			Grid{Height: 4, Width: 4},
			Robot{
				X:  0,
				Y:  4,
				vX: 3,
				vY: -3,
			},
			1,
			Robot{
				X:  3,
				Y:  1,
				vX: 3,
				vY: -3,
			},
		},
		{
			"position moving after 1",
			Grid{Height: 4, Width: 4},
			Robot{
				X:  0,
				Y:  4,
				vX: -3,
				vY: -3,
			},
			1,
			Robot{
				X:  1,
				Y:  1,
				vX: -3,
				vY: -3,
			},
		},
		{
			"move further than 1 grid",
			Grid{Height: 4, Width: 4},
			Robot{
				X:  0,
				Y:  4,
				vX: -3,
				vY: -7,
			},
			1,
			Robot{
				X:  1,
				Y:  1,
				vX: -3,
				vY: -7,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.robot.Move(test.grid, test.seconds)

			if test.expectedRobot != test.robot {
				t.Errorf("expected: %v, got: %v", test.expectedRobot, test.robot)
			}
		})
	}
}
