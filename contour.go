package martinez_rueda

import (
	"github.com/paulmach/orb"
	"math"
)

type Contour struct {
	points        []orb.Point
	holes         []int
	isExternal    bool
	cc            bool
	precomputedCC *bool
}

func NewContour(points []orb.Point) Contour {
	return Contour{
		points: points,
	}
}

func (c *Contour) Add(point orb.Point) {
	c.points = append(c.points, point)
}

func (c *Contour) erase(index int) {
	//         unset($this->points[$index]);
	c.points = append(c.points[:index], c.points[(index+1):]...)
}

func (c *Contour) clear() {
	c.points = []orb.Point{}
	c.holes = []int{}
}

func (c *Contour) addHole(index int) {
	c.holes = append(c.holes, index)
}

// Get the p-th vertex of the external contour
func (c *Contour) vertex(index int) orb.Point {
	return c.points[index]
}

func (c *Contour) segment(index int) Segment {
	if index == c.nvertices()-1 {
		return NewSegment(c.points[len(c.points)-1], c.points[0])
	}
	return NewSegment(c.points[index], c.points[index+1])
}

func (c *Contour) hole(index int) int {
	return c.holes[index]
}

// Get minimum bounding rectangle
// ['min' => Point, 'max' => Point]
func (c *Contour) getBoundingBox() []orb.Point {
	minX := math.Inf(1)
	minY := math.Inf(1)
	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for idx := 0; idx < len(c.points); idx++ {
		point := c.points[idx]
		if point.X() < minX {
			minX = point.X()
		}
		if point.X() > maxX {
			maxX = point.X()
		}
		if point.Y() < minY {
			minY = point.Y()
		}
		if point.Y() > maxY {
			maxY = point.Y()
		}
	}

	return []orb.Point{orb.Point{minX, minY}, orb.Point{maxX, maxY}}
}

func (c *Contour) counterClockwise() bool {
	if c.precomputedCC != nil {
		return *c.precomputedCC
	}
	precomputedCC := true
	c.precomputedCC = &precomputedCC

	var area float64
	for idx := 0; idx < len(c.points)-1; idx++ {
		area = area + c.vertex(idx).X()*c.vertex(idx+1).Y() - c.vertex(idx+1).X()*c.vertex(idx).Y()
	}

	area = area + c.vertex(len(c.points)-1).X()*c.vertex(0).Y() - c.vertex(0).X()*c.vertex(len(c.points)-1).Y()

	c.cc = area >= 0.0

	return c.cc
}

func (c *Contour) clockwise() bool {
	return !c.counterClockwise()
}

func (c *Contour) changeOrientation() {
	c.points = segmentsReverse(c.points)
	c.cc = !c.cc
}

func (c *Contour) setClockwise() {
	if c.counterClockwise() {
		c.changeOrientation()
	}
}

func (c *Contour) setCounterClockwise() {
	if c.clockwise() {
		c.changeOrientation()
	}
}

func (c *Contour) external() bool {
	return c.isExternal
}

func (c *Contour) setExternal(flag bool) {
	c.isExternal = flag
}

func (c *Contour) move(x, y float64) {
	for idx := 0; idx < len(c.points); idx++ {
		c.points[idx] = orb.Point{c.points[idx].X() + x, c.points[idx].Y() + y}

	}
}

func (c *Contour) nvertices() int {
	return len(c.points)
}

func (c *Contour) GetPoint(index int ) orb.Point {
	return c.points[index]
}
func (c *Contour) Nvertices() int {
	return c.nvertices()
}