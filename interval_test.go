package spaniel

import (
	"reflect"
	"testing"
	"time"
)

var berlin, _ = time.LoadLocation("Europe/Berlin")

var overlapTests = []struct {
	description string
	begin       time.Time
	end         time.Time
	begin2      time.Time
	end2        time.Time
	expectation bool
}{
	{
		"Same start, b ends before a",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 17, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		true,
	},
	{
		"Same start, a ends before b",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 17, 6, 5, 0, berlin),
		true,
	},
	{
		"Same start, same end",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		true,
	},
	{
		"a starts before b, same end",
		time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		true,
	},
	{
		"b starts before a, same end",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		true,
	},
	{
		"different start, different end, with overlap",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		true,
	},
	{
		"a follows b directly",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 17, 34, 5, 0, berlin),
		false,
	},
	{
		"Same start, same end, different day",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 27, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 27, 16, 6, 5, 0, berlin),
		false,
	},
	{
		"a follows b directly but on a different day",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 27, 16, 6, 5, 0, berlin),
		time.Date(2020, 9, 27, 17, 34, 5, 0, berlin),
		false,
	},
	{
		"Overlapping the next day",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 27, 8, 6, 5, 0, berlin),
		time.Date(2020, 9, 27, 6, 4, 5, 0, berlin),
		time.Date(2020, 9, 27, 16, 6, 5, 0, berlin),
		true,
	},
	{
		"Overlapping the same day into the next",
		time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
		time.Date(2020, 9, 26, 22, 6, 5, 0, berlin),
		time.Date(2020, 9, 26, 20, 4, 5, 0, berlin),
		time.Date(2020, 9, 27, 16, 6, 5, 0, berlin),
		true,
	},
}

func TestOverlap(t *testing.T) {
	for _, tt := range overlapTests {
		t.Log(tt.description)
		spanA := New(tt.begin, tt.end)
		spanB := New(tt.begin2, tt.end2)
		hasOverlap := overlap(spanA, spanB)
		if tt.expectation != hasOverlap {
			t.Fail()
		}
	}
}

var intersectionTests = []struct {
	description   string
	spans         Spans
	intersections Spans
}{
	{
		"two spans, no intersection",
		Spans{
			New(
				time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 20, 45, 5, 0, berlin),
				time.Date(2020, 9, 26, 23, 6, 5, 0, berlin),
			),
		},
		Spans{},
	},
	{
		"two spans, b starts before a ends = b.Start -> a.End",
		Spans{
			New(
				time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 18, 45, 5, 0, berlin),
				time.Date(2020, 9, 26, 23, 6, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 18, 45, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
	},
	{
		"two spans, a starts before b ends = a.Start-> b.Ends",
		Spans{
			New(
				time.Date(2020, 9, 26, 18, 45, 5, 0, berlin),
				time.Date(2020, 9, 26, 23, 6, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 15, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 18, 45, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
	},
	{
		"two spans, a engulfs b completely = b",
		Spans{
			New(
				time.Date(2020, 9, 26, 16, 45, 5, 0, berlin),
				time.Date(2020, 9, 26, 23, 6, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
	},
	{
		"two spans, b engulfs a completely = b",
		Spans{
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 16, 45, 5, 0, berlin),
				time.Date(2020, 9, 26, 23, 6, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
	},
	{
		"two spans, b equals a = b",
		Spans{
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
	},
	{
		"two spans, b equals a = b",
		Spans{
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 17, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 19, 13, 5, 0, berlin),
			),
		},
	},
	{
		"three spans, all three disjunctive = empty",
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 16, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 18, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 20, 13, 5, 0, berlin),
			),
		},
		Spans{},
	},
	{
		"three spans, a intersects with c",
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 16, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 13, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
	},
	{
		"three spans, a intersects with b, b with c",
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 13, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 13, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 13, 5, 0, berlin),
			),
		},
	},
	{
		"three spans, a intersects with b and c, b with c",
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 16, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
	},
	{
		"three spans, a engulfs b and intersects with c",
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 7, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 0, 5, 0, berlin),
				time.Date(2020, 9, 26, 16, 13, 5, 0, berlin),
			),
		},
		Spans{
			New(
				time.Date(2020, 9, 26, 11, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 7, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 0, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 0, 5, 0, berlin),
				time.Date(2020, 9, 26, 14, 7, 5, 0, berlin),
			),
		},
	},
}

func TestIntersections(t *testing.T) {
	for _, tt := range intersectionTests {
		t.Log(tt.description)
		if !reflect.DeepEqual(tt.intersections, tt.spans.Intersection()) {
			t.Fail()
		}
	}
}

