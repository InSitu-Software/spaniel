package spaniel

import (
	"fmt"
	"sort"
	"time"
)

// Span represents a basic span, with a start and end time.
type Span interface {
	Start() time.Time
	End() time.Time
	String() string
}

// Spans represents a list of spans, on which other functions operate.
type Spans []Span

// ByStart sorts a list of spans by their start point
type ByStart Spans

func (s ByStart) Len() int           { return len(s) }
func (s ByStart) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByStart) Less(i, j int) bool { return s[i].Start().Before(s[j].Start()) }

// ByEnd sorts a list of spans by their end point
type ByEnd Spans

func (s ByEnd) Len() int           { return len(s) }
func (s ByEnd) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByEnd) Less(i, j int) bool { return s[i].End().Before(s[j].End()) }

// UnionHandlerFunc is used by UnionWithHandler to allow for custom functionality when two spans are merged.
// It is passed the two spans to be merged, and span which will result from the union.
type UnionHandlerFunc func(mergeInto, mergeFrom, mergeSpan Span) Span

// IntersectionHandlerFunc is used by IntersectionWithHandler to allow for custom functionality when two spans
// intersect. It is passed the two spans that intersect, and span representing the intersection.
type IntersectionHandlerFunc func(intersectingEvent1, intersectingEvent2, intersectionSpan Span) Span

func filter(spans Spans, filterFunc func(Span) bool) Spans {
	filtered := Spans{}
	for _, span := range spans {
		if !filterFunc(span) {
			filtered = append(filtered, span)
		}
	}
	return filtered
}

// getMax gets the latest of two timestamps, returning the latter if they are equal
func getMax(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}

	return b
}

// getMax gets the earliest of two timestamps, returning the latter if they are equal
func getMin(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}

	return b
}

// IsInstant returns true if the interval is deemed instantaneous
func IsInstant(a Span) bool {
	return a.Start().Equal(a.End())
}

// Without returns a list of spans representing the spans of baseSpan not
// covered by intersector.
// It can return one Span element in Spans, if the intersector can overlap at the
// beginning or end of baseSpan (baseSpan.Start/End is moved to intersector.End/Start) or
// It returns two Span elements in Spans, if the intersector is within the baseSpan. this creates
// a Span baseSpan.Start->interseector.Start and intersector.End->baseSpan.End
func Without(a, b Span) Spans {
	residues := Spans{}

	baseStart := a.Start()
	baseEnd := a.End()
	interStart := b.Start()
	interEnd := b.End()

	// We do not want to split a Span with a no-duration time span
	if IsInstant(b) {
		residues = append(residues, a)
		return residues
	}

	switch {
	case baseStart.Before(interStart) && baseEnd.After(interEnd):
		// intersector is "in" baseSpan -> two new spans will be created
		// ----++++++---- a
		// ------++------ b
		// =>
		// ----++--++----
		baseSpanPart1 := New(
			a.Start(),
			b.Start(),
		)

		baseSpanPart2 := New(
			b.End(),
			a.End(),
		)

		residues = append(residues, baseSpanPart1, baseSpanPart2)

	case (baseStart.After(interStart) || baseStart.Equal(interStart)) && baseEnd.After(interEnd) && interEnd.After(baseStart):
		// intersector overlaps at the begin of basespan
		// ----++++++---- a
		// ---++++------- b
		//
		// -------+++----
		baseSpanPart := New(
			b.End(),
			a.End(),
		)

		residues = append(residues, baseSpanPart)

	case (baseEnd.Before(interEnd) || baseEnd.Equal(interEnd)) && baseStart.Before(interStart) && interStart.Before(baseEnd):
		// intersector intersects at the end of basespan
		// ----++++++---- a
		// -------++++--- b
		//
		// ----+++-------
		baseSpanPart := New(
			a.Start(),
			b.Start(),
		)

		residues = append(residues, baseSpanPart)

	case baseStart.Equal(interStart) && baseEnd.Equal(interEnd):
		// if basespan and intersection match, return empty Spans
		break

	case !interStart.After(baseStart) && !interEnd.Before(baseEnd):
		// if intersection engulfs the basespan, return empty Spans
		break
	default:
		residues = append(residues, a)
	}

	return residues
}

