package main

import (
	"adventofcode/toolbox/tests"
	"testing"
)

func TestNewComputer(t *testing.T) {
	tests.TestTables(t,
		tests.TestTable[string, string]{
			{
				Name: "example 1",
				Input: `Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`,
				Expect: "4,6,3,5,6,3,5,2,1,0",
			},
			{
				Name: "example 2",
				Input: `Register A: 0
Register B: 0
Register C: 9

Program: 2,6,5,5`,
				Expect: "1",
			},
			{
				Name: "example 3",
				Input: `Register A: 2024
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`,
				Expect: "4,2,5,6,7,7,7,7,3,1,0",
			},
			{
				Name: "example 4",
				Input: `Register A: 0
Register B: 2024
Register C: 43690

Program: 4,0,5,5`,
				Expect: "2",
			},
			{
				Name: "example 5",
				Input: `Register A: 0
Register B: 29
Register C: 0

Program: 1,7,5,5`,
				Expect: "2",
			},
		}, func(test tests.Test[string, string], t *testing.T) {
			c := NewComputer([]byte(test.GetInput()))

			c.Run()

			// Read out of the Stdout
			outputData := c.StdoutString()
			if outputData != (test.GetExpecting()) {
				t.Errorf("Expected %s but got %s", test.GetExpecting(), outputData)
			}
		})
}

func TestSolveFor2(t *testing.T) {
	var find int
	a := 0
	for find < 8 {

		b1 := a % 8
		b1 = a ^ b1
		c1 := a/2 ^ b1
		b1 = b1 ^ (c1 % 8)
		b1 = b1 ^ 6

		if b1%8 == 2 {
			t.Logf("found %d", a)
			find++
			a++
		} else {
			a++
		}
	}
}
