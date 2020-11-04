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
}

func TestMoreThanTwoIntersections(t *testing.T) {
	for _, tt := range intersectionTests {
		t.Log(tt.description)
		if !reflect.DeepEqual(tt.intersections, tt.spans.Intersection()) {
			t.Fail()
		}
	}
}
