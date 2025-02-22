package main

import "testing"

func TestMap_Score(t *testing.T) {
	tests := []struct {
		name string
		m    []byte
		want int
	}{
		{
			"basic",
			[]byte(`0123
1234
8765
9876`),
			1,
		},
		{
			"2 score",
			[]byte(`...0...
...1...
...2...
6543456
7.....7
8.....8
9.....9`),
			2,
		},
		{
			"4 score",
			[]byte(`..90..9
...1.98
...2..7
6543456
765.987
876....
987....`),
			4,
		},
		{
			"36 score",
			[]byte(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`),
			36,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := LoadMap(tt.m)

			score, _ := m.Score()
			if score != tt.want {
				t.Errorf("score is not as expected, got %d, want %d", score, tt.want)
			}
		})
	}
}
