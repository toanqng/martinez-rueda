package martinez_rueda

import (
	"github.com/paulmach/orb"
	"math"
)

const earthRadius = 6371008.8 // 6,378,137 meters).

func toRadians(deg float64) float64 { return deg * math.Pi / 180 }
func toDegrees(rad float64) float64 { return rad * 180 / math.Pi }

// DestinationPoint return the destination from a point based on a distance and bearing.
// Given a start point and a distance d along constant bearing θ, this will calculate the destina­tion point. If you maintain a constant bearing along a rhumb line, you will gradually spiral in towards one of the poles.
func DestinationPoint(point orb.Point, meters, bearingDegrees float64) orb.Point {
	// see http://williams.best.vwh.net/avform.htm#LL
	δ := meters / earthRadius // angular distance in radians
	θ := toRadians(bearingDegrees)
	φ1 := toRadians(point.Lat())
	λ1 := toRadians(point.Lon())
	φ2 := math.Asin(math.Sin(φ1)*math.Cos(δ) + math.Cos(φ1)*math.Sin(δ)*math.Cos(θ))
	λ2 := λ1 + math.Atan2(math.Sin(θ)*math.Sin(δ)*math.Cos(φ1), math.Cos(δ)-math.Sin(φ1)*math.Sin(φ2))
	λ2 = math.Mod(λ2+3*math.Pi, 2*math.Pi) - math.Pi // normalise to -180..+180°

	return orb.Point{toDegrees(λ2), toDegrees(φ2)}
}

// Signed area of the triangle (p0, p1, p2)
func signedArea(p0 orb.Point, p1 orb.Point, p2 orb.Point) float64 {
	//         return ($p0->x - $p2->x) * ($p1->y - $p2->y) - ($p1->x - $p2->x) * ($p0->y - $p2->y);
	return (p0.X()-p2.X())*(p1.Y()-p2.Y()) - (p1.X()-p2.X())*(p0.Y()-p2.Y())
}

/**
 * Used for sorting SweepEvents in PriorityQueue
 * If same x coordinate - from bottom to top.
 * If two endpoints share the same point - rights are before lefts.
 * If two left endpoints share the same point then they must be processed
 * in the ascending order of their associated edges in SweepLine
 */
func compareSweepEvents(event1 *SweepEvent, event2 *SweepEvent) bool {
	// x is not the same
	// $event1->p->x > $event2->p->x
	if event1.p.X() > event2.p.X() {
		return true
	}

	// x is not the same too
	// ($event2->p->x > $event1->p->x)
	if event2.p.X() > event1.p.X() {
		return false
	}

	// x is the same, but y is not
	// the event with lower y-coordinate is processed first
	// (!$event1->p->equalsTo($event2->p))
	if !event1.p.Equal(event2.p) {
		return event1.p.Y() > event2.p.Y()
	}

	// x and y are the same, but one is a left endpoint and the other a right endpoint
	// the right endpoint is processed first

	if event1.isLeft != event2.isLeft {
		return event1.isLeft
	}

	// x and y are the same and both points are left or right
	return event1.above(event2.other.p)
}

func compareSegments(event1 *SweepEvent, event2 *SweepEvent) bool {
	// ($event1->equalsTo($event2))
	if event1.equalsTo(event2) {
		return false
	}

	// (self::signedArea($event1->p, $event1->other->p, $event2->p) != 0
	//            || self::signedArea($event1->p, $event1->other->p, $event2->other->p) != 0) {
	if signedArea(event1.p, event1.other.p, event2.p) != 0 ||
		signedArea(event1.p, event1.other.p, event2.other.p) != 0 {
		//            if ($event1->p->equalsTo($event2->p)) {
		//                return $event1->below($event2->other->p)
		if event1.p.Equal(event2.p) {
			return event1.below(event2.other.p)
		}

		if compareSweepEvents(event1, event2) {
			return event2.above(event1.p)
			//return $event1->below($event2->p);
		}

		return event1.below(event2.p)
		//return $event2->above($event1->p);
	}

	if event1.p.Equal(event2.p) {
		//return $event1->lessThan($event2);
		return false
	}

	return compareSweepEvents(event1, event2)
}
