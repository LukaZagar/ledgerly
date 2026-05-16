package parse

import (
	"fmt"
	"strings"
	"time"

	"github.com/SciTee/ledgerly/internal/model"
)

// DefaultDateLayouts are tried, in order, when a profile does not pin down its
// own date format. ISO comes first because it is unambiguous; the German
// day-first formats follow. Padded and non-padded variants are both listed so
// "2.1.2026" and "02.01.2026" both parse.
var DefaultDateLayouts = []string{
	"2006-01-02",
	"02.01.2006",
	"2.1.2006",
	"02/01/2006",
	"2/1/2006",
	"2006/01/02",
	"02-01-2006",
}

// ParseDate parses a date string using the supplied layouts, falling back to
// DefaultDateLayouts when none are given. The first layout that matches wins.
func ParseDate(s string, layouts []string) (model.Date, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return model.Date{}, fmt.Errorf("parse date: empty")
	}
	if len(layouts) == 0 {
		layouts = DefaultDateLayouts
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return model.Date{Time: t}, nil
		}
	}
	return model.Date{}, fmt.Errorf("parse date %q: no matching layout", s)
}
