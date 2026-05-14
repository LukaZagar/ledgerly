// Package parse turns the messy textual fields of a bank CSV (amounts, dates,
// raw rows) into the typed values the rest of ledgerly works with.
package parse

import (
	"fmt"
	"strings"

	"github.com/SciTee/ledgerly/internal/model"
)

// DecimalSeparator tells the amount parser how to interpret '.' and ',' in a
// number. Auto guesses from the shape of the string, which is good enough for
// most exports; a bank profile can pin it down when the data is ambiguous.
type DecimalSeparator int

const (
	// DecimalAuto guesses the decimal separator from the input.
	DecimalAuto DecimalSeparator = iota
	// DecimalComma treats ',' as the decimal separator (1.234,56).
	DecimalComma
	// DecimalDot treats '.' as the decimal separator (1,234.56).
	DecimalDot
)

// ParseAmount parses a monetary string into model.Money (minor units). It
// copes with grouping separators, either decimal convention, surrounding
// whitespace, currency symbols and a leading sign.
func ParseAmount(s string, dec DecimalSeparator) (model.Money, error) {
	raw := s
	s = cleanAmount(s)
	if s == "" {
		return 0, fmt.Errorf("parse amount %q: empty", raw)
	}

	neg := false
	switch s[0] {
	case '+':
		s = s[1:]
	case '-':
		neg = true
		s = s[1:]
	}

	intPart, fracPart, err := splitDecimal(s, dec)
	if err != nil {
		return 0, fmt.Errorf("parse amount %q: %w", raw, err)
	}

	cents, err := toCents(intPart, fracPart)
	if err != nil {
		return 0, fmt.Errorf("parse amount %q: %w", raw, err)
	}
	if neg {
		cents = -cents
	}
	return model.Money(cents), nil
}

// cleanAmount strips currency symbols and whitespace (including the
// non-breaking space some banks emit) but keeps digits, separators and signs.
func cleanAmount(s string) string {
	s = strings.TrimSpace(s)
	replacer := strings.NewReplacer(
		" ", "", // non-breaking space
		" ", "",
		"€", "",
		"$", "",
		"£", "",
		"EUR", "",
		"USD", "",
	)
	return replacer.Replace(s)
}

// splitDecimal separates the integer and fractional digits, removing grouping
// separators, according to the chosen convention.
func splitDecimal(s string, dec DecimalSeparator) (intPart, fracPart string, err error) {
	lastComma := strings.LastIndexByte(s, ',')
	lastDot := strings.LastIndexByte(s, '.')

	var decByte byte
	switch dec {
	case DecimalComma:
		decByte = ','
	case DecimalDot:
		decByte = '.'
	default: // DecimalAuto
		decByte = guessDecimalByte(s, lastComma, lastDot)
	}

	var groupByte byte = '.'
	if decByte == '.' {
		groupByte = ','
	}

	s = strings.ReplaceAll(s, string(groupByte), "")
	if i := strings.IndexByte(s, decByte); i >= 0 {
		intPart = s[:i]
		fracPart = s[i+1:]
	} else {
		intPart = s
	}

	if intPart == "" {
		intPart = "0"
	}
	if !isDigits(intPart) || (fracPart != "" && !isDigits(fracPart)) {
		return "", "", fmt.Errorf("unexpected characters in %q", s)
	}
	return intPart, fracPart, nil
}

// guessDecimalByte picks the decimal separator for DecimalAuto.
func guessDecimalByte(s string, lastComma, lastDot int) byte {
	switch {
	case lastComma >= 0 && lastDot >= 0:
		// Both present: the rightmost one is the decimal separator.
		if lastComma > lastDot {
			return ','
		}
		return '.'
	case lastComma >= 0:
		// Only commas. A single comma with 1-2 trailing digits is a decimal
		// separator; anything else is grouping (e.g. "1,234,567").
		if strings.Count(s, ",") == 1 && len(s)-lastComma-1 <= 2 {
			return ','
		}
		return '.' // treat commas as grouping, leave no decimal separator
	case lastDot >= 0:
		if strings.Count(s, ".") == 1 && len(s)-lastDot-1 <= 2 {
			return '.'
		}
		return ','
	default:
		return '.'
	}
}

// toCents converts split integer/fraction digit strings into int64 minor
// units, rounding half away from zero when more than two fraction digits are
// present.
func toCents(intPart, fracPart string) (int64, error) {
	var cents int64
	for i := 0; i < len(intPart); i++ {
		cents = cents*10 + int64(intPart[i]-'0')
	}
	cents *= 100

	switch {
	case len(fracPart) == 0:
		// nothing to add
	case len(fracPart) == 1:
		cents += int64(fracPart[0]-'0') * 10
	default:
		cents += int64(fracPart[0]-'0')*10 + int64(fracPart[1]-'0')
		if len(fracPart) > 2 && fracPart[2] >= '5' {
			cents++ // round the dropped digits
		}
	}
	return cents, nil
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
