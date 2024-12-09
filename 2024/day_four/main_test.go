package main

import "testing"

func TestPuzzleSearch(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{
			"simple",
			`XMAS`,
			1,
		},
		{
			"example",
			`MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX`,
			18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			puzzle := NewPuzzle(tt.in)

			found := puzzle.Find("XMAS")
			if tt.want != found {
				t.Errorf("Puzzle.Find() = %v, want %v", tt.want, found)
			}
		})
	}
}

func TestPuzzleFindX(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{
			"example",
			`M.S
.A.
M.S`,
			1,
		},
		{
			"",
			`.M.S......
..A..MSMS.
.M.S.MAA..
..A.ASMSM.
.M.S.M....
..........
S.S.S.S.S.
.A.A.A.A..
M.M.M.M.M.
..........`,
			9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			puzzle := NewPuzzle(tt.in)

			found := puzzle.FindX("MAS")
			if tt.want != found {
				t.Errorf("Puzzle.Find() = %v, want %v", tt.want, found)
			}
		})
	}
}
