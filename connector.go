package martinez_rueda

import (
	"github.com/paulmach/orb"
)

type Connector struct {
	openPolygons   []*PointChain
	closedPolygons []*PointChain
	closed         bool
}

func NewConnector() Connector {
	return Connector{
		openPolygons:   []*PointChain{},
		closedPolygons: []*PointChain{},
		closed:         false,
	}
}

func (c *Connector) isClosed() bool {
	return c.closed
}

func (c *Connector) add(segment Segment) {
	size := len(c.openPolygons)
	for jdx := 0; jdx < size; jdx++ {
		chain := c.openPolygons[jdx]
		isLinkSegment := chain.linkSegment(segment)

		if !isLinkSegment {
			continue
		}

		if chain.closed {

			if len(chain.segments) == 2 {
				chain.closed = false
				return
			}
			c.closedPolygons = append(c.closedPolygons, c.openPolygons[jdx])

			//append(s[:index], s[index+1:]...)
			c.openPolygons = append(c.openPolygons[:jdx], c.openPolygons[(jdx+1):]...)

			return
		}

		// if chain not closed
		k := len(c.openPolygons)

		for idx := jdx + 1; idx < k; idx++ {
			v := c.openPolygons[idx]
			if chain.linkChain(v) {
				//append(s[:index], s[index+1:]...)
				c.openPolygons = append(c.openPolygons[:idx], c.openPolygons[(idx+1):]...)

				return
			}
		}

		return
	}

	newChain := NewPointChain(segment)
	c.openPolygons = append(c.openPolygons, newChain)
}

func (c *Connector) toPolygon() *Polygon {
	contours := []Contour{}
	for _, cp := range c.closedPolygons {
		contourPoints := []orb.Point{}
		for _, point := range cp.segments {
			contourPoints = append(contourPoints, point)
		}

		// close contour
		first := contourPoints[0]
		last := contourPoints[len(contourPoints)-1]
		if first.Equal(last) == false {
			contourPoints = append(contourPoints, first)
		}

		contours = append(contours, NewContour(contourPoints))
	}

	return NewPolygon(contours)
}
