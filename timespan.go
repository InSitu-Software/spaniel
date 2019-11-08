package spaniel

import (
	"encoding/json"
	"fmt"
	"time"
)

// TimeSpan represents a simple span of time, with no additional properties. It should be constructed with NewEmpty.
type TimeSpan struct {
	start time.Time
	end   time.Time
}

// Start returns the start time of a span
func (ts TimeSpan) Start() time.Time { return ts.start }

// End returns the end time of a span
func (ts TimeSpan) End() time.Time { return ts.end }

func (ts TimeSpan) Duration() time.Duration {
	return ts.end.Sub(ts.start)
}

// MarshalJSON implements json.Marshal
func (ts TimeSpan) MarshalJSON() ([]byte, error) {
	o := struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}{
		Start: ts.start,
		End:   ts.end,
	}

	return json.Marshal(o)
}

// UnmarshalJSON implements json.Unmarshal
func (ts *TimeSpan) UnmarshalJSON(b []byte) (err error) {
	var i struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	err = json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	ts.start = i.Start
	ts.end = i.End

	return
}

func (ts *TimeSpan) MarshalBinary() ([]byte, error) {
	startBin, _ := ts.start.MarshalBinary()
	endBin, _ := ts.end.MarshalBinary()

	return append(startBin, endBin...), nil
}

func (ts *TimeSpan) UnmarshalBinary(data []byte) error {
	if len(data) != 30 {
		return fmt.Errorf("unknown binary format, len %d", len(data))
	}

	startBin := data[0:15]
	endBin := data[15:]

	if err := ts.start.UnmarshalBinary(startBin); err != nil {
		return err
	}

	if err := ts.end.UnmarshalBinary(endBin); err != nil {
		return err
	}

	return nil
}

func (ts TimeSpan) String() string {
	return fmt.Sprintf(
		"%s - %s",
		ts.Start().Format("2006-01-02 15:04"),
		ts.End().Format("2006-01-02 15:04"),
	)
}

// New creates a span with a start and end time, with the types set to [] for instants and [) for spans.
func New(start time.Time, end time.Time) *TimeSpan {
	return &TimeSpan{
		start: start,
		end:   end,
	}
}
