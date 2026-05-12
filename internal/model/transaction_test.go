package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMoneyString(t *testing.T) {
	tests := []struct {
		name string
		in   Money
		want string
	}{
		{"zero", 0, "0.00"},
		{"whole euro", 100, "1.00"},
		{"with cents", 123456, "1234.56"},
		{"single minor digit", 105, "1.05"},
		{"negative", -4200, "-42.00"},
		{"negative cents", -7, "-0.07"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.String(); got != tt.want {
				t.Errorf("Money(%d).String() = %q, want %q", int64(tt.in), got, tt.want)
			}
		})
	}
}

func TestMoneyAbs(t *testing.T) {
	tests := []struct {
		in   Money
		want Money
	}{
		{0, 0},
		{500, 500},
		{-500, 500},
	}
	for _, tt := range tests {
		if got := tt.in.Abs(); got != tt.want {
			t.Errorf("Money(%d).Abs() = %d, want %d", int64(tt.in), int64(got), int64(tt.want))
		}
	}
}

func TestTransactionDirection(t *testing.T) {
	tests := []struct {
		name      string
		amount    Money
		isExpense bool
		isIncome  bool
	}{
		{"outflow", -1500, true, false},
		{"inflow", 1500, false, true},
		{"zero is neither", 0, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := Transaction{Amount: tt.amount}
			if tx.IsExpense() != tt.isExpense {
				t.Errorf("IsExpense() = %v, want %v", tx.IsExpense(), tt.isExpense)
			}
			if tx.IsIncome() != tt.isIncome {
				t.Errorf("IsIncome() = %v, want %v", tx.IsIncome(), tt.isIncome)
			}
		})
	}
}

func TestCategoryOrDefault(t *testing.T) {
	if got := CategoryOrDefault(""); got != Uncategorized {
		t.Errorf("empty category = %q, want %q", got, Uncategorized)
	}
	if got := CategoryOrDefault(CategoryGroceries); got != CategoryGroceries {
		t.Errorf("non-empty category = %q, want %q", got, CategoryGroceries)
	}
}

func TestDateString(t *testing.T) {
	if got := NewDate(2026, time.May, 9).String(); got != "2026-05-09" {
		t.Errorf("Date.String() = %q, want 2026-05-09", got)
	}
	var zero Date
	if got := zero.String(); got != "" {
		t.Errorf("zero Date.String() = %q, want empty", got)
	}
}

func TestDateJSONRoundTrip(t *testing.T) {
	orig := NewDate(2026, time.January, 31)
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(b) != `"2026-01-31"` {
		t.Fatalf("marshal = %s, want \"2026-01-31\"", b)
	}
	var back Date
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !back.Equal(orig.Time) {
		t.Errorf("round-trip = %v, want %v", back, orig)
	}

	var null Date
	if err := json.Unmarshal([]byte("null"), &null); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if !null.IsZero() {
		t.Errorf("null did not decode to zero date")
	}
}
