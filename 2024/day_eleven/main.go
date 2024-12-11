package main

import (
	"adventofcode"
	"bytes"
	"fmt"
	"strconv"
	"sync"
)

type Direction []int

var (
	UP    = Direction{0, -1}
	DOWN  = Direction{0, 1}
	LEFT  = Direction{-1, 0}
	RIGHT = Direction{1, 0}
)

type Map struct {
	Cells []Cells

	Width  int
	Height int
}

type Cells []Cell
type Cell struct{}

type LinkedList struct {
	head  *Node
	nodes []Node
}
type Node struct {
	leftNode  *Node
	rightNode *Node

	Value    int
	Children []Node
	Split    bool
}

type Stone struct {
	Value int
}

var stoneCache map[string]uint

func CountStones(startingStones []Stone, blinks int) uint {
	if blinks == 0 {
		return uint(len(startingStones))
	}

	count := uint(0)
	for _, stone := range startingStones {
		lookup := fmt.Sprintf("%d|%d", stone.Value, blinks)
		nextStones, ok := stoneCache[lookup]
		if ok {
			count += nextStones
			continue
		}

		if stone.Value == 0 {
			stones := []Stone{
				{Value: 1},
			}
			c := CountStones(stones, blinks-1)
			stoneCache[lookup] = c
			count += c
		} else if IsEvenDigitNum(stone.Value) {
			left, right := SplitEvenDigitString(stone.Value)

			stones := []Stone{
				{Value: left},
				{Value: right},
			}
			c := CountStones(stones, blinks-1)
			stoneCache[lookup] = c
			count += c
		} else {
			stones := []Stone{
				{Value: stone.Value * 2024},
			}
			c := CountStones(stones, blinks-1)
			stoneCache[lookup] = c
			count += c
		}
	}
	return count
}

func (n *Node) Render() string {
	if n.Split {
		return fmt.Sprintf("%s %s", n.Children[0].Render(), n.Children[1].Render())
	} else {
		return fmt.Sprintf("%d", n.Value)
	}
}

func (n *LinkedList) Render() string {
	first := true

	rendered := ""
	for _, node := range n.nodes {
		if first {
			first = false
		} else {
			rendered += " "
		}

		rendered += node.Render()
	}

	return rendered
}

func (l *LinkedList) Count(blinks int) uint {
	count := uint(0)
	for i := range l.nodes {
		value := l.nodes[i].Value
		count += CountStones([]Stone{{Value: value}}, blinks)
	}

	return count
}

// 0 becomes 1
// if even length, replaced by 2 stone
// if not activated it is *2024
func (n *Node) Blink() {
	if n.Split {
		for i := range n.Children {
			n.Children[i].Blink()
		}
	} else if IsEvenDigitNum(n.Value) {
		left, right := SplitEvenDigitString(n.Value)

		n.Children = make([]Node, 2)

		n.Children[0].Value = left
		n.Children[1].Value = right

		n.Split = true
	} else if n.Value == 0 {
		n.Value = 1
	} else {
		n.Value *= 2024
	}
}

func (n *LinkedList) Blink() {
	wg := sync.WaitGroup{}
	for i := range n.nodes {
		wg.Add(1)

		go func() {
			n.nodes[i].Blink()
			wg.Done()
		}()
	}
	wg.Wait()
}

func IsEvenDigitNum(number int) bool {
	return len(fmt.Sprintf("%d", number))%2 == 0
}

var splitCache = map[int][]int{}

func SplitEvenDigitString(number int) (int, int) {
	if v, ok := splitCache[number]; ok {
		return v[0], v[1]
	}

	r := fmt.Sprintf("%d", number)

	start := 0
	mid := len(r) / 2

	leftPart := r[start:mid]
	rightPart := r[mid:]

	left, err := strconv.Atoi(leftPart)
	if err != nil {
		panic(err)
	}
	right, err := strconv.Atoi(rightPart)
	if err != nil {
		panic(err)
	}

	return left, right
}

func LoadNodeList(data []byte) *LinkedList {
	nodes := []Node{}

	for _, num := range bytes.Split(data, []byte(" ")) {
		n, err := strconv.Atoi(string(num))
		if err != nil {
			panic(err)
		}
		node := &Node{
			Value:    n,
			Children: make([]Node, 2),
		}

		nodes = append(nodes, *node)
	}
	return &LinkedList{head: &nodes[0], nodes: nodes}
}

func main() {
	adventofcode.Time(func() {
		data := adventofcode.LoadFile("2024/day_eleven/input.txt")

		n := LoadNodeList(data)

		//for i := range 75 {
		//	fmt.Printf("Loop %d\n", i)
		//
		//	n.Blink()
		//}

		// 54ms (part 2 code: 1.6ms)
		//fmt.Printf("Part 1: %d\n", n.Count(25))
		// 43ms
		fmt.Printf("Part 2: %d\n", n.Count(75))
	})
}

func init() {
	stoneCache = map[string]uint{}
}
