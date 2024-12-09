package main

import "testing"

func TestCheckDelta(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want bool
	}{
		{"positive", 1, 0, true},
		{"negative", -1, 0, true},
		{"zero", 0, 0, true},
		{"positive", 0, 4, false},
		{"negative", -1, -5, false},
		{"largest", 0, 3, true},
		{"across zero", -2, 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDelta(tt.a, tt.b); got != tt.want {
				t.Errorf("CheckDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDirection(t *testing.T) {
	tests := []struct {
		name string
		r    Report
		want bool
	}{
		{"positive", Report{1, 2, 3}, true},
		{"negative", Report{3, 2, 1}, true},
		{"zero", Report{1, 1, 1}, false},
		{"across zero", Report{-2, 2, 4}, true},
		{"failing", Report{1, 2, 1, 4}, false},
		{"actual line", Report{1, 3, 6, 7, 9}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDirection(tt.r); got != tt.want {
				t.Errorf("CheckDirection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReport_CalculateSafety(t *testing.T) {
	tests := []struct {
		name string
		r    Report
		want bool
	}{
		// 62 61 62 63 65 67 68 71
		{"failing", Report{62, 61, 62, 63, 65, 67, 68, 71}, true},
		// 31 34 33 32 31 29 26 24
		{"passing", Report{31, 34, 33, 32, 31, 29, 26, 24}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.CalculateSafety(true); got != tt.want {
				t.Errorf("Report.CalculateSafety() = %v, want %v", got, tt.want)
			}
		})
	}
}
