package main

import (
	"adventofcode/toolbox/tests"
	"testing"
)

func TestCalibrationScore(t *testing.T) {
	tests.TestTables(t, tests.TestTable[string, int]{
		{
			Name: "simple",
			Input: `1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet`,
			Expect: 142,
			GetResult: func(input string) (int, error) {
				return CalibrationScore([]byte(input), false), nil
			},
		},
	})
	tests.TestTables(t, tests.TestTable[string, int]{
		{
			Name: "using nonnumeric",
			Input: `two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen`,
			Expect: 281,
		},
		{
			Name:   "accept all non-numeric even overlapping",
			Input:  "hthphptmmtwo7sixsevenoneightls",
			Expect: 28,
		},
	}, func(test tests.Test[string, int], t *testing.T) {
		if CalibrationScore([]byte(test.GetInput()), true) != test.GetExpecting() {
			t.Errorf("did not get expected output, got %d, want %d", CalibrationScore([]byte(test.GetInput()), true), test.GetExpecting())
		}
	})
}
