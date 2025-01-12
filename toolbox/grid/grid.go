package grid

import (
	"fmt"
	"iter"
	"reflect"
)

// Direction is an array of 2 items [x, y]
type Direction []int

var (
	N  = Direction{0, -1}
	NE = Direction{1, -1}
	E  = Direction{1, 0}
	SE = Direction{1, 1}
	S  = Direction{0, 1}
	SW = Direction{-1, 1}
	W  = Direction{-1, 0}
	NW = Direction{-1, -1}

	ConnectedDirections = []Direction{N, E, S, W}
	DiagonalDirections  = []Direction{NE, SE, SW, NW}
	AllDirections       = []Direction{N, NE, E, SE, S, SW, W, NW}
)

type Position struct {
	X int
	Y int
}

func (p Position) Move(d Direction) Position {
	return Position{p.X + d[0], p.Y + d[1]}
}

func (p Position) Equal(o Position) bool {
	return p.X == o.X && p.Y == o.Y
}

func (p Position) Delta(other Position) (x int, y int) {
	x, y = other.X-p.X, other.Y-p.Y
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return
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

func NewGridValue[T any](value T, width, height int) Grid[T] {
	data := make([][]T, height)
	for i := range data {
		data[i] = make([]T, width)

		for j := range data[i] {
			data[i][j] = value
		}
	}

	return Grid[T]{
		Data:   data,
		Height: len(data),
		Width:  width,
	}
}

func (g Grid[T]) Clone() Grid[T] {
	data := make([][]T, len(g.Data))
	for i := range data {
		data[i] = make([]T, len(g.Data[i]))

		for j := range data[i] {
			data[i][j] = g.Data[i][j]
		}
	}
	return Grid[T]{data, g.Height, g.Width}
}

func (g Grid[T]) All() iter.Seq2[Position, T] {
	return func(yield func(Position, T) bool) {
		var y, x int
		for y = 0; y < g.Height; y++ {
			for x = 0; x < g.Width; x++ {
				if !yield(Position{x, y}, g.Data[y][x]) {
					return
				}
			}
		}
	}
}

// Neighbours gets all neighbours that are connected N, E, S, W
func (g Grid[T]) Neighbours(x, y int, directions []Direction) iter.Seq2[Position, T] {
	return func(yield func(Position, T) bool) {
		for _, d := range directions {
			nx, ny := x+d[0], y+d[1]
			if g.CheckBounds(nx, ny) {
				if !yield(Position{nx, ny}, g.Data[ny][nx]) {
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

func (g Grid[T]) RenderFunc(f func(v T) string) string {
	var s string
	for _, row := range g.Data {
		s += "|"
		for _, cell := range row {
			s += f(cell)
			s += "|"
		}
		s += "\n"
	}
	return s

}

func (g Grid[T]) Render() string {
	var s string
	for _, row := range g.Data {
		s += "|"
		for _, cell := range row {
			s += fmt.Sprintf("%v", cell)
			s += "|"
		}
		s += "\n"
	}
	return s
}

func (g *Grid[T]) FindAndReplace(find func(v T) bool, newValue T) (Position, error) {
	for p, cell := range g.All() {
		if find(cell) {
			g.Set(p.X, p.Y, newValue)
			return p, nil
		}
	}
	return Position{}, fmt.Errorf("not found")
}
