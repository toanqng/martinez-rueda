package martinez_rueda

import (
	"errors"
	"github.com/paulmach/orb"
	"math"
)

type OPERATION string
type EDGE_TYPE uint
type POLYGON_TYPE uint

const (
	OP_INTERSECTION OPERATION = "INTERSECTION"
	OP_UNION        OPERATION = "UNION"
	OP_DIFFERENCE   OPERATION = "DIFFERENCE"
	OP_XOR          OPERATION = "XOR"
)

const (
	POLYGON_TYPE_SUBJECT  POLYGON_TYPE = 1
	POLYGON_TYPE_CLIPPING POLYGON_TYPE = 2
)

const (
	EDGE_NORMAL               EDGE_TYPE = 1
	EDGE_NON_CONTRIBUTING     EDGE_TYPE = 2
	EDGE_SAME_TRANSITION      EDGE_TYPE = 3
	EDGE_DIFFERENT_TRANSITION EDGE_TYPE = 4
)

var eq PriorityQueue

func Compute(subject *Polygon, clipping *Polygon, operation OPERATION) (result *Polygon) {

	eq = NewPriorityQueue()

	// Test for 1 trivial result case
	if subject.ncontours()*clipping.ncontours() == 0 {
		if subject.ncontours() == 0 {
			return clipping
		}
		return subject
	}

	// Test 2 for trivial result case
	boxSub := subject.getBoundingBox()
	minSubj := boxSub[0]
	maxSubj := boxSub[1]

	boxClip := clipping.getBoundingBox()
	minClip := boxClip[0]
	maxClip := boxClip[1]

	if minSubj.X() > maxClip.X() || minClip.X() > maxSubj.X() || minSubj.Y() > maxClip.Y() || minClip.Y() > maxSubj.Y() {
		result = subject

		for idx := 0; idx < clipping.ncontours(); idx++ {
			result.contours = append(result.contours, clipping.contour(idx))
		}
		return result
	}

	// Boolean operation is not trivial
	// Insert all the endpoints associated to the line segments into the event queue
	for idx := 0; idx < subject.ncontours(); idx++ {
		con := subject.contour(idx)
		for jdx := 0; jdx < con.nvertices(); jdx++ {
			processSegment(con.segment(jdx), POLYGON_TYPE_SUBJECT)
		}
	}

	for idx := 0; idx < clipping.ncontours(); idx++ {
		con := clipping.contour(idx)
		for jdx := 0; jdx < con.nvertices(); jdx++ {
			processSegment(con.segment(jdx), POLYGON_TYPE_CLIPPING)
		}
	}

	connector := NewConnector()
	sweepline := NewSweepLine()

	minMaxX := math.Min(maxSubj.X(), maxClip.X())

	glbIDX := 0
	for !eq.isEmpty() {
		e := eq.dequeue()

		if (operation == OP_INTERSECTION && (e.p.X() > minMaxX)) || (operation == OP_DIFFERENCE && (e.p.X() > maxSubj.X())) {
			result = connector.toPolygon()
			return result
		}

		if e.isLeft {
			position := sweepline.insert(e)
			var prev *SweepEvent
			var next *SweepEvent

			if position > 0 {
				prev = sweepline.get(position - 1)
			}

			if position < sweepline.size()-1 {
				next = sweepline.get(position + 1)
			}

			if prev == nil {
				e.inside = false
				e.inOut = false
			} else if prev.edgeType != EDGE_NORMAL {
				if position-2 < 0 {
					e.inside = false
					e.inOut = false

					if prev.polygonType != e.polygonType {
						e.inside = true
					} else {
						e.inOut = true
					}
				} else {
					prev2 := sweepline.get(position - 2)

					if prev.polygonType == e.polygonType {
						e.inOut = !prev.inOut
						e.inside = !prev2.inOut
					} else {
						e.inOut = !prev2.inOut
						e.inside = !prev.inOut
					}
				}
			} else if e.polygonType == prev.polygonType {
				e.inside = prev.inside
				e.inOut = !prev.inOut
			} else {
				e.inside = !prev.inOut
				e.inOut = prev.inside
			}

			if next != nil {
				possibleIntersection(e, next)
			}

			if prev != nil {
				possibleIntersection(prev, e)
			}
		} else {
			var prev *SweepEvent
			var next *SweepEvent
			// not left, the line segment must be removed from S
			other_pos := -1

			for index, event := range sweepline.events {
				if event.equalsTo(e.other) {
					other_pos = index
					break
				}
			}

			if other_pos != -1 {
				if other_pos > 0 {
					prev = sweepline.get(other_pos - 1)
				}

				if other_pos < len(sweepline.events)-1 {
					next = sweepline.get(other_pos + 1)
				}
			}

			// Check if the line segment belongs to the Boolean operation
			switch e.edgeType {
			case EDGE_NORMAL:
				switch operation {
				case OP_INTERSECTION:
					if e.other.inside {
						connector.add(e.segment())
					}

					break

				case OP_UNION:
					if !e.other.inside {

						connector.add(e.segment())

					}

					break

				case OP_DIFFERENCE:

					if e.polygonType == POLYGON_TYPE_SUBJECT && !e.other.inside || e.polygonType == POLYGON_TYPE_CLIPPING && e.other.inside {
						connector.add(e.segment())
					}

					break

				case OP_XOR:

					connector.add(e.segment())
					break
				}

				break // end of EDGE_NORMAL

			case EDGE_SAME_TRANSITION:

				if operation == OP_INTERSECTION || operation == OP_UNION {
					connector.add(e.segment())
				}

				break

			case EDGE_DIFFERENT_TRANSITION:

				if operation == OP_DIFFERENCE {
					connector.add(e.segment())
				}

				break
			} // end switch (e.edgeType)

			if other_pos != -1 {
				sweepline.remove(sweepline.get(other_pos))
			}

			if next != nil && prev != nil {
				possibleIntersection(next, prev)
			}

		}

		glbIDX += 1

	}

	return connector.toPolygon()
}

