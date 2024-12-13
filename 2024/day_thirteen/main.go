package main

import (
	"adventofcode"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Machine struct {
	// A costs 3, B costs 1
	ButtonA, ButtonB Button

	Prize Prize
}

func (m Machine) Solve() []Solution {
	var solution Solution
	var token, tokenCost uint
	tokenCost--

	xPrize, yPrize := m.Prize.X, m.Prize.Y
	xAButton, xBButton := m.ButtonA.X, m.ButtonB.X
	yAButton, yBButton := m.ButtonA.Y, m.ButtonB.Y
	portion := xPrize
	remainder := xPrize % xAButton

	if remainder != 0 {
		portion -= remainder
	}

	for portion > 0 {

		// check if other button can be modulo
		if (xPrize-portion)%xBButton == 0 {
			aPresses := portion / xAButton
			bPresses := (xPrize - portion) / xBButton

			exceededA := aPresses > 100
			exceededB := bPresses > 100

			// no more than 100 presses allowed
			if exceededA && exceededB {
				break
			} else if exceededA || exceededB {
				portion -= xAButton
				continue
			}

			// this could be valid now so we should check Y buttons are valid
			if aPresses*yAButton+bPresses*yBButton == yPrize {
				token = uint(aPresses*3 + bPresses)

				if token < tokenCost {
					solution.A = aPresses
					solution.B = bPresses
				}
			}
		}

		portion -= xAButton
	}

	return []Solution{solution}
}

type (
	Button struct {
		X, Y int
	}
	Prize struct {
		X, Y int
	}

	Solution struct {
		A, B int
	}
)

var (
	reAButtonPositions = regexp.MustCompile(`Button A: X\+?(-?\d+), Y\+?(-?\d+)`)
	reBButtonPositions = regexp.MustCompile(`Button B: X\+?(-?\d+), Y\+?(-?\d+)`)
	rePrizePositions   = regexp.MustCompile(`Prize: X=(-?\d+), Y=(-?\d+)`)
)

func (s Solution) Tokens() int {
	return s.A*3 + s.B
}

// Placeholder for the CountCheapestPrizes function
func CountCheapestPrizes(machines []Machine) int {
	var tokens int

	for i, machine := range machines {
		_ = i
		solutions := machine.Solve()
		//fmt.Printf("Solved %d/%d\n", i+1, len(machines))

		var cheapestSolution *Solution
		for _, solution := range solutions {
			if cheapestSolution == nil {
				cheapestSolution = &solution
				continue
			}
			if solution.Tokens() < cheapestSolution.Tokens() {
				cheapestSolution = &solution
			}
		}
		if cheapestSolution != nil {
			tokens += cheapestSolution.Tokens()
		}
	}

	return tokens
}

func LoadMachines(data []byte) []Machine {
	machines := []Machine{}
	for _, line := range strings.Split(string(data), "\n\n") {
		if line == "" {
			continue
		}
		//Button A: X+94, Y+34
		//Button B: X+22, Y+67
		//Prize: X=8400, Y=5400
		aButton := reAButtonPositions.FindStringSubmatch(line)
		bButton := reBButtonPositions.FindStringSubmatch(line)
		prize := rePrizePositions.FindStringSubmatch(line)

		aButtonX, _ := strconv.Atoi(aButton[1])
		aButtonY, _ := strconv.Atoi(aButton[2])

		bButtonX, _ := strconv.Atoi(bButton[1])
		bButtonY, _ := strconv.Atoi(bButton[2])

		prizeX, _ := strconv.Atoi(prize[1])
		prizeY, _ := strconv.Atoi(prize[2])

		machine := Machine{
			ButtonA: Button{X: aButtonX, Y: aButtonY},
			ButtonB: Button{X: bButtonX, Y: bButtonY},
			Prize:   Prize{X: prizeX, Y: prizeY},
		}

		machines = append(machines, machine)
	}
	return machines
}

func main() {
	adventofcode.Time(func() {
		data := adventofcode.LoadFile("2024/day_thirteen/input.txt")
		machines := LoadMachines(data)

		// part 1: 404us
		fmt.Printf("Part 1: Tokens=%d\n", CountCheapestPrizes(machines))
		//
	})
}
