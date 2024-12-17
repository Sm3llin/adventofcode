package main

import (
	"adventofcode"
	"adventofcode/toolbox/fs"
	"bytes"
	"fmt"
)

type Direction []int

var (
	UP    = Direction{0, -1}
	DOWN  = Direction{0, 1}
	LEFT  = Direction{-1, 0}
	RIGHT = Direction{1, 0}
)

type Region struct {
	Plots  []*Plot
	IsEdge map[string]int

	Map *Map
}

type Plot struct {
	Value byte

	X int
	Y int

	Region *Region
}

func (p *Plot) Fences() int {
	topPlot := p.Region.Map.Get(p.X, p.Y-1)
	bottomPlot := p.Region.Map.Get(p.X, p.Y+1)
	leftPlot := p.Region.Map.Get(p.X-1, p.Y)
	rightPlot := p.Region.Map.Get(p.X+1, p.Y)

	var fences int

	if p.X == 0 || p.X == p.Region.Map.Width-1 {
		fences++
	}
	if p.Y == 0 || p.Y == p.Region.Map.Height-1 {
		fences++
	}
	for _, plot := range []*Plot{topPlot, bottomPlot, leftPlot, rightPlot} {
		if plot == nil {
			continue
		}
		if plot.Value != p.Value || plot.Region != p.Region {
			fences++
		}
	}

	return fences
}

func (r *Region) GetBounds() (int, int, int, int) {
	var regionXMin, regionXMax, regionYMin, regionYMax int

	for _, plot := range r.Plots {
		if plot.X < regionXMin || regionXMin == 0 {
			regionXMin = plot.X
		}
		if plot.X > regionXMax {
			regionXMax = plot.X
		}
		if plot.Y < regionYMin || regionYMin == 0 {
			regionYMin = plot.Y
		}
		if plot.Y > regionYMax {
			regionYMax = plot.Y
		}
	}
	return regionXMin - 2, regionXMax + 2, regionYMin - 2, regionYMax + 2
}

func (r *Region) Edges() (edges int) {
	r.IsEdge = make(map[string]int)
	regionPlot := r.Plots[0]
	m := r.Map

	regionXMin, regionXMax, regionYMin, regionYMax := r.GetBounds()

	var topTrackingEdge bool
	var bottomTrackingEdge bool
	var p *Plot
	// check for horizontal edges
	for y := regionYMin; y <= regionYMax; y++ {
		for x := regionXMin; x <= regionXMax; x++ {
			p = m.Get(x, y)
			if p == nil || regionPlot.Region != p.Region {
				topTrackingEdge = false
				bottomTrackingEdge = false
				continue
			}

			topPlot := m.Get(p.X, p.Y-1)

			// remove non matching
			if topPlot == nil || topPlot.Region != p.Region {
				if !topTrackingEdge {
					edges++
					topTrackingEdge = true
					r.IsEdge[fmt.Sprintf("%d,%d", p.X, p.Y-1)] += 1
				}
			} else if topTrackingEdge && topPlot.Region == p.Region {
				topTrackingEdge = false
			}

			bottomPlot := m.Get(p.X, p.Y+1)

			// remove non matching
			if bottomPlot == nil || bottomPlot.Region != p.Region {
				if !bottomTrackingEdge {
					edges++
					bottomTrackingEdge = true
					r.IsEdge[fmt.Sprintf("%d,%d", p.X, p.Y+1)] += 1
				}
			} else if bottomTrackingEdge && bottomPlot.Region == p.Region {
				bottomTrackingEdge = false
			}
		}
	}

	bottomTrackingEdge = false
	topTrackingEdge = false

	// check for vertical edges
	for x := regionXMin; x <= regionXMax; x++ {
		for y := regionYMin; y <= regionYMax; y++ {
			p = m.Get(x, y)
			if p == nil || regionPlot.Region != p.Region {
				topTrackingEdge = false
				bottomTrackingEdge = false
				continue
			}

			leftPlot := m.Get(p.X-1, p.Y)

			// remove non matching
			if leftPlot == nil || leftPlot.Region != p.Region {
				if !topTrackingEdge {
					edges++
					topTrackingEdge = true
					r.IsEdge[fmt.Sprintf("%d,%d", p.X-1, p.Y)] += 1
				}
			} else if topTrackingEdge && leftPlot.Region == p.Region {
				topTrackingEdge = false

			}

			rightPlot := m.Get(p.X+1, p.Y)

			// remove non matching
			if rightPlot == nil || rightPlot.Region != p.Region {
				if !bottomTrackingEdge {
					edges++
					bottomTrackingEdge = true
					r.IsEdge[fmt.Sprintf("%d,%d", p.X+1, p.Y)] += 1
				}
			} else if bottomTrackingEdge && rightPlot.Region == p.Region {
				bottomTrackingEdge = false
			}
		}
	}

	return edges
}

