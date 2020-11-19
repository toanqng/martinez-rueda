package martinez_rueda

import "github.com/paulmach/orb"

type Segment struct {
	p1 orb.Point
	p2 orb.Point
}

func NewSegment(p1 orb.Point, p2 orb.Point) Segment {
	return Segment{
		p1: p1,
		p2: p2,
	}
}

func (s *Segment) setBegin(p orb.Point) {
	s.p1 = p
}

func (s *Segment) setEnd(p orb.Point) {
	s.p2 = p
}

func (s *Segment) begin() orb.Point {
	return s.p1
}

func (s *Segment) end() orb.Point {
	return s.p2
}

func (s *Segment) changeOrientation() {
	tmp := s.p1
	s.p1 = s.p2
	s.p2 = tmp
}
