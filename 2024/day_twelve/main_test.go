package main

import "testing"

var (
	fieldA = `AAAA
BBCD
BBCC
EEEC`

	fieldB = `RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE`
)

func TestFieldFenceCost(t *testing.T) {
	tests := []struct {
		name               string
		in                 string
		want, wantDiscount int
	}{
		{"basic", fieldA, 140, 80},
		{"complex", fieldB, 1930, 1206},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := FenceCost([]byte(tt.in))
			if value != tt.want {
				t.Errorf("Test %q failed: expected %v but got %v", tt.name, tt.want, value)
			}

			discount := FenceDiscount([]byte(tt.in))
			if discount != tt.wantDiscount {
				t.Errorf("Test %q failed: expected %v but got %v", tt.name, tt.wantDiscount, discount)
			}
		})

	}
}
