package spaniel

import (
	"testing"
	"time"
)

func TestOverlap(t *testing.T) {
	begin, _ := time.Parse("2006-01-02T15:04:05Z", "2020-09-26T02:00:00Z")
	end1, _ := time.Parse("2006-01-02T15:04:05Z", "2020-09-26T06:00:00Z")
	end2, _ := time.Parse("2006-01-02T15:04:05Z", "2020-09-26T06:15:00Z")
	span1 := New(begin, end1)
	span2 := New(begin, end2)
	shouldOverlap := overlap(span1, span2)
	if !shouldOverlap {
		t.Fail()
	}
	shouldOverlap = overlap(span2, span1)
	if !shouldOverlap {
		t.Fail()
	}
	beginNotOverlap1, _ := time.Parse("2006-01-02T15:04:05Z", "2020-09-26T02:00:00Z")
	endNotOverlap1, _ := time.Parse("2006-01-02T15:04:05Z", "2020-09-26T06:00:00Z")
	beginNotOverlap2, _ := time.Parse("2006-01-02T15:04:05Z", "2020-09-26T06:01:00Z")
	endNotOverlap2, _ := time.Parse("2006-01-02T15:04:05Z", "2020-09-26T20:00:00Z")
	spanNotOverlap1 := New(beginNotOverlap1, endNotOverlap1)
	spanNotOverlap2 := New(beginNotOverlap2, endNotOverlap2)

	shouldNotOverlap := overlap(spanNotOverlap1, spanNotOverlap2)
	if shouldNotOverlap {
		t.Fail()
	}
	shouldNotOverlap = overlap(spanNotOverlap2, spanNotOverlap1)
	if shouldNotOverlap {
		t.Fail()
	}

}
