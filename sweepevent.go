package martinez_rueda

import "github.com/paulmach/orb"

type SweepEvent struct {
	// Point associated with the event
	p orb.Point
	// Event associated to the other endpoint of the edge
	other *SweepEvent
	// Is the point the left endpoint of the edge (p, other->p)
	isLeft bool
	// Indicates if the edge belongs to subject or clipping polygon
	polygonType POLYGON_TYPE
	// Inside-outside transition into the polygon
	inOut bool
	// Is the edge (p, other->p) inside the other polygon
	inside bool
	// Used for overlapped edges
	edgeType EDGE_TYPE
	// For sorting, increases monotonically
	id int
}

var SEID int

func NewSweepEvent(p orb.Point, isLeft bool, associatedPolygon POLYGON_TYPE, other *SweepEvent, edgeType EDGE_TYPE) *SweepEvent {
	SEID += 1
	return &SweepEvent{
		p:           p,
		other:       other,
		isLeft:      isLeft,
		polygonType: associatedPolygon,
		edgeType:    edgeType,
		id:          SEID,
	}
}

func (se *SweepEvent) getId() int {
	return se.id
}

func (se *SweepEvent) segment() Segment {
	return NewSegment(se.p, se.other.p)
}

func (se *SweepEvent) below(point orb.Point) bool {
	if se.isLeft {
		return signedArea(se.p, se.other.p, point) > 0
	}
	return signedArea(se.other.p, se.p, point) > 0
}

func (se *SweepEvent) above(point orb.Point) bool {
	return !se.below(point)
}

func (se *SweepEvent) equalsTo(event *SweepEvent) bool {
	return se.getId() == event.getId()
}

func (se *SweepEvent) lessThan(event SweepEvent) bool {
	return se.getId() < event.getId()
}
