package main

import "testing"

func TestUpdate_Valid(t *testing.T) {
	tests := []struct {
		Name  string
		In    Update
		Rules []Rule
		Want  bool
	}{
		{
			"simple",
			Update{1, 2, 4, 8},
			[]Rule{
				{2, 8},
			},
			true,
		},
		{
			"failure",
			Update{1, 2, 4, 8},
			[]Rule{
				{4, 2},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := tt.In.Valid(tt.Rules); got != tt.Want {
				t.Errorf("Update.Valid() = %v, want %v", got, tt.Want)
			}
		})
	}
}
