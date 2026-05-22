package summary

import (
	"strings"
	"testing"
	"time"

	"github.com/SciTee/ledgerly/internal/model"
)

func txns() []model.Transaction {
	return []model.Transaction{
		{Date: model.NewDate(2026, time.January, 2), Amount: -4312, Category: "Groceries"},
		{Date: model.NewDate(2026, time.January, 5), Amount: 245000, Category: "Income"},
		{Date: model.NewDate(2026, time.February, 7), Amount: -1299, Category: "Subscriptions"},
		{Date: model.NewDate(2026, time.February, 9), Amount: -688, Category: "Groceries"},
	}
}

func TestComputeTotals(t *testing.T) {
	s := Compute(txns())

	if s.Count != 4 {
		t.Errorf("count = %d, want 4", s.Count)
	}
	if s.Income != 245000 {
		t.Errorf("income = %d, want 245000", int64(s.Income))
	}
	if s.Expense != -(4312 + 1299 + 688) {
		t.Errorf("expense = %d", int64(s.Expense))
	}
	if s.Net != 245000-4312-1299-688 {
		t.Errorf("net = %d", int64(s.Net))
	}
	if s.From.Format("2006-01-02") != "2026-01-02" || s.To.Format("2006-01-02") != "2026-02-09" {
		t.Errorf("period = %s..%s", s.From, s.To)
	}
}

func TestComputeBuckets(t *testing.T) {
	s := Compute(txns())

	// Categories are sorted by absolute total: Income (2450) first.
	if len(s.Categories) != 3 || s.Categories[0].Label != "Income" {
		t.Fatalf("unexpected categories: %+v", s.Categories)
	}
	// Groceries appears in both months: -43.12 + -6.88 = -50.00 over 2 rows.
	for _, b := range s.Categories {
		if b.Label == "Groceries" {
			if b.Total != -5000 || b.Count != 2 {
				t.Errorf("groceries bucket = %d/%d, want -5000/2", int64(b.Total), b.Count)
			}
		}
	}
	// Two months, ascending.
	if len(s.Months) != 2 || s.Months[0].Label != "2026-01" || s.Months[1].Label != "2026-02" {
		t.Fatalf("unexpected months: %+v", s.Months)
	}
}

func TestWriteTextContainsSections(t *testing.T) {
	var sb strings.Builder
	if err := Compute(txns()).WriteText(&sb); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	for _, want := range []string{"By category:", "By month:", "Totals:", "Net"} {
		if !strings.Contains(out, want) {
			t.Errorf("summary text missing %q\n%s", want, out)
		}
	}
}