var withoutTests = []struct {
	description string
	spans       Spans
	exclude     Span
	expected    Spans
}{
	{
		"no spans",
		Spans{},
		New(
			time.Date(2020, 9, 26, 13, 6, 5, 0, berlin),
			time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
		),
		Spans{},
	},
	{
		"span interrupted by an start=end span",
		Spans{
			New(
				time.Date(2020, 9, 26, 13, 6, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 14, 6, 5, 0, berlin),
			time.Date(2020, 9, 26, 14, 6, 5, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 13, 6, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
	},
	{
		"one span overlapping end",
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 13, 6, 5, 0, berlin),
			time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 13, 6, 5, 0, berlin),
			),
		},
	},
	{
		"one span overlapping start",
		Spans{
			New(
				time.Date(2020, 9, 26, 10, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 8, 6, 5, 0, berlin),
			time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
	},
	{
		"one span no overlapping span = a",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 13, 6, 5, 0, berlin),
			time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
	},
	{
		"one span complete overlap = empty spans",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 8, 6, 5, 0, berlin),
			time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
		),
		Spans{},
	},
	{
		"one span identical overlap = empty span",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
			time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
		),
		Spans{},
	},
	{
		"one span punch hole",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 10, 6, 5, 0, berlin),
			time.Date(2020, 9, 26, 11, 13, 5, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 10, 6, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 11, 13, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
	},
	{
		"one span contiguous span = a",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
		},
	},
	{
		"two spans minus one overlapping both = two spans",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 11, 13, 5, 0, berlin),
			time.Date(2020, 9, 26, 14, 30, 5, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 11, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 30, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
	},
	{
		"two spans minus one cutting the first in two = three spans",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 10, 0, 38, 0, berlin),
			time.Date(2020, 9, 26, 11, 30, 0, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 10, 0, 38, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 11, 30, 0, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
		},
	},
	{
		"four spans minus one wrecking havoc",
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 12, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 8, 0, 7, 0, berlin),
				time.Date(2020, 9, 26, 13, 33, 38, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 1, 27, 0, berlin),
				time.Date(2020, 9, 26, 14, 30, 38, 0, berlin),
			),
		},
		New(
			time.Date(2020, 9, 26, 10, 0, 38, 0, berlin),
			time.Date(2020, 9, 26, 13, 27, 0, 0, berlin),
		),
		Spans{
			New(
				time.Date(2020, 9, 26, 9, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 10, 0, 38, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 4, 5, 0, berlin),
				time.Date(2020, 9, 26, 15, 13, 5, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 8, 0, 7, 0, berlin),
				time.Date(2020, 9, 26, 10, 0, 38, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 13, 27, 0, 0, berlin),
				time.Date(2020, 9, 26, 13, 33, 38, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 14, 1, 27, 0, berlin),
				time.Date(2020, 9, 26, 14, 30, 38, 0, berlin),
			),
		},
	},
}

func TestWithout(t *testing.T) {
	for _, tt := range withoutTests {
		t.Log(tt.description)
		result := tt.spans.Without(tt.exclude)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Error("Expected ", tt.expected, "\nReceived ", result)
			t.Fail()
		}
	}
}

var withinTests = []struct {
	description string
	spanA       Span
	spanB       Span
	expected    bool
}{
	{
		"is within",
		New(
			time.Date(2020, 9, 26, 11, 27, 0, 0, berlin),
			time.Date(2020, 9, 26, 15, 33, 38, 0, berlin),
		),
		New(
			time.Date(2020, 9, 26, 12, 27, 0, 0, berlin),
			time.Date(2020, 9, 26, 14, 33, 38, 0, berlin),
		),
		true,
	},
	{
		"is identical",
		New(
			time.Date(2020, 9, 26, 11, 27, 0, 0, berlin),
			time.Date(2020, 9, 26, 15, 33, 38, 0, berlin),
		),
		New(
			time.Date(2020, 9, 26, 11, 27, 0, 0, berlin),
			time.Date(2020, 9, 26, 15, 33, 38, 0, berlin),
		),
		true,
	},
	{
		"start outside",
		New(
			time.Date(2020, 9, 26, 11, 27, 0, 0, berlin),
			time.Date(2020, 9, 26, 15, 33, 38, 0, berlin),
		),
		New(
			time.Date(2020, 9, 26, 11, 26, 0, 0, berlin),
			time.Date(2020, 9, 26, 13, 33, 38, 0, berlin),
		),
		false,
	},
	{
		"end outside",
		New(
			time.Date(2020, 9, 26, 11, 27, 0, 0, berlin),
			time.Date(2020, 9, 26, 15, 33, 38, 0, berlin),
		),
		New(
			time.Date(2020, 9, 26, 12, 27, 0, 0, berlin),
			time.Date(2020, 9, 26, 15, 34, 38, 0, berlin),
		),
		false,
	},
}

func TestWithin(t *testing.T) {
	for _, tt := range withinTests {
		t.Log(tt.description)
		if Within(tt.spanA, tt.spanB) != tt.expected {
			t.Error("Expected ", tt.expected, ", got ", !tt.expected)
		}
	}
}

var testDurations = []struct {
	description string
	spans       Spans
	expected    time.Duration
}{
	{
		"single Span",
		Spans{
			New(
				time.Date(2020, 9, 26, 11, 27, 0, 0, berlin),
				time.Date(2020, 9, 26, 13, 12, 45, 0, berlin),
			),
		},
		6345000000000,
	},
	{
		"two Spans overlapping ",
		Spans{
			New(
				time.Date(2020, 9, 26, 11, 27, 0, 0, berlin),
				time.Date(2020, 9, 26, 13, 12, 45, 0, berlin),
			),
			New(
				time.Date(2020, 9, 26, 12, 12, 36, 0, berlin),
				time.Date(2020, 9, 26, 14, 23, 56, 0, berlin),
			),
		},
		14225000000000,
	},
}

func TestSpans_Duration(t *testing.T) {
	for _, tt := range testDurations {
		t.Log(tt.description)
		duration := tt.spans.Duration()
		if tt.expected == duration {
			t.Error("Expected ", tt.expected, " Received ", duration)
			t.Fail()
		}
	}
}
