// Package summary computes the per-category and per-month totals that ledgerly
// prints to the terminal after a conversion.
package summary

import (
	"fmt"
	"io"
	"sort"

	"github.com/LukaZagar/ledgerly/internal/model"
)

// Bucket is a labelled total with a transaction count.
type Bucket struct {
	Label string
	Total model.Money
	Count int
}

// Summary is the aggregated view of a set of transactions.
type Summary struct {
	Count      int
	Income     model.Money
	Expense    model.Money // negative
	Net        model.Money
	From, To   model.Date
	Categories []Bucket // sorted by absolute total, descending
	Months     []Bucket // sorted by month, ascending
}

// Compute aggregates the transactions into a Summary.
func Compute(txs []model.Transaction) Summary {
	var s Summary
	s.Count = len(txs)

	catTotals := map[string]*Bucket{}
	monthTotals := map[string]*Bucket{}

	for _, tx := range txs {
		if tx.Amount >= 0 {
			s.Income += tx.Amount
		} else {
			s.Expense += tx.Amount
		}
		s.Net += tx.Amount

		if s.From.IsZero() || tx.Date.Before(s.From.Time) {
			s.From = tx.Date
		}
		if s.To.IsZero() || tx.Date.After(s.To.Time) {
			s.To = tx.Date
		}

		add(catTotals, model.CategoryOrDefault(tx.Category), tx.Amount)
		if !tx.Date.IsZero() {
			add(monthTotals, tx.Date.Format("2006-01"), tx.Amount)
		}
	}

	s.Categories = sortedByMagnitude(catTotals)
	s.Months = sortedByLabel(monthTotals)
	return s
}

func add(m map[string]*Bucket, label string, amount model.Money) {
	b := m[label]
	if b == nil {
		b = &Bucket{Label: label}
		m[label] = b
	}
	b.Total += amount
	b.Count++
}

func sortedByMagnitude(m map[string]*Bucket) []Bucket {
	out := flatten(m)
	sort.Slice(out, func(i, j int) bool {
		ai, aj := out[i].Total.Abs(), out[j].Total.Abs()
		if ai != aj {
			return ai > aj
		}
		return out[i].Label < out[j].Label
	})
	return out
}

func sortedByLabel(m map[string]*Bucket) []Bucket {
	out := flatten(m)
	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	return out
}

func flatten(m map[string]*Bucket) []Bucket {
	out := make([]Bucket, 0, len(m))
	for _, b := range m {
		out = append(out, *b)
	}
	return out
}

// WriteText renders the summary as a compact, aligned text block.
func (s Summary) WriteText(w io.Writer) error {
	bw := &errWriter{w: w}

	bw.printf("Transactions: %d\n", s.Count)
	if !s.From.IsZero() {
		bw.printf("Period:       %s to %s\n", s.From, s.To)
	}

	bw.printf("\nBy category:\n")
	for _, b := range s.Categories {
		bw.printf("  %-16s %12s  (%d)\n", b.Label, b.Total.String(), b.Count)
	}

	bw.printf("\nBy month:\n")
	for _, b := range s.Months {
		bw.printf("  %-16s %12s  (%d)\n", b.Label, b.Total.String(), b.Count)
	}

	bw.printf("\nTotals:\n")
	bw.printf("  %-16s %12s\n", "Income", s.Income.String())
	bw.printf("  %-16s %12s\n", "Expense", s.Expense.String())
	bw.printf("  %-16s %12s\n", "Net", s.Net.String())

	return bw.err
}

// errWriter swallows the repetitive error checks of many Fprintf calls.
type errWriter struct {
	w   io.Writer
	err error
}

func (e *errWriter) printf(format string, args ...any) {
	if e.err != nil {
		return
	}
	_, e.err = fmt.Fprintf(e.w, format, args...)
}
