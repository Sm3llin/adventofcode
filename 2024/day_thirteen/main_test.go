package main

import (
	"testing"
)

func TestCountCheapestPrizes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name: "Single machine, simple input",
			input: `Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400`,
			expected: 280,
		},
		{
			name: "Another one",
			input: `Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450`,
			expected: 200,
		},
		{
			name: "cannot win",
			input: `Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load machines from input
			machines := LoadMachines([]byte(tt.input))

			// Call the function under test
			result := CountCheapestPrizes(machines)

			// Validate the result
			if result != tt.expected {
				t.Errorf("Test %s failed: expected %d, got %d", tt.name, tt.expected, result)
			}
		})
	}
}
