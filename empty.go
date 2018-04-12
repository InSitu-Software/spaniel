package spaniel

import (
	"time"
)

// Empty represents a simple span of time, with no additional properties. It should be constructed with NewEmpty.
type Empty struct {
	start     time.Time
	end       time.Time
	startType IntervalType
	endType   IntervalType
}

// Start returns the start time of a span
func (ets *Empty) Start() time.Time { return ets.start }

// End returns the end time of a span
func (ets *Empty) End() time.Time { return ets.end }

// StartType returns the type of the start of the interval (Open in this case)
func (ets *Empty) StartType() IntervalType { return ets.startType }

// EndType returns the type of the end of the interval (Closed in this case)
func (ets *Empty) EndType() IntervalType { return ets.endType }

// NewEmpty creates a span with just a start and end time, and associated types, and is used when no handlers are provided to Union or Intersection.
func NewEmpty(start time.Time, end time.Time, startType IntervalType, endType IntervalType) *Empty {
	return &Empty{start, end, startType, endType}
}

// NewEmptyTyped creates a span with a start and end time, with the types set to [] for instants and [) for spans.
func NewEmptyTyped(start time.Time, end time.Time) *Empty {
	if start.Equal(end) {
		// An instantaneous event has to be Closed (i.e. inclusive)
		return NewEmpty(start, end, Closed, Closed)
	}
	return NewEmpty(start, end, Closed, Open)
}

func (ets *Empty) String() string {
	s := ""
	if ets.StartType() == Closed {
		s += "["
	} else {
		s += "("
	}
	s += ets.Start().String()
	s += ","
	s += ets.End().String()

	if ets.EndType() == Closed {
		s += "]"
	} else {
		s += ")"
	}
	return s
}
