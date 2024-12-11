package main

import (
	"fmt"
	"testing"
)

func TestDayEleven(t *testing.T) {
	tests := []struct {
		input  string
		blinks int
		want   string
	}{
		{"125 17", 0, "125 17"},
		{"125 17", 1, "253000 1 7"},
		{"125 17", 4, "512 72 2024 2 0 2 4 2867 6032"},
		{"125 17", 6, "2097446912 14168 4048 2 0 2 4 40 48 2024 40 48 80 96 2 8 6 7 6 0 3 2"},
	}
	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			node := LoadNodeList([]byte(test.input))

			for range test.blinks {
				node.Blink()
			}
			if node.Render() != test.want {
				t.Errorf("not as expected, got=%v, want=%v", node.Render(), test.want)
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
