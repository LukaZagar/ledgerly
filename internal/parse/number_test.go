package parse

import (
	"testing"

	"github.com/SciTee/ledgerly/internal/model"
)

func TestParseAmountAuto(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want model.Money
	}{
		{"german grouping+decimal", "1.234,56", 123456},
		{"us grouping+decimal", "1,234.56", 123456},
		{"plain dot decimal", "1234.56", 123456},
		{"plain comma decimal", "1234,56", 123456},
		{"negative comma", "-12,50", -1250},
		{"explicit plus", "+99,99", 9999},
		{"leading zero cents", "0,07", 7},
		{"one fraction digit", "1,5", 150},
		{"currency symbol and space", "€ 1.234,56", 123456},
		{"german millions", "1.234.567,89", 123456789},
		{"us millions", "1,234,567.89", 123456789},
		{"comma as grouping only", "1,234", 123400},
		{"dot as grouping only", "1.234", 123400},
		{"integer", "42", 4200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAmount(tt.in, DecimalAuto)
			if err != nil {
				t.Fatalf("ParseAmount(%q) error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseAmount(%q) = %d, want %d", tt.in, int64(got), int64(tt.want))
			}
		})
	}
}

func TestParseAmountExplicitSeparator(t *testing.T) {
	tests := []struct {
		name string
		in   string
		dec  DecimalSeparator
		want model.Money
	}{
		{"force comma decimal", "1.234,56", DecimalComma, 123456},
		{"force dot decimal", "1,234.56", DecimalDot, 123456},
		// With comma forced as the decimal separator, "1,234" reads as 1.234
		// and rounds down to 1.23 (123 cents).
		{"force comma on ambiguous", "1,234", DecimalComma, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAmount(tt.in, tt.dec)
			if err != nil {
				t.Fatalf("ParseAmount(%q) error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("ParseAmount(%q, %v) = %d, want %d", tt.in, tt.dec, int64(got), int64(tt.want))
			}
		})
	}
}

func TestParseAmountRounding(t *testing.T) {
	got, err := ParseAmount("1.005", DecimalDot)
	if err != nil {
		t.Fatal(err)
	}
	if got != 101 {
		t.Errorf("ParseAmount(1.005) = %d, want 101 (rounded)", int64(got))
	}
}

func TestParseAmountErrors(t *testing.T) {
	for _, in := range []string{"", "   ", "abc", "12.x4"} {
		if _, err := ParseAmount(in, DecimalAuto); err == nil {
			t.Errorf("ParseAmount(%q) expected error, got nil", in)
		}
	}
}
