package grid

import (
	"adventofcode/toolbox/arrays"
	"slices"
)

type Maze[T comparable] struct {
	Grid[T]

	Start, End Position
	Wall       T
	OnMove     MoveFunc

	SearchedGrid Grid[bool]

	walkers []*mazeWalker[T]
}

func (m *Maze[T]) register(w *mazeWalker[T]) {
	m.walkers = append(m.walkers, w)
}

func (m *Maze[T]) unregister(w *mazeWalker[T]) {
	idx := slices.Index(m.walkers, w)

	if idx == -1 {
		return
	}

	m.walkers = append(m.walkers[:idx], m.walkers[idx+1:]...)
}

type MoveFunc func(from, to Position) (allow bool, score int)

func NewMaze[T comparable](g Grid[T], wall T, onMove MoveFunc) *Maze[T] {
	return &Maze[T]{
		Grid:         g,
		Wall:         wall,
		OnMove:       onMove,
		SearchedGrid: NewGridValue(false, g.Width, g.Height),
	}
}

type mazeWalker[T comparable] struct {
	position Position
	path     []Position

	Step int
}

func (w *mazeWalker[T]) walk(m *Maze[T]) bool {
	options := []Position{}
	for _, d := range ConnectedDirections {
		// check all non explored locations
		nextPosition := w.position.Move(d)
		cell, err := m.Get(nextPosition.X, nextPosition.Y)
		searched, errSearch := m.SearchedGrid.Get(nextPosition.X, nextPosition.Y)
		searched = searched || errSearch != nil

		if err != nil {
			continue
		} else if cell != m.Wall && !searched {
			options = append(options, nextPosition)
		}
	}

	if len(options) == 0 {
		m.unregister(w)
		return false
	}
	w.position = options[0]
	w.path = append(w.path, w.position)
	w.Step++
	// TODO: can't have this if it is the end?
	m.SearchedGrid.Set(w.position.X, w.position.Y, true)

	for _, option := range options[1:] {
		nWalker := &mazeWalker[T]{
			position: option,
			Step:     w.Step,
			path:     make([]Position, len(w.path)),
		}
		m.SearchedGrid.Set(option.X, option.Y, true)
		copy(nWalker.path, w.path)
		nWalker.path[len(nWalker.path)-1] = option

		m.register(nWalker)
	}
	return true
}

func (m *Maze[T]) Solve(start, end Position) ([]Position, bool) {
	m.Start = start
	m.End = end

	m.walkers = []*mazeWalker[T]{
		{
			position: m.Start,
			Step:     0,
			path:     []Position{m.Start},
		},
	}
	m.SearchedGrid.Set(m.Start.X, m.Start.Y, true)

	successWalkers := []*mazeWalker[T]{}

	for len(m.walkers) > 0 {
		walkers := m.walkers

		slices.SortFunc(walkers, func(a, b *mazeWalker[T]) int {
			aScore := m.End.X - a.position.X + m.End.Y - a.position.Y
			bScore := m.End.X - b.position.X + m.End.Y - b.position.Y

			return aScore - bScore
		})

		for _, w := range m.walkers {
			if !w.walk(m) && len(m.walkers) == 0 {
				//fmt.Println("last known path", w.path)
			}
		}

		// check if a walker is at the end
		for _, w := range m.walkers {
			// let them all finish
			if w.position.Equal(m.End) {
				successWalkers = append(successWalkers, w)
				m.unregister(w)
			}
		}
	}
	if len(successWalkers) > 0 {
		return successWalkers[0].path, true
	}

	return nil, false
}

func (m *Maze[T]) FloodFill(start Position) Grid[int] {
	type tracker struct {
		position Position
		step     int
	}
	q := arrays.Queue[tracker]{}
	q.Push(tracker{
		position: start,
		step:     0,
	})

	// create a copy of maze as an int grid
	floodGrid := NewGridValue(-1, m.Width, m.Height)

	for t := range q.Iter() {
		cell, err := m.Get(t.position.X, t.position.Y)
		if err != nil {
			continue
		}
		if cell == m.Wall {
			continue
		}
		floodGrid.Set(t.position.X, t.position.Y, t.step)
		for _, d := range ConnectedDirections {
			nextPosition := t.position.Move(d)

			current, err := floodGrid.Get(nextPosition.X, nextPosition.Y)
			if err != nil {
				continue
			}
			if current != -1 && current <= t.step {
				continue
			}

			q.Push(tracker{
				position: nextPosition,
				step:     t.step + 1,
			})
		}
	}
	return floodGrid
}
