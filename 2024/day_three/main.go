package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var reMul = regexp.MustCompile(`mul\((\d+),(\d+)\)|(do(?:n't)?)\(\)`)

type Mul struct {
	A int
	B int
}

func (m Mul) Calculate() int {
	return m.A * m.B
}

func ProcessLocator(in string) []Mul {
	matches := reMul.FindAllStringSubmatch(in, -1)

	var ignore bool
	elements := []Mul{}
	for _, match := range matches {
		if len(match) >= 3 {
			if match[3] == "do" {
				ignore = false
				continue
			} else if match[3] == "don't" {
				ignore = true
				continue
			} else if ignore {
				continue
			}
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])

			elements = append(elements, Mul{A: a, B: b})
		}
	}

	return elements
}

func main() {
	file, err := os.ReadFile("2024/day_three/input.txt")
	if err != nil {
		panic(err)
	}

	elements := ProcessLocator(string(file))

	var total int
	for _, element := range elements {
		total += element.Calculate()
	}

	fmt.Printf("Total: %d\n", total)
}