func processSegment(segment Segment, polygonType POLYGON_TYPE) {
	// if the two edge endpoints are equal the segment is discarded
	if segment.begin().Equal(segment.end()) {
		return
	}

	e1 := NewSweepEvent(segment.begin(), true, polygonType, nil, EDGE_NORMAL)
	e2 := NewSweepEvent(segment.end(), true, polygonType, e1, EDGE_NORMAL)
	e1.other = e2

	if e1.p.X() < e2.p.X() {
		e2.isLeft = false
	} else if e1.p.X() > e2.p.X() {
		e1.isLeft = false
	} else if e1.p.Y() < e2.p.Y() {
		e2.isLeft = false
	} else {
		e1.isLeft = false
	}

	eq.enqueue(e1)
	eq.enqueue(e2)

}

func possibleIntersection(event1 *SweepEvent, event2 *SweepEvent) error {
	// uncomment these two lines if self-intersecting polygons are not allowed
	// if (event1.polygon_type == event2.polygon_type) {
	//    return false;
	// }

	ip1, _, intersections := findIntersection(event1.segment(), event2.segment())

	if intersections == 0 {
		return nil
	}

	if intersections == 1 && (event1.p.Equal(event2.p) || event1.other.p.Equal(event2.other.p)) {
		return nil
	}

	// the line segments overlap, but they belong to the same polygon
	// the program does not work with this kind of polygon
	if intersections == 2 && event1.polygonType == event2.polygonType {
		return errors.New("Polygon has overlapping edges.")
	}

	if intersections == 1 {
		if !event1.p.Equal(ip1) && !event1.other.p.Equal(ip1) {
			divideSegment(event1, ip1)
		}

		if !event2.p.Equal(ip1) && !event2.other.p.Equal(ip1) {
			divideSegment(event2, ip1)
		}

		return nil
	}

	// The line segments overlap
	sorted_events := []*SweepEvent{}

	if event1.p.Equal(event2.p) {
		//   $sorted_events[] = 0;
		sorted_events = []*SweepEvent{}
	} else if compareSweepEvents(event1, event2) {
		sorted_events = append(sorted_events, event2)
		sorted_events = append(sorted_events, event1)
	} else {
		sorted_events = append(sorted_events, event1)
		sorted_events = append(sorted_events, event2)
	}

	if event1.other.p.Equal(event2.other.p) {
		//   $sorted_events[] = 0;
		sorted_events = []*SweepEvent{}
	} else if compareSweepEvents(event1.other, event2.other) {
		sorted_events = append(sorted_events, event2.other)
		sorted_events = append(sorted_events, event1.other)
	} else {
		sorted_events = append(sorted_events, event1.other)
		sorted_events = append(sorted_events, event2.other)
	}

	if len(sorted_events) == 2 {
		event1.edgeType = EDGE_NON_CONTRIBUTING
		event1.other.edgeType = EDGE_NON_CONTRIBUTING
		if event1.inOut == event2.inOut {
			event2.edgeType = EDGE_SAME_TRANSITION
			event2.other.edgeType = EDGE_SAME_TRANSITION
		} else {
			event2.edgeType = EDGE_DIFFERENT_TRANSITION
			event2.other.edgeType = EDGE_DIFFERENT_TRANSITION
		}

		return nil
	}

	if len(sorted_events) == 3 {
		sorted_events[1].edgeType = EDGE_NON_CONTRIBUTING
		sorted_events[1].other.edgeType = EDGE_NON_CONTRIBUTING

		if sorted_events[0].other == nil {
			if event1.inOut == event2.inOut {
				sorted_events[0].other.edgeType = EDGE_SAME_TRANSITION
			} else {
				sorted_events[0].other.edgeType = EDGE_DIFFERENT_TRANSITION
			}
		} else {
			if event1.inOut == event2.inOut {
				sorted_events[2].other.edgeType = EDGE_SAME_TRANSITION
			} else {
				sorted_events[2].other.edgeType = EDGE_DIFFERENT_TRANSITION
			}
		}

		if sorted_events[0].other == nil {
			divideSegment(sorted_events[2].other, sorted_events[1].p)
		} else {
			divideSegment(sorted_events[0], sorted_events[1].p)
		}
		return nil
	}

	// NEED DEBUG
	if len(sorted_events) == 4 {
		if !sorted_events[0].equalsTo(sorted_events[3].other) {
			sorted_events[1].edgeType = EDGE_NON_CONTRIBUTING
			if event1.inOut == event2.inOut {
				sorted_events[2].edgeType = EDGE_SAME_TRANSITION
			} else {
				sorted_events[2].edgeType = EDGE_DIFFERENT_TRANSITION
			}

			divideSegment(sorted_events[0], sorted_events[1].p)
			divideSegment(sorted_events[1], sorted_events[2].p)

			return nil
		}

		sorted_events[1].edgeType = EDGE_NON_CONTRIBUTING
		sorted_events[1].other.edgeType = EDGE_NON_CONTRIBUTING
		divideSegment(sorted_events[0], sorted_events[1].p)

		if event1.inOut == event2.inOut {
			sorted_events[3].other.edgeType = EDGE_SAME_TRANSITION
		} else {
			sorted_events[3].other.edgeType = EDGE_DIFFERENT_TRANSITION
		}

		divideSegment(sorted_events[3].other, sorted_events[2].p)
	}

	return nil
}

