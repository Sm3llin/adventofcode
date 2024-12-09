package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

type Rule struct {
	Left  int
	Right int
}

type Update []int

func (u Update) Valid(rules []Rule) bool {
	for i, update := range u {
		for _, rule := range rules {
			if update == rule.Right {
				valid := true
				for x := i + 1; x < len(u); x++ {
					if u[x] == rule.Left {
						valid = false
						break
					}
				}

				if !valid {
					return false
				}
			} else if update == rule.Left {
				valid := true
				for x := i - 1; x >= 0; x-- {
					if u[x] == rule.Right {
						valid = false
						break
					}
				}
				if !valid {
					return false
				}
			}
		}
	}

	return true

}

func (u Update) Fix(rules []Rule) bool {
	for i, update := range u {
		for _, rule := range rules {
			if update == rule.Right {
				for x := i + 1; x < len(u); x++ {
					if u[x] == rule.Left {
						t := u[x]
						u[x] = u[i]
						u[i] = t

						return u.Fix(rules)
					}
				}
			} else if update == rule.Left {
				for x := i - 1; x >= 0; x-- {
					if u[x] == rule.Right {
						t := u[x]
						u[x] = u[i]
						u[i] = t

						return u.Fix(rules)
					}
				}
			}
		}
	}

	return true

}

func main() {
	data, err := os.ReadFile("2024/day_five/input.txt")
	if err != nil {
		panic(err)
	}

	rules, updates := Load(data)

	fmt.Printf("Rules: %d\n", len(rules))
	fmt.Printf("Updates: %d\n", len(updates))

	valid := 0
	totalMiddlePage := 0
	totalWrongPage := 0
	for _, update := range updates {

		if update.Valid(rules) {
			//fmt.Printf("Valid: %v\n", update)
			valid++
			totalMiddlePage += update[len(update)/2]
		} else {
			update.Fix(rules)

			totalWrongPage += update[len(update)/2]
			//fmt.Printf("Invalid: %v\n", update)
		}
	}

	fmt.Printf("Valid: %d\n", valid)
	fmt.Printf("Total Middle Page: %d\n", totalMiddlePage)
	fmt.Printf("Total Wrong Page: %d\n", totalWrongPage)
}

func Load(data []byte) ([]Rule, []Update) {
	datas := bytes.Split(data, []byte("\n\n"))

	if len(datas) != 2 {
		panic("Invalid input")
	}

	bytesRules := datas[0]
	bytesUpdates := datas[1]

	var rules []Rule
	for _, rule := range bytes.Split(bytesRules, []byte("\n")) {
		bytesRule := bytes.Split(rule, []byte("|"))

		left, _ := strconv.Atoi(string(bytesRule[0]))
		right, _ := strconv.Atoi(string(bytesRule[1]))

		rules = append(rules, Rule{Left: left, Right: right})
	}

	var updates []Update
	for _, byteUpdate := range bytes.Split(bytesUpdates, []byte("\n")) {
		bytesUpdate := bytes.Split(byteUpdate, []byte(","))

		var update []int
		for _, value := range bytesUpdate {
			v, err := strconv.Atoi(string(value))

			if err != nil {
				panic("Error parsing integers from input")
			}
			update = append(update, v)
		}
		updates = append(updates, update)
	}

	return rules, updates
}