func (p *Plot) Spread(d Direction) {
	nextPlot := p.Region.Map.Get(p.X+d[0], p.Y+d[1])

	if nextPlot != nil && nextPlot.Value == p.Value && nextPlot.Region == nil {
		nextPlot.Region = p.Region
		p.Region.Plots = append(p.Region.Plots, nextPlot)

		nextPlot.SpreadRegion()
	}
}

func (p *Plot) SpreadRegion() {
	// take our region and spread it
	for _, d := range []Direction{UP, DOWN, LEFT, RIGHT} {
		p.Spread(d)
	}
}

type Map struct {
	Regions []*Region
	Plots   [][]*Plot

	Height, Width int
}

func (m *Map) Get(x, y int) *Plot {
	// check bounds
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return nil
	}
	return m.Plots[y][x]
}

func (m *Map) FenceDiscountCost() int {
	var sum int

	for _, region := range m.Regions {
		edges := region.Edges()
		fmt.Printf(region.Render())
		fmt.Printf("plots=%d, edges=%d\n\n", len(region.Plots), edges)
		sum += edges * len(region.Plots)
	}

	return sum
}
func (m *Map) FenceCost() int {
	var sum int

	for _, region := range m.Regions {
		var fenceCount int
		for _, plot := range region.Plots {
			fenceCount += plot.Fences()
		}
		sum += fenceCount * len(region.Plots)
	}

	return sum
}

func FenceCost(data []byte) int {
	m := &Map{
		Regions: []*Region{},
		Plots:   [][]*Plot{},
	}

	var x, y int
	var line []byte
	for y, line = range bytes.Split(data, []byte("\n")) {
		plots := []*Plot{}
		var b byte
		for x, b = range line {
			plot := &Plot{
				X: x,
				Y: y,

				Value: b,
			}
			plots = append(plots, plot)
		}

		m.Plots = append(m.Plots, plots)
	}

	m.Height = y + 1
	m.Width = x + 1

	// TODO: Figure out regions
	for _, plots := range m.Plots {
		for _, plot := range plots {
			if plot.Region == nil {
				plot.Region = &Region{
					Map:   m,
					Plots: []*Plot{plot},
				}
				m.Regions = append(m.Regions, plot.Region)
				plot.SpreadRegion()
			}
		}
	}

	return m.FenceCost()
}

func FenceDiscount(data []byte) int {
	m := &Map{
		Regions: []*Region{},
		Plots:   [][]*Plot{},
	}

	var x, y int
	var line []byte
	for y, line = range bytes.Split(data, []byte("\n")) {
		plots := []*Plot{}
		var b byte
		for x, b = range line {
			plot := &Plot{
				X: x,
				Y: y,

				Value: b,
			}
			plots = append(plots, plot)
		}

		m.Plots = append(m.Plots, plots)
	}

	m.Height = y + 1
	m.Width = x + 1

	// TODO: Figure out regions
	for _, plots := range m.Plots {
		for _, plot := range plots {
			if plot.Region == nil {
				plot.Region = &Region{
					Map:   m,
					Plots: []*Plot{plot},
				}
				m.Regions = append(m.Regions, plot.Region)
				plot.SpreadRegion()
			}
		}
	}

	return m.FenceDiscountCost()
}

func (r *Region) Render() string {
	minX, maxX, minY, maxY := r.GetBounds()

	s := ""
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			plot := r.Map.Get(x, y)

			isEdge, ok := r.IsEdge[fmt.Sprintf("%d,%d", x, y)]
			if ok && isEdge > 0 {
				s += fmt.Sprintf("%d", isEdge)
			} else if plot == nil || plot.Region != r {
				s += "."
			} else {
				s += "o"
			}
		}
		s += "\n"
	}

	return s
}

func main() {
	adventofcode.Time(func() {

		data := fs.LoadFile("2024/day_twelve/input.txt")

		//cost := FenceCost(data)
		//fmt.Printf("Part 1: %d\n", cost)
		cost := FenceDiscount(data)
		fmt.Printf("Part 2: %d\n", cost)
	})
}
