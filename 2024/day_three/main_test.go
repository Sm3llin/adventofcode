package main

import "testing"

func TestProcessLocator(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []Mul
	}{
		{"basic", "mul(1,2)", []Mul{{1, 2}}},
		{"double", "mul(2,5)mul(6,2)", []Mul{{2, 5}, {6, 2}}},
		{"ignore", "mul(4*, mul(6,9!, ?(12,34), or mul ( 2 , 4 )", []Mul{}},
		{"complex", "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))", []Mul{{2, 4}, {5, 5}, {11, 8}, {8, 5}}},
		{"logic", "mul(2,5)don't()mul(6,2)", []Mul{{2, 5}}},
		{"logic_extra", "mul(2,5)don't()mul(6,2)do()mul(4,2)", []Mul{{2, 5}, {4, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessLocator(tt.in)
			if len(got) != len(tt.want) {
				t.Fatalf("Test %q failed: expected %v but got %v", tt.name, tt.want, got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Test %q failed at index %d: expected %v but got %v", tt.name, i, tt.want[i], got[i])
				}
			}
		})
	}
}
