package main

import (
	"os"
	"regexp"
	"sort"
	"strconv"
)

var reNumbers = regexp.MustCompile(`(\d+) +(\d+)`)

func main() {
	data, err := os.ReadFile("2024/day_one/input.txt")
	if err != nil {
		panic(err)
	}

	var locationIDsA = []int{}
	var locationIDsB = []int{}

	matches := reNumbers.FindAllSubmatch(data, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			// Convert the byte slices to strings and then parse to integers
			idA, errA := strconv.Atoi(string(match[1]))
			idB, errB := strconv.Atoi(string(match[2]))
			if errA != nil || errB != nil {
				panic("Error parsing integers from input")
			}
			locationIDsA = append(locationIDsA, idA)
			locationIDsB = append(locationIDsB, idB)
		}
	}

	// sort the lists
	sort.Ints(locationIDsA)
	sort.Ints(locationIDsB)

	EstimateDistanceDelta(locationIDsA, locationIDsB)
	CalculateSimilarityScore(locationIDsA, locationIDsB)

	return
}

func CalculateSimilarityScore(listA, listB []int) int {
	// create a count of list B
	countB := map[int]int{}

	for _, id := range listB {
		countB[id]++
	}

	totalSimilarity := 0
	for _, id := range listA {
		totalSimilarity += id * countB[id]
	}

	println(totalSimilarity)
	return totalSimilarity
}

func EstimateDistanceDelta(listA, listB []int) int {
	totalDistance := 0
	for i := 0; i < len(listA); i++ {
		delta := listA[i] - listB[i]

		// set abs
		if delta < 0 {
			delta = -delta
		}

		totalDistance += delta
	}

	println(totalDistance)
	return totalDistance
}
