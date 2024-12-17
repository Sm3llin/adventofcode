package main

import (
	"adventofcode/toolbox/datatypes"
	"adventofcode/toolbox/tests"
	"adventofcode/toolbox/text"
	"testing"
)

func TestLoadGame(t *testing.T) {
	tests.TestTables(
		t,
		tests.TestTable[string, *datatypes.Inventory[text.Text]]{
			{
				Name:  "game 1",
				Input: "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
				Expect: &datatypes.Inventory[text.Text]{
					Label: "Game 1",
					Stock: map[string]int{
						"blue":  14,
						"red":   12,
						"green": 13,
					},
				},
				GetResult: func(input string) (*datatypes.Inventory[text.Text], error) {
					i, _, err := NewGame([]byte(input))
					return i, err
				},
			},
		},
	)
}
