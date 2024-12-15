package main

import (
	"adventofcode"
	"bytes"
	"fmt"
	"slices"
)

func main() {
	adventofcode.Time(func() {
		data := adventofcode.LoadFile("2024/day_fifteen/input.txt")
		factory := NewFactory(data, false)

		RunAutomation(factory)
		gpsTotal := factory.Score()

		fmt.Printf("Part 1: %d\n", gpsTotal)
	})
	adventofcode.Time(func() {
		data := adventofcode.LoadFile("2024/day_fifteen/input.txt")
		factory := NewFactory(data, true)

		RunAutomation(factory)
		gpsTotal := factory.Score()

		fmt.Printf("Part 2: %d\n", gpsTotal)
	})
}
func RunAutomation(factory *Factory) {
	robot := factory.GetRobot()
	if robot == nil {
		panic("no robot")
	}

	for robot.Step(factory) {
		// running
	}
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

func (f *Factory) Score() int {
	score := 0
	for y, line := range f.Objects {
		for x, obj := range line {
			if obj == nil {
				continue
			}
			if obj.GetType() == 'O' || obj.GetType() == '[' {
				score += 100*y + x
			}
		}
	}

	return score
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

func NewFactory(data []byte, isBig bool) *Factory {
	config := bytes.Split(data, []byte("\n\n"))

	logistics, strategy := config[0], config[1]

	lines := bytes.Split(logistics, []byte("\n"))
	factory := Factory{
		Width:  len(lines[0]),
		Height: len(lines),
	}

	if isBig {
		factory.Width *= 2
	}

	factory.Objects = make([][]Objector, len(lines))
	for y, line := range lines {
		if !isBig {
			factory.Objects[y] = make([]Objector, len(line))
		} else {
			factory.Objects[y] = make([]Objector, len(line)*2)
		}
		for x, b := range line {
			if isBig {
				// big mode stretches everything by 2
				x = x * 2
			}

			if b == '#' {
				factory.Objects[y][x] = &Wall{}
				if isBig {
					factory.Objects[y][x+1] = &Wall{}
				}
			} else if b == 'O' {
				if isBig {
					factory.Objects[y][x] = &BigBox{
						Left: true,
						Mover: Mover{
							StartingX: x,
							StartingY: y,
							X:         x,
							Y:         y,
						},
						Pair: &BigBox{
							Mover: Mover{
								StartingX: x + 1,
								StartingY: y,
								X:         x + 1,
								Y:         y,
							},
						},
					}
					factory.Objects[y][x+1] = factory.Objects[y][x].(*BigBox).Pair
					factory.Objects[y][x].(*BigBox).Pair.Pair = factory.Objects[y][x].(*BigBox)
				} else {
					factory.Objects[y][x] = &Box{
						Mover{
							StartingX: x,
							StartingY: y,
							X:         x,
							Y:         y,
						},
					}
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
				if isBig {
					// add space after robot
					factory.Objects[y][x+1] = &Space{}
				}

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
				if isBig {
					factory.Objects[y][x+1] = &Space{}
				}
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
	for range 2 {
		object := factory.Get(m.X+direction[0], m.Y+direction[1])
		switch obj := object.(type) {
		case *Wall:
			return false
		case *BigBox:
			return obj.CanMove(factory, direction)
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
		case *BigBox:
			if !obj.Move(obj, factory, direction) {
				return false
			}
		case *Box:
			// attempt to push box
			if !obj.Move(obj, factory, direction) {
				return false
			}
		case *Space:
			// switch with space
			factory.Objects[m.Y][m.X], factory.Objects[m.Y+direction[1]][m.X+direction[0]] = obj, self
			//fmt.Printf("Moved %d,%d (%c) to %d,%d (%c)\n", m.X, m.Y, self.GetType(), m.X+direction[0], m.Y+direction[1], obj.GetType())

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

type BigBox struct {
	Mover

	Left bool
	Pair *BigBox
}

func (b *BigBox) GetType() byte {
	if b.Left {
		return '['
	}
	return ']'
}

func (b *BigBox) CanMove(factory *Factory, direction Direction) bool {
	if slices.Equal(direction, LEFT) && !b.Left {
		return b.Pair.canMove(factory, direction)
	} else if slices.Equal(direction, RIGHT) && b.Left {
		return b.Pair.canMove(factory, direction)
	} else if slices.Equal(direction, LEFT) && b.Left {
		return b.canMove(factory, direction)
	} else if slices.Equal(direction, RIGHT) && !b.Left {
		return b.canMove(factory, direction)
	}

	return b.canMove(factory, direction) && b.Pair.canMove(factory, direction)
}

func (b *BigBox) canMove(factory *Factory, direction Direction) bool {
	return b.Mover.CanMove(factory, direction)
}

func (b *BigBox) Move(self Objector, factory *Factory, direction Direction) bool {
	if !b.Pair.canMove(factory, direction) {
		return false
	}

	if slices.Equal(direction, UP) || slices.Equal(direction, DOWN) {
		return b.Mover.Move(self, factory, direction) && b.Pair.Mover.Move(b.Pair, factory, direction)
	}

	return b.Mover.Move(self, factory, direction)
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

	if !r.Move(r, factory, currentStep) {
		// failed to move
	}

	return true

}
