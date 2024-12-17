package grid

import (
	"fmt"
	"iter"
	"reflect"
)

type Position struct {
	X int
	Y int
}
type Grid[T any] struct {
	Data [][]T

	Height int
	Width  int
}

func NewGrid[T any](data [][]T) Grid[T] {
	var width int
	if len(data) >= 1 {
		width = len(data[0])
	}
	return Grid[T]{
		Data:   data,
		Height: len(data),
		Width:  width,
	}
}

func (g Grid[T]) All() iter.Seq2[Position, T] {
	var y, x int

	return func(yield func(Position, T) bool) {
		for y = 0; y < g.Height; y++ {
			for x = 0; x < g.Width; x++ {
				if !yield(Position{x, y}, g.Data[y][x]) {
					return
				}
			}
		}
	}
}

func (g Grid[T]) Get(x, y int) (T, error) {
	if !g.CheckBounds(x, y) {
		return reflect.Zero(reflect.TypeOf(*new(T))).Interface().(T), fmt.Errorf("out of bounds: %d, %d", x, y)
	}
	return g.Data[y][x], nil
}

func (g Grid[T]) CheckBounds(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

func (g Grid[T]) Set(x, y int, value T) {
	g.Data[y][x] = value
}

func (g Grid[T]) Swap(xA, yA, xB, yB int) {
	g.Data[yA][xA], g.Data[yB][xB] = g.Data[yB][xB], g.Data[yA][xA]
}

func (g Grid[T]) String() string {
	return fmt.Sprintf("Grid[%d, %d]", g.Height, g.Width)
}

func (g Grid[T]) Render() string {
	var s string
	for _, row := range g.Data {
		for _, cell := range row {
			s += fmt.Sprintf("%v", cell)
		}
		s += "\n"
	}
	return s
}
