package main

import (
	"adventofcode"
	"adventofcode/toolbox/fs"
	"bytes"
	"fmt"
	"slices"
)

// initial thoughts are to walk the lowest scoring path each step until we find the end
// we will need to do analysis on each intersection and send out paths at each step.

// Maze will be the current state of the maze
type Maze [][]byte
type Position struct {
	X, Y int
}

func NewMaze(data []byte) (Maze, Position, Position) {
	lines := bytes.Split(data, []byte("\n"))
	var start Position
	var end Position

	m := make(Maze, len(lines))
	for y, line := range lines {
		m[y] = make([]byte, len(line))
		for x, b := range line {
			if b == 'S' {
				start.X = x
				start.Y = y
			} else if b == 'E' {
				end.X = x
				end.Y = y
			}
			m[y][x] = b
		}
	}

	m.RemoveDeadEnds()

	// replace start with .
	m[start.Y][start.X] = '*'

	return m, start, end
}

func (m Maze) Solve(start Position, end Position) (int, int) {
	reindeer := m.DeployReindeer(start)

	// march lowest scoring reindeer until it's position is the end
	for true {
		// locate lowest scoring reindeer
		lowestX := -1
		var lowest *Reindeer
		for x, r := range reindeer {
			if r.X == end.X && r.Y == end.Y {
				fmt.Printf("Found lowest scoring reindeer at %d, %d (score=%d) after %d steps\n", r.X, r.Y, r.Score, r.Steps)
				for _, friend := range r.Friends {
					for _, p := range friend.Path {
						m[p.Y][p.X] = 'O'
					}
				}
				for _, p := range r.Path {
					m[p.Y][p.X] = 'O'
				}

				var score int
				for _, line := range m {
					for _, b := range line {
						if b == 'O' {
							score++
						}
					}
				}

				return r.Score, score
			}
			if lowest == nil || r.Score < lowest.Score {
				lowestX = x
				lowest = r
			}
		}

		if lowest == nil {
			panic("lowest is nil")
		}

		if !lowest.Move(m) {
			for _, r := range reindeer {
				// this doesn't work as the main branch
				if r.X == lowest.X && r.Y == lowest.Y && ((r.Score == lowest.Score && slices.Equal(r.Facing, lowest.Facing)) || r.Score-1000 == lowest.Score && !slices.Equal(r.Facing, lowest.Facing)) {
					r.Friends = append(r.Friends, lowest)

					for _, friend := range r.Friends {
						if friend == lowest {
							continue
						}
					}
				}
			}
			// remove from pool
			reindeer = append(reindeer[:lowestX], reindeer[lowestX+1:]...)
			continue
		}

		//fmt.Println(m)

		if lowest.Check(m, true) == '.' {
			reindeer = append(reindeer, lowest.Clone().Turn(true))
		}
		if lowest.Check(m, false) == '.' {
			reindeer = append(reindeer, lowest.Clone().Turn(false))
		}
	}

	panic("unreachable")
}

func (m Maze) String() string {
	s := ""
	for _, line := range m {
		s += fmt.Sprintf("%s\n", line)
	}
	return s
}

func (m Maze) Get(x, y int) byte {
	if x < 0 || x >= len(m) || y < 0 || y >= len(m[0]) {
		return 0
	}
	return m[y][x]
}

func (m Maze) RemoveDeadEnds() {
	for y := range m {
		for x := range m[y] {
			m.removeDeadEnd(x, y)
		}
	}
}

func (m Maze) DeployReindeer(start Position) []*Reindeer {
	reindeer := make([]*Reindeer, len(DIRECTIONS))

	reindeer[0] = &Reindeer{
		X:      start.X,
		Y:      start.Y,
		Facing: RIGHT,
		Path: []Position{
			{start.X, start.Y},
		},
		Friends: []*Reindeer{},
	}
	reindeer[1] = reindeer[0].Clone().Turn(true)
	reindeer[2] = reindeer[0].Clone().Turn(false)
	reindeer[3] = reindeer[1].Clone().Turn(true)

	return reindeer
}

