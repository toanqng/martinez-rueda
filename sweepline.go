package martinez_rueda

type SweepLine struct {
	events []*SweepEvent
}

func NewSweepLine() SweepLine {
	return SweepLine{
		events: []*SweepEvent{},
	}
}

func (sl *SweepLine) size() int {
	return len(sl.events)
}

func (sl *SweepLine) get(index int) *SweepEvent {
	return sl.events[index]
}

func (sl *SweepLine) remove(removable *SweepEvent) {
	for idx, event := range sl.events {
		if event.equalsTo(removable) {
			//append(s[:index], s[index+1:]...)
			sl.events = append(sl.events[:idx], sl.events[(idx+1):]...)
			break
		}
	}
}

func (sl *SweepLine) insert(event *SweepEvent) int {
	if len(sl.events) == 0 {
		sl.events = append(sl.events, event)
		return 0
	}

	// priority queue is sorted, shift elements to the right and find place for event
	index := len(sl.events) - 1
	for i := len(sl.events) - 1; i >= 0 && compareSegments(event, sl.events[i]); i-- {
		if i+1 == len(sl.events) {
			sl.events = append(sl.events, sl.events[i])
		} else {
			sl.events[i+1] = sl.events[i]
		}

		index = i - 1
	}

	if index+1 == len(sl.events) {
		sl.events = append(sl.events, event)
	} else {
		sl.events[index+1] = event
	}

	return index + 1
}
