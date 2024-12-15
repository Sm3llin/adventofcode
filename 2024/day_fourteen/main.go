package main

import (
	"adventofcode"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"sync"
)

const (
	Width  = 101
	Height = 103
)

type Grid struct {
	Width, Height int
}

func (g Grid) CheckBounds(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

// Create a list of Robot positions
type Robot struct {
	X, Y int

	vX, vY int
}

func (robot *Robot) Move(grid Grid, seconds int) {
	for range seconds {
		robot.X += robot.vX
		robot.Y += robot.vY

		// if x less than 0 add length of grid
		// if x greater than grid remove length of grid
		// repeat for Y
		for !grid.CheckBounds(robot.X, robot.Y) {
			if robot.X < 0 {
				robot.X += grid.Width
			} else if robot.X >= grid.Width {
				robot.X -= grid.Width
			}

			if robot.Y < 0 {
				robot.Y += grid.Height
			} else if robot.Y >= grid.Height {
				robot.Y -= grid.Height
			}
		}
	}
}

func WaitFor(robots []*Robot, seconds int, grid Grid) {
	for _, robot := range robots {
		robot.Move(grid, seconds)
	}
}

func LoadRobots(data []byte) []*Robot {
	robots := []*Robot{}

	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		robot, err := ParseRobotData(string(line))

		if err != nil {
			panic(err)
		}

		robots = append(robots, &robot)
	}

	return robots
}

func IsEgg(grid Grid, robots []*Robot) bool {
	robotTileCount := make([][]int, grid.Height)
	for i := range robotTileCount {
		robotTileCount[i] = make([]int, grid.Width)
	}
	for _, robot := range robots {
		robotTileCount[robot.Y][robot.X]++
	}

	// first guess will be looking at whitespace counts and seeing if we see a major
	// grouping at same time
	spaceGuess := make(map[int]int)
	for y, row := range robotTileCount {
		var space bool
		for _, tile := range row {
			if !space && tile == 0 {
				space = true
				spaceGuess[y]++
			} else if tile != 0 {
				space = false
			}
		}
	}

	threshold := 0
	isEgg := true
	for _, count := range spaceGuess {
		// Apparently I don't read and the egg is a xmas tree, but these settings found it
		isValid := count < 6
		if !isValid && threshold > 10 {
			isEgg = false
			break
		} else if !isValid {
			threshold++
		}
	}
	return isEgg
}

func EggHunt(grid Grid, robots []*Robot) int {
	// the robots will form the shape of an egg at a number of seconds. It
	// is probably not at a simple count and I shouldn't try visually complete
	// this but it does sounds fun
	var x int
	wg := sync.WaitGroup{}
	for x = range 10000 {
		fmt.Printf("%d  \r", x)
		if IsEgg(grid, robots) {
			fmt.Printf("%d seconds\n", x)
			Print(grid, robots)
			fmt.Printf("\n")
		}

		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//	time.Sleep(200 * time.Millisecond)
		//}()
		for _, robot := range robots {
			robot.Move(grid, 1)
		}
		wg.Wait()
	}
	return x
}

func Print(grid Grid, robots []*Robot) {
	robotTileCount := make([][]int, grid.Height)
	for i := range robotTileCount {
		robotTileCount[i] = make([]int, grid.Width)
	}
	for _, robot := range robots {
		robotTileCount[robot.Y][robot.X]++
	}

	for _, row := range robotTileCount {
		for _, tile := range row {
			if tile == 0 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%d", tile)
			}
		}
		fmt.Printf("\n")
	}
	for range grid.Width {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
}

func SafetyScore(grid Grid, robots []*Robot, debug ...bool) int {
	var debugMode bool
	if len(debug) > 0 {
		debugMode = debug[0]
	}
	robotTileCount := make([][]int, grid.Height)
	for i := range robotTileCount {
		robotTileCount[i] = make([]int, grid.Width)
	}
	for _, robot := range robots {
		robotTileCount[robot.Y][robot.X]++
	}

	scores := make([][]int, 2)

	for i := range scores {
		scores[i] = make([]int, 2)
	}

	xMiddle := grid.Width / 2
	yMiddle := grid.Height / 2

	// separate into quadrants
	for y, row := range robotTileCount {
		yPosition := 0

		if y == yMiddle {
			if debugMode {
				for _, tile := range row {
					fmt.Printf("%d", tile)
				}
				fmt.Printf("\n")
			}
			continue
		} else if y < yMiddle {
			yPosition = 0
		} else {
			yPosition = 1
		}

		for x, tile := range row {
			if debugMode {
				fmt.Printf("%d", tile)
			}

			var xPosition int
			if x == xMiddle {
				continue
			} else if x < xMiddle {
				xPosition = 0
			} else {
				xPosition = 1
			}

			scores[yPosition][xPosition] += tile
		}
		if debugMode {

			fmt.Printf("\n")
		}
	}

	return scores[0][0] * scores[1][1] * scores[0][1] * scores[1][0]
}

func ParseRobotData(line string) (Robot, error) {
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 5 {
		return Robot{}, fmt.Errorf("invalid line format")
	}

	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])
	vx, _ := strconv.Atoi(matches[3])
	vy, _ := strconv.Atoi(matches[4])

	return Robot{X: x, Y: y, vX: vx, vY: vy}, nil
}

func main() {
	adventofcode.Time(func() {
		data := adventofcode.LoadFile("2024/day_fourteen/input.txt")
		robots := LoadRobots(data)
		grid := Grid{Width: Width, Height: Height}

		WaitFor(robots, 100, grid)
		fmt.Printf("Safety Score: %d\n", SafetyScore(grid, robots))
	})
	adventofcode.Time(func() {
		data := adventofcode.LoadFile("2024/day_fourteen/input.txt")
		robots := LoadRobots(data)
		grid := Grid{Width: Width, Height: Height}

		fmt.Printf("Egg Hunt: %d\n", EggHunt(grid, robots))
	})
}
