package martinez_rueda

import (
	"sort"
)

type Events []*SweepEvent

// Priority queue that holds sweep-events sorted from left to right.
type PriorityQueue struct {
	events Events
	sorted bool
}

func NewPriorityQueue() PriorityQueue {
	return PriorityQueue{
		events: []*SweepEvent{},
		sorted: false,
	}
}

func (pq PriorityQueue) sort() {
	sort.SliceStable(pq.events, func(i, j int) bool {
		return compareSweepEvents(pq.events[i], pq.events[j])
	})

}

func (pq *PriorityQueue) enqueue(event *SweepEvent) {
	if !pq.isSorted() {
		pq.events = append(pq.events, event)
		return
	}

	if len(pq.events) <= 0 {
		pq.events = append(pq.events, event)
		return
	}

	// priority queue is sorted, shift elements to the right and find place for event
	index := len(pq.events) - 1
	for idx := len(pq.events) - 1; idx >= 0 && compareSweepEvents(event, pq.events[idx]); idx-- {
		if idx+1 == len(pq.events) {
			pq.events = append(pq.events, pq.events[idx])
		} else {
			pq.events[idx+1] = pq.events[idx]
		}

		index = idx - 1
	}

	if index+1 == len(pq.events) {
		pq.events = append(pq.events, event)
	} else {
		pq.events[index+1] = event
	}

}

func (pq *PriorityQueue) dequeue() *SweepEvent {
	if !pq.isSorted() {
		pq.sort()
		pq.sorted = true
	}

	ep := len(pq.events) - 1
	e := pq.events[ep]
	pq.events = pq.events[:ep]
	return e
}

func (pq *PriorityQueue) isSorted() bool {
	return pq.sorted
}

func (pq *PriorityQueue) isEmpty() bool {
	if len(pq.events) == 0 {
		return true
	}
	return false
}