func (m Maze) removeDeadEnd(x, y int) {
	cell := m.Get(x, y)

	if cell != '.' {
		return
	}

	var nextPath Direction
	walls := 0
	// check if 3 sides are walls and if they are make position a wall
	for _, d := range DIRECTIONS {
		nextCell := m.Get(x+d[0], y+d[1])

		if nextCell == '#' {
			walls++
		} else if nextCell == '.' {
			nextPath = d
		}
	}

	if walls == 3 {
		// turn current position to a wall try next path
		m[y][x] = '#'
		if nextPath != nil {
			m.removeDeadEnd(x+nextPath[0], y+nextPath[1])
		}
	}
}

type Reindeer struct {
	Score int
	Steps int

	X, Y   int
	Facing Direction

	Path    []Position
	Friends []*Reindeer
}

func (r *Reindeer) Clone() *Reindeer {
	friends := make([]*Reindeer, len(r.Friends))
	for i, f := range r.Friends {
		friends[i] = f
	}
	path := make([]Position, len(r.Path))
	copy(path, r.Path)

	return &Reindeer{
		Score:   r.Score,
		Steps:   r.Steps,
		X:       r.X,
		Y:       r.Y,
		Facing:  r.Facing,
		Path:    path,
		Friends: friends,
	}
}

func (r *Reindeer) Move(m Maze) bool {
	r.Steps++
	r.X += r.Facing[0]
	r.Y += r.Facing[1]

	r.Score += 1

	cell := m.Get(r.X, r.Y)

	switch cell {
	case '*':
		// potentially we are on a cheaper path but crossing another one
		if m.Get(r.X+r.Facing[0], r.Y+r.Facing[1]) == '.' {
			r.Path = append(r.Path, Position{r.X, r.Y})
			return true
		}
		fallthrough
	case '#':
		return false
	case '.':
		m[r.Y][r.X] = '*'
	}
	r.Path = append(r.Path, Position{r.X, r.Y})
	return true
}

func (r *Reindeer) Check(m Maze, left bool) byte {
	if slices.Equal(r.Facing, UP) {
		if left {
			return m.Get(r.X+LEFT[0], r.Y+LEFT[1])
		} else {
			return m.Get(r.X+RIGHT[0], r.Y+RIGHT[1])
		}
	} else if slices.Equal(r.Facing, DOWN) {
		if left {
			return m.Get(r.X+RIGHT[0], r.Y+RIGHT[1])
		} else {
			return m.Get(r.X+LEFT[0], r.Y+LEFT[1])
		}
	} else if slices.Equal(r.Facing, LEFT) {
		if left {
			return m.Get(r.X+DOWN[0], r.Y+DOWN[1])
		} else {
			return m.Get(r.X+UP[0], r.Y+UP[1])
		}
	} else if slices.Equal(r.Facing, RIGHT) {
		if left {
			return m.Get(r.X+UP[0], r.Y+UP[1])
		} else {
			return m.Get(r.X+DOWN[0], r.Y+DOWN[1])
		}
	}
	panic("unreachable")
}

func (r *Reindeer) Turn(left bool) *Reindeer {
	r.Score += 1000
	if slices.Equal(r.Facing, UP) {
		if left {
			r.Facing = LEFT
		} else {
			r.Facing = RIGHT
		}
	} else if slices.Equal(r.Facing, DOWN) {
		if left {
			r.Facing = RIGHT
		} else {
			r.Facing = LEFT
		}
	} else if slices.Equal(r.Facing, LEFT) {
		if left {
			r.Facing = DOWN
		} else {
			r.Facing = UP
		}
	} else if slices.Equal(r.Facing, RIGHT) {
		if left {
			r.Facing = UP
		} else {
			r.Facing = DOWN
		}
	}
	return r
}

type Direction []int

var (
	UP    = Direction{0, -1}
	DOWN  = Direction{0, 1}
	LEFT  = Direction{-1, 0}
	RIGHT = Direction{1, 0}

	DIRECTIONS = []Direction{UP, DOWN, LEFT, RIGHT}
)

func main() {
	adventofcode.Time(func() {
		data := fs.LoadFile("2024/day_sixteen/input.txt")
		m, start, end := NewMaze(data)
		score, seats := m.Solve(start, end)

		fmt.Printf("Part 1: %d\n", score)
		fmt.Printf("Part 2: %d\n", seats)
	})
}
