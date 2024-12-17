package main

import (
	"adventofcode"
	"adventofcode/toolbox/assert"
	"adventofcode/toolbox/conversion"
	"adventofcode/toolbox/datatypes"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/text"
	"bytes"
	"fmt"
)

func NewGame(data []byte) (*datatypes.Inventory[text.Text], []text.Texts, error) {
	d := text.Text(data).Split(":", 1)

	assert.LengthMin(d, 2)
	game, setString := d[0], d[1]
	sets := setString.Split(";", -1)

	var steps []text.Texts
	inventory := datatypes.NewInventory(game)

	for _, set := range sets {
		for _, s := range set.Split(",", -1) {
			step := s.TrimSpace().Split(" ", 1)

			steps = append(steps, step)
		}
	}

	return inventory, steps, nil
}

func main() {
	adventofcode.Time(func() {
		d := fs.LoadFile("2023/day_two/input.txt")

		var total int

		for i, line := range bytes.Split(d, []byte{'\n'}) {
			inventory, steps, _ := NewGame(line)

			var err error
			for _, step := range steps {
				// set static inventory
				inventory.SetX("blue", 14)
				inventory.SetX("red", 12)
				inventory.SetX("green", 13)

				n, errInt := conversion.ToInt(step[0])
				assert.NoError(errInt)
				err = inventory.RemoveX(step[1].String(), n)

				if err != nil {
					break
				}
			}
			if err == nil {
				total += i + 1
			}
		}

		fmt.Printf("Part 1: %d\n", total)
	})
	adventofcode.Time(func() {
		d := fs.LoadFile("2023/day_two/input.txt")

		var total int

		for _, line := range bytes.Split(d, []byte{'\n'}) {
			_, steps, _ := NewGame(line)

			var maxBlue, maxRed, maxGreen int

			var err error
			for _, step := range steps {
				n, errInt := conversion.ToInt(step[0])
				assert.NoError(errInt)

				switch step[1] {
				case "blue":
					if n > maxBlue {
						maxBlue = n
					}
				case "red":
					if n > maxRed {
						maxRed = n
					}
				case "green":
					if n > maxGreen {
						maxGreen = n
					}
				}

				if err != nil {
					break
				}
			}
			if err == nil {
				total += maxBlue * maxRed * maxGreen
			}
		}

		fmt.Printf("Part 2: %d\n", total)
	})
}
