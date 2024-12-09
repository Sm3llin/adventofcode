package main

import (
	"os"
	"testing"
)

func TestMap_Solve(t *testing.T) {
	tests := []struct {
		Filename  string
		Want      int
		Blockages int
	}{
		{"test_input.txt", 41, 6},
		{"test_input.2.txt", 7, 0},
		{"test_input.3.txt", 6, 1},
		{"test_input.4.txt", 10, 3},
		{"test_input.5.txt", 7, 1},
		{"test_input.6.txt", 7, 2},
		{"test_input.7.txt", 7, 2},
		{"test_input.8.txt", 7, 1},
		{"test_input.9.txt", 6, 1},
	}

	for _, test := range tests {
		t.Run(test.Filename, func(t *testing.T) {
			data, err := os.ReadFile(test.Filename)
			if err != nil {
				t.Error(err)
			}

			m := NewMap(data)

			want, block := m.Solve()
			if want != test.Want {
				t.Errorf("Expected %d unique steps, got %d", test.Want, want)
			}
			if block != test.Blockages {
				t.Errorf("Expected %d blockortunities, got %d", test.Blockages, block)
			}

		})
	}
}
