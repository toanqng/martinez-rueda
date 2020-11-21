package martinez_rueda

import (
	"github.com/paulmach/orb"
)

type PointChain struct {
	segments []orb.Point
	closed   bool
}

func NewPointChain(initSegment Segment) *PointChain {
	return &PointChain{
		segments: []orb.Point{initSegment.begin(), initSegment.end()},
	}
}

func (pc *PointChain) begin() orb.Point {
	return pc.segments[0]
}
func (pc *PointChain) end() orb.Point {
	return pc.segments[len(pc.segments)-1]
}

func (pc *PointChain) linkSegment(segment Segment) bool {

	front := pc.begin()
	back := pc.end()

	// CASE 1
	if segment.begin().Equal(front) {
		if segment.end().Equal(back) {
			pc.closed = true
		} else {
			//  Prepend one  elements to the beginning of an array
			pc.segments = append([]orb.Point{segment.end()}, pc.segments...)
		}

		return true
	}

	// CASE 2
	if segment.end().Equal(back) {
		if segment.begin().Equal(front) {
			pc.closed = true
		} else {
			pc.segments = append(pc.segments, segment.begin())
		}
		return true
	}
	// CASE 3
	if segment.end().Equal(front) {
		if segment.begin().Equal(back) {
			pc.closed = true
		} else {
			//  Prepend one  elements to the beginning of an array
			pc.segments = append([]orb.Point{segment.begin()}, pc.segments...)
		}
		return true
	}

	// CASE 4
	if segment.begin().Equal(back) {
		if segment.end().Equal(front) {
			pc.closed = true
		} else {
			pc.segments = append(pc.segments, segment.end())
		}
		return true
	}
	return false
}

func (pc *PointChain) linkChain(other *PointChain) bool {
	front := pc.begin()
	back := pc.end()
	otherFront := other.begin()
	otherBack := other.end()

	if otherFront.Equal(back) {
		// Shift an element off the beginning of array
		other.segments = other.segments[1:]
		// insert at the end of $this->segments
		pc.segments = append(pc.segments, other.segments...)
		return true
	}
	if otherBack.Equal(front) {
		// Shift an element off the beginning of array
		pc.segments = pc.segments[1:]
		// insert at the beginning of $this->segments
		pc.segments = append(other.segments, pc.segments...)
		return true
	}

	if otherFront.Equal(front) {

		// Shift an element off the beginning of array
		pc.segments = pc.segments[1:]
		other.segments = segmentsReverse(other.segments)
		// insert reversed at the beginning of $this->segments
		pc.segments = append(other.segments, pc.segments...)
		return true
	}

	if otherBack.Equal(back) {
		// array_pop — Pop the element off the end of array
		pc.segments = pc.segments[:(len(pc.segments) - 1)]
		other.segments = segmentsReverse(other.segments)
		// insert reversed at the end of $this->segments
		pc.segments = append(pc.segments, other.segments...)

		return true

	}

	return false
}

func segmentsReverse(segments []orb.Point) []orb.Point {
	for i, j := 0, len(segments)-1; i < j; i, j = i+1, j-1 {
		segments[i], segments[j] = segments[j], segments[i]
	}
	return segments
}

func (pc *PointChain) Size() int {
	return len(pc.segments)
}
