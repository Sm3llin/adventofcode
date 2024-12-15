package main

import (
	"adventofcode"
	"bytes"
	"fmt"
)

func main() {
	adventofcode.Time(func() {
		data := adventofcode.LoadFile("2024/day_fifteen/input.txt")
		factory := NewFactory(data)

		robot := factory.GetRobot()
		if robot == nil {
			panic("no robot")
		}

		for robot.Step(factory) {
			// running
		}

		gpsTotal := 0
		for y, row := range factory.Objects {
			for x, obj := range row {
				if obj == nil {
					continue
				}
				if obj.GetType() == 'O' {
					gpsTotal += 100*y + x
				}
			}
		}

		fmt.Printf("Part 1: %d\n", gpsTotal)
	})
}

type Direction []int

var (
	UP    = Direction{0, -1}
	DOWN  = Direction{0, 1}
	LEFT  = Direction{-1, 0}
	RIGHT = Direction{1, 0}
)

type Factory struct {
	Width, Height int
	Objects       [][]Objector
}

func (f Factory) Render() string {
	rendered := ""
	for _, line := range f.Objects {
		for _, obj := range line {
			if obj == nil {
				panic("nil object")
			}
			rendered += string(obj.GetType())
		}
		rendered += "\n"
	}
	return rendered
}

func (f Factory) Get(x, y int) Objector {
	if x < 0 || x >= f.Width || y < 0 || y >= f.Height {
		return &Wall{}
	}
	return f.Objects[y][x]
}

func (f *Factory) GetRobot() *Robot {
	for _, line := range f.Objects {
		for _, obj := range line {
			if obj == nil {
				continue
			}
			if obj.GetType() == '@' {
				return obj.(*Robot)
			}
		}
	}
	return nil
}

func NewFactory(data []byte) *Factory {
	config := bytes.Split(data, []byte("\n\n"))

	logistics, strategy := config[0], config[1]

	lines := bytes.Split(logistics, []byte("\n"))
	factory := Factory{
		Width:  len(lines[0]),
		Height: len(lines),
	}

	factory.Objects = make([][]Objector, len(lines))
	for y, line := range lines {
		factory.Objects[y] = make([]Objector, len(line))
		for x, b := range line {
			if b == '#' {
				factory.Objects[y][x] = &Wall{}
			} else if b == 'O' {
				factory.Objects[y][x] = &Box{
					Mover{
						StartingX: x,
						StartingY: y,
						X:         x,
						Y:         y,
					},
				}
			} else if b == '@' {
				robot := Robot{
					Mover: Mover{
						StartingX: x,
						StartingY: y,
						X:         x,
						Y:         y,
					},
					CurrentStep: 0,
					Steps:       []Direction{},
				}
				factory.Objects[y][x] = &robot

				// apply the strategy possible steps <,>,^,v
				for _, steps := range bytes.Split(strategy, []byte("\n")) {
					directions := make([]Direction, len(steps))
					for xS, step := range steps {
						switch step {
						case '<':
							directions[xS] = LEFT
						case '>':
							directions[xS] = RIGHT
						case '^':
							directions[xS] = UP
						case 'v':
							directions[xS] = DOWN
						}
					}
					robot.Steps = append(robot.Steps, directions...)
				}
			} else {
				factory.Objects[y][x] = &Space{}
			}
		}
	}
	return &factory
}

type Objector interface {
	GetType() byte
}

type Mover struct {
	X, Y                 int
	StartingX, StartingY int
}

func (m *Mover) CanMove(factory *Factory, direction Direction) bool {
	for true {
		object := factory.Get(m.X+direction[0], m.Y+direction[1])
		switch obj := object.(type) {
		case *Wall:
			return false
		case *Box:
			return obj.CanMove(factory, direction)
		case *Space:
			return true
		}
	}
	return false

}

func (m *Mover) Move(self Objector, factory *Factory, direction Direction) bool {
	for true {
		object := factory.Get(m.X+direction[0], m.Y+direction[1])
		switch obj := object.(type) {
		case *Wall:
			return false
		case *Box:
			// attempt to push box
			if !obj.Move(obj, factory, direction) {
				return false
			}
		case *Space:
			// switch with space
			factory.Objects[m.Y][m.X], factory.Objects[m.Y+direction[1]][m.X+direction[0]] = obj, self

			m.X += direction[0]
			m.Y += direction[1]

			return true
		}
	}
	return false
}

type Box struct {
	Mover
}

func (b Box) GetType() byte {
	return 'O'
}

type Space struct {
}

func (s Space) GetType() byte {
	return '.'
}

type Wall struct {
}

func (w Wall) GetType() byte {
	return '#'
}

type Robot struct {
	Mover

	CurrentStep int
	Steps       []Direction
}

func (r Robot) GetType() byte {
	return '@'
}

func (r *Robot) Step(factory *Factory) bool {
	if r.CurrentStep >= len(r.Steps) {
		return false
	}

	// obtain the instruction
	currentStep := r.Steps[r.CurrentStep]
	r.CurrentStep++

	r.Move(r, factory, currentStep)

	return true

}
