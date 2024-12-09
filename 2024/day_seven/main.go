package main

import (
	"adventofcode"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	adventofcode.Time(func() {

		data, err := os.ReadFile("2024/day_seven/input.txt")
		if err != nil {
			panic(err)
		}

		var sum int
		for _, line := range bytes.Split(data, []byte("\n")) {
			total, operators, err := Extract(line)
			if err != nil {
				panic(err)
			}

			if Reason(operators, total) {
				//fmt.Printf("True: %s\n", line)
				sum += total
			} else {
				//fmt.Printf("False: %s\n", line)
			}
		}
		fmt.Printf("Sum: %d\n", sum)
	})
}

type recursingCalc func(current int, remaining []int, expected int, next recursingCalc) bool

func Reason(operators []int, expectedResult int) bool {
	// unable to reason the operator type we need to try all the different combinations to get the result
	add := func(a, b int) int {
		return a + b
	}
	mul := func(a, b int) int {
		return a * b
	}
	concat := func(a, b int) int {
		v, err := strconv.Atoi(fmt.Sprintf("%d%d", a, b))
		if err != nil {
			panic("Error parsing integers from input")
		}
		return v
	}

	calc := func(current int, remaining []int, expected int, next recursingCalc) bool {
		if current > expected {
			//result <- false
			return false
		} else if len(remaining) == 0 {
			//result <- current == expected
			return current == expected
		}

		b := remaining[0]
		remaining = remaining[1:]

		return next(add(current, b), remaining, expected, next) || next(mul(current, b), remaining, expected, next) || next(concat(current, b), remaining, expected, next)
	}

	return calc(operators[0], operators[1:], expectedResult, calc)
}

func Extract(data []byte) (int, []int, error) {
	pieces := bytes.Split(data, []byte(": "))
	if len(pieces) != 2 {
		return 0, nil, fmt.Errorf("invalid input")
	}

	total, err := strconv.Atoi(string(pieces[0]))
	if err != nil {
		return 0, nil, fmt.Errorf("invalid input")
	}

	operators := []int{}
	for _, value := range bytes.Split(pieces[1], []byte{' '}) {
		if len(value) == 0 {
			continue
		}

		operator, err := strconv.Atoi(string(value))
		if err != nil {
			return 0, nil, fmt.Errorf("invalid input")
		}

		operators = append(operators, operator)
	}

	return total, operators, nil
}