func findIntersection(segment0 Segment, segment1 Segment) (orb.Point, orb.Point, int) {
	pi0 := orb.Point{}
	pi1 := orb.Point{}

	p0 := segment0.begin()
	d0 := orb.Point{segment0.end().X() - p0.X(), segment0.end().Y() - p0.Y()}

	p1 := segment1.begin()
	d1 := orb.Point{segment1.end().X() - p1.X(), segment1.end().Y() - p1.Y()}

	sqrwEpsilon := 1e-7 // it was 1e-3 before
	E := orb.Point{p1.X() - p0.X(), p1.Y() - p0.Y()}
	kross := d0.X()*d1.Y() - d0.Y()*d1.X()
	sqr_kross := kross * kross
	sqr_len0 := d0.X()*d0.X() + d0.Y()*d0.Y()
	sqr_len1 := d1.X()*d1.X() + d1.Y()*d1.Y()

	if sqr_kross > sqrwEpsilon*sqr_len0*sqr_len1 {
		s := (E.X()*d1.Y() - E.Y()*d1.X()) / kross

		if s < 0 || s > 1 {
			return pi0, pi1, 0
		}

		t := (E.X()*d0.Y() - E.Y()*d0.X()) / kross

		if t < 0 || t > 1 {
			return pi0, pi1, 0
		}

		// intersection of lines is a point an each segment
		pi0 = orb.Point{p0.X() + s*d0.X(), p0.Y() + s*d0.Y()}

		if distanceTo(pi0, segment0.begin()) < 1e-8 {
			pi0 = segment0.begin()
		}

		if distanceTo(pi0, segment0.end()) < 1e-8 {
			pi0 = segment0.end()
		}

		if distanceTo(pi0, segment1.begin()) < 1e-8 {
			pi0 = segment1.begin()
		}

		if distanceTo(pi0, segment1.end()) < 1e-8 {
			pi0 = segment1.end()
		}
		return pi0, pi1, 1
	}

	sqr_len_e := E.X()*E.X() + E.Y()*E.Y()
	kross = E.X()*d0.Y() - E.Y()*d0.X()
	sqr_kross = kross * kross

	if sqr_kross > sqrwEpsilon*sqr_len0*sqr_len_e {
		return pi0, pi1, 0
	}

	s0 := (d0.X()*E.X() + d0.Y()*E.Y()) / sqr_len0
	s1 := s0 + (d0.X()*d1.X()+d0.Y()*d1.Y())/sqr_len0

	smin := math.Min(s0, s1)
	smax := math.Max(s0, s1)

	w, imax := findIntersection2(0.0, 1.0, smin, smax)

	if imax > 0 {
		pi0 = orb.Point{p0.X() + w[0]*d0.X(), p0.Y() + w[0]*d0.Y()}

		if distanceTo(pi0, segment0.begin()) < 1e-8 {
			pi0 = segment0.begin()
		}

		if distanceTo(pi0, segment0.end()) < 1e-8 {
			pi0 = segment0.end()
		}

		if distanceTo(pi0, segment1.begin()) < 1e-8 {
			pi0 = segment1.begin()
		}

		if distanceTo(pi0, segment1.end()) < 1e-8 {
			pi0 = segment1.end()
		}

		if imax > 1 {
			pi1 = orb.Point{p0.X() + w[1]*d0.X(), p0.Y() + w[1]*d0.Y()}
		}
	}

	return pi0, pi1, imax
}

