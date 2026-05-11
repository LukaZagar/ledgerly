package model

import (
	"strings"
	"time"
)

// dateLayout is the canonical wire format for dates in exports and JSON.
const dateLayout = "2006-01-02"

// Date is a calendar date without a time-of-day or timezone component. Bank
// statements deal in days, not instants, so dragging a full time.Time (with
// its zone) through the pipeline only invites off-by-one bugs around midnight.
type Date struct {
	time.Time
}

// NewDate builds a Date from year, month and day.
func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// String renders the date as YYYY-MM-DD, or the empty string if unset.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format(dateLayout)
}

// MarshalJSON emits the date as a "YYYY-MM-DD" string, or null when unset.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Format(dateLayout) + `"`), nil
}

// UnmarshalJSON parses a "YYYY-MM-DD" string (or null) back into a Date.
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}
