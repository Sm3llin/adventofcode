package main

import (
	"fmt"
	"testing"
)

func TestDayEleven(t *testing.T) {
	tests := []struct {
		input  string
		blinks int
		want   uint
		actual string
	}{
		{"125 17", 0, 2, "125 17"},
		{"1200 1", 1, 3, "125 17"},
		{"125 17", 1, 3, "253000 1 7"},
		{"125 17", 4, 9, "512 72 2024 2 0 2 4 2867 6032"},
		{"125 17", 6, 22, "2097446912 14168 4048 2 0 2 4 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2"},
		{"125 17", 25, 55312, ""},
		// 50 == 24sec, 49 == 16, 48 == 10
		{"1 1 1 1 1 1", 150, uint(2363334184417015096), ""},
	}
	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			node := LoadNodeList([]byte(test.input))

			count := node.Count(test.blinks)
			if count != uint(test.want) {
				t.Errorf("not as expected, got=%v, want=%v", count, test.want)
			}
		})
	}
}

func TestSplitEvenDigitString(t *testing.T) {
	tests := []struct {
		in  int
		out []int
	}{
		{2024, []int{20, 24}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.in), func(t *testing.T) {
			a, b := SplitEvenDigitString(test.in)

			if a != test.out[0] || b != test.out[1] {
				t.Errorf("not as expected, got=%v, want=%v", a, test.out[0])
			}

		})
	}
}