func distanceTo(p1, p2 orb.Point) float64 {
	dx := p1.X() - p2.X()
	dy := p1.Y() - p2.Y()

	return math.Sqrt(dx*dx + dy*dy)
}

func findIntersection2(u0, u1, v0, v1 float64) ([]float64, int) {
	w := []float64{}
	if u1 < v0 || u0 > v1 {
		return w, 0
	}

	if u1 > v0 {
		if u0 < v1 {
			if u0 < v0 {
				// w[0] = v0
				w = append(w, v0)
			} else {
				//w[0] = u0
				w = append(w, u0)
			}

			if u1 > v1 {
				// w[1] = v1
				w = append(w, v1)
			} else {
				//w[1] = u1
				w = append(w, u1)
			}

			return w, 2
		} else {
			// w[0] = u0
			w = append(w, u0)

			return w, 1
		}
	} else {
		// w[0] = u1
		w = append(w, u1)
		return w, 1
	}
}

func divideSegment(event *SweepEvent, point orb.Point) {
	right := NewSweepEvent(point, false, event.polygonType, event, event.edgeType)
	left := NewSweepEvent(point, true, event.polygonType, event.other, event.other.edgeType)

	if compareSweepEvents(left, event.other) {
		event.other.isLeft = true
		left.isLeft = false
	}

	event.other.other = left
	event.other = right

	eq.enqueue(left)
	eq.enqueue(right)
}
