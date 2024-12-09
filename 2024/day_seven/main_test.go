package main

import "testing"

func TestCalculator(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want bool
	}{
		{"true", []byte("190: 10 19"), true},
		{"true", []byte("3267: 81 40 27"), true},
		{"false", []byte("83: 17 5"), false},
		{"false", []byte("156: 15 6"), true},
		{"false", []byte("7290: 6 8 6 15"), true},
		{"false", []byte("161011: 16 10 13"), false},
		{"false", []byte("192: 17 8 14"), true},
		{"false", []byte("21037: 9 7 18 13"), false},
		{"true", []byte("292: 11 6 16 20"), true},
		// 59196060: 7 372 3 6 28 3 9 3 4 315
		{"unsure", []byte("59196060: 7 372 3 6 28 3 9 3 4 315"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, operators, err := Extract(tt.in)

			if err != nil {
				t.Errorf("unable to extract: %s", err)
			}
			if Reason(operators, total) != tt.want {
				t.Errorf("unable to calculate reason: a=%v q=%s", Reason(operators, total), tt.in)
			}
		})
	}

}