type WithoutHandlerFunc func(a, b Span, diff Spans) Spans

func WithoutWithHandler(a, b Span, handlerFunc WithoutHandlerFunc) Spans {
	s := Without(a, b)
	return handlerFunc(a, b, s)
}

// Within returns if b is completely in a
// Same instants of start or end are considered within.
func Within(a, b Span) bool {
	return ((a.Start().Before(b.Start())) || a.Start().Equal(b.Start())) &&
		((a.End().After(b.End())) || a.End().Equal(b.End()))
}

// Returns true if two spans overlap
func overlap(a, b Span) bool {
	return a.Start().Before(b.End()) && b.Start().Before(a.End())
}

// IntersectionWithHandler returns a list of Spans representing the overlaps between the contained spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned, with the span C covering
// the intersection of the A and B. The provided handler function is notified of the two spans that have been found
// to overlap, and the span representing the overlap.
func (s Spans) IntersectionWithHandler(intersectHandlerFunc IntersectionHandlerFunc) Spans {
	var sorted Spans
	sorted = append(sorted, s...)
	sort.Stable(ByStart(sorted))

	actives := Spans{sorted[0]}

	intersections := Spans{}

	for _, b := range sorted[1:] {
		// Tidy up the active span list
		actives = filter(actives, func(t Span) bool {
			// If this value is identical to one in actives, don't filter it.
			if b.Start().Equal(t.Start()) && b.End().Equal(t.End()) {
				return false
			}
			// If this value starts after the one in actives finishes, filter the active.
			return b.Start().After(t.End())
		})

		for _, a := range actives {
			if overlap(a, b) {
				spanStart := getMax(a.Start(), b.Start())
				spanEnd := getMin(a.End(), b.End())

				span := New(spanStart, spanEnd)
				intersection := intersectHandlerFunc(a, b, span)
				intersections = append(intersections, intersection)
			}
		}
		actives = append(actives, b)
	}
	return intersections
}

// Intersection returns a list of Spans representing the overlaps between the contained spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned,
// with the span C covering the intersection of A and B.
func (s Spans) Intersection() Spans {
	return s.IntersectionWithHandler(func(intersectingEvent1, intersectingEvent2, intersectionSpan Span) Span {
		return intersectionSpan
	})
}

// IntersectionBetweenWithHandler returns a list of pointers to Spans representing the overlaps between the contained spans
// and a given set of spans. It calls intersectHandlerFunc for each pair of spans that are intersected.
func (s Spans) IntersectionBetweenWithHandler(candidates Spans, intersectHandlerFunc IntersectionHandlerFunc) Spans {
	intersections := Spans{}
	for _, candidate := range candidates {
		for _, span := range s {
			i := Spans{candidate, span}.IntersectionWithHandler(func(a, b, s Span) Span {
				if a == candidate {
					return intersectHandlerFunc(a, b, s)
				}

				return intersectHandlerFunc(b, a, s)
			})
			intersections = append(intersections, i...)
		}
	}
	return intersections
}

// IntersectionBetween returns the slice of spans representing the overlaps between the contained spans
// and a given set of spans.
func (s Spans) IntersectionBetween(b Spans) Spans {
	return s.IntersectionBetweenWithHandler(b, func(intersectingEvent1, intersectingEvent2, intersectionSpan Span) Span {
		return intersectionSpan
	})
}

//Without removes the given Span from the Spans
func (s Spans) Without(b Span) Spans {
	var o = Spans{}
	for _, a := range s {
		o = append(o, Without(a, b)...)
	}

	return o
}

func (s Spans) WithoutWithHandler(b Span, handlerFunc WithoutHandlerFunc) Spans {
	var o Spans
	for _, a := range s {
		o = append(o, WithoutWithHandler(a, b, handlerFunc)...)
	}

	return o
}

//Duration sums up the duration of all given Spans
func (s Spans) Duration() time.Duration {
	var d time.Duration

	for _, span := range s {
		sd := span.End().Sub(span.Start())
		d += sd
	}

	return d
}

func (s Spans) String() string {
	var out string
	for _, span := range s {
		out = fmt.Sprintf("%s\n %s", out, span.String())
	}

	return out
}
