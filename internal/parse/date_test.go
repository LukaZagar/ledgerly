package parse

import (
	"testing"
	"time"

	"github.com/LukaZagar/ledgerly/internal/model"
)

func TestParseDateDefaults(t *testing.T) {
	want := model.NewDate(2026, time.January, 2)
	tests := []struct {
		name string
		in   string
	}{
		{"iso", "2026-01-02"},
		{"german dotted", "02.01.2026"},
		{"german dotted non-padded", "2.1.2026"},
		{"slash day first", "02/01/2026"},
		{"dash day first", "02-01-2026"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.in, nil)
			if err != nil {
				t.Fatalf("ParseDate(%q): %v", tt.in, err)
			}
			if !got.Equal(want.Time) {
				t.Errorf("ParseDate(%q) = %s, want %s", tt.in, got, want)
			}
		})
	}
}

func TestParseDateCustomLayout(t *testing.T) {
	got, err := ParseDate("2026/12/24", []string{"2006/01/02"})
	if err != nil {
		t.Fatalf("ParseDate: %v", err)
	}
	if want := model.NewDate(2026, time.December, 24); !got.Equal(want.Time) {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestParseDateErrors(t *testing.T) {
	for _, in := range []string{"", "not-a-date", "31.31.2026"} {
		if _, err := ParseDate(in, nil); err == nil {
			t.Errorf("ParseDate(%q) expected error, got nil", in)
		}
	}
}
