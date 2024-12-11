package main

import (
	"adventofcode"
	"bytes"
	"fmt"
	"strconv"
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
	head *Node
}
type Node struct {
	leftNode  *Node
	rightNode *Node

	Value int
}

func (n *LinkedList) Render() string {
	node := n.head
	rendered := ""

	for node != nil {
		if rendered != "" {
			rendered += " "
		}
		rendered += strconv.Itoa(node.Value)
		node = node.rightNode
	}
	return rendered
}

func (l *LinkedList) Count() int {
	node := l.head
	count := 0

	for node != nil {
		count++
		node = node.rightNode
	}
	return count
}

// 0 becomes 1
// if even length, replaced by 2 stone
// if not activated it is *2024
func (n *LinkedList) Blink() {
	node := n.head

	for node != nil {
		if node.Value == 0 {
			node.Value = 1
			node = node.rightNode
		} else if IsEvenDigitNum(node.Value) {
			left, right := SplitEvenDigitString(node.Value)
			rightNode := &Node{
				rightNode: node.rightNode,
				Value:     right,
			}
			leftNode := &Node{
				leftNode:  node.leftNode,
				rightNode: rightNode,
				Value:     left,
			}
			rightNode.leftNode = leftNode

			if node.leftNode != nil {
				node.leftNode.rightNode = leftNode
			}
			if node.rightNode != nil {
				node.rightNode.leftNode = rightNode
			}

			if node.leftNode == nil {
				n.head = leftNode
			}
			node.leftNode = nil
			node.rightNode = nil

			node = rightNode.rightNode
		} else {
			node.Value *= 2024
			node = node.rightNode
		}

	}
}

func IsEvenDigitNum(number int) bool {
	return len(fmt.Sprintf("%d", number))%2 == 0
}

func SplitEvenDigitString(number int) (int, int) {
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
	var firstNode *Node
	var prevNode *Node

	for _, num := range bytes.Split(data, []byte(" ")) {
		n, err := strconv.Atoi(string(num))
		if err != nil {
			panic(err)
		}
		node := &Node{
			Value: n,
		}

		if firstNode == nil {
			firstNode = node
		} else {
			prevNode.rightNode = node
			node.leftNode = prevNode
		}
		prevNode = node
	}
	return &LinkedList{head: firstNode}
}

func main() {
	adventofcode.Time(func() {

		data := adventofcode.LoadFile("2024/day_eleven/input.txt")

		n := LoadNodeList(data)

		for i := range 75 {
			fmt.Printf("Loop %d\n", i)

			n.Blink()
		}

		// 54ms
		fmt.Printf("Part 1: %d\n", n.Count())
	})
}
