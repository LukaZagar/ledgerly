package merge

import (
	"testing"
	"time"

	"github.com/SciTee/ledgerly/internal/model"
)

func mk(day int, amount model.Money, payee, purpose, account string) model.Transaction {
	return model.Transaction{
		Date:    model.NewDate(2026, time.January, day),
		Amount:  amount,
		Payee:   payee,
		Purpose: purpose,
		Account: account,
	}
}

func TestMergeChronological(t *testing.T) {
	giro := []model.Transaction{
		mk(10, -500, "REWE", "Einkauf", "giro"),
		mk(3, 245000, "Arbeitgeber", "Gehalt", "giro"),
	}
	credit := []model.Transaction{
		mk(7, -1299, "Netflix", "Abo", "credit"),
	}

	got := Merge(giro, credit)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	wantDays := []int{3, 7, 10}
	for i, d := range wantDays {
		if got[i].Date.Day() != d {
			t.Errorf("position %d: day = %d, want %d", i, got[i].Date.Day(), d)
		}
	}
}

func TestMergeDoesNotMutateInputs(t *testing.T) {
	a := []model.Transaction{mk(5, 100, "B", "", "x"), mk(1, 100, "A", "", "x")}
	_ = Merge(a)
	if a[0].Date.Day() != 5 {
		t.Errorf("input slice was reordered: first day = %d, want 5", a[0].Date.Day())
	}
}

func TestDedupRemovesOverlap(t *testing.T) {
	// Two exports of the same account that overlap on the 7th and 8th.
	exportA := []model.Transaction{
		mk(6, -500, "REWE", "Einkauf", "giro"),
		mk(7, -1299, "Netflix", "Abo Januar", "giro"),
		mk(8, 245000, "Arbeitgeber", "Gehalt", "giro"),
	}
	exportB := []model.Transaction{
		mk(7, -1299, "Netflix", "Abo Januar", "giro"), // duplicate of above
		mk(8, 245000, "Arbeitgeber", "Gehalt", "giro"), // duplicate of above
		mk(9, -4200, "Shell", "Tanken", "giro"),
	}

	timeline := Merge(exportA, exportB)
	if len(timeline) != 6 {
		t.Fatalf("merged len = %d, want 6", len(timeline))
	}

	kept, removed := Dedup(timeline)
	if removed != 2 {
		t.Errorf("removed = %d, want 2", removed)
	}
	if len(kept) != 4 {
		t.Errorf("kept = %d, want 4", len(kept))
	}
}

func TestDedupKeepsDistinctSameDay(t *testing.T) {
	// Two genuinely separate coffees on the same day, same amount, must both
	// survive because the purpose differs... and even identical purposes are a
	// known limitation, so here we differ the payee.
	txs := []model.Transaction{
		mk(4, -350, "Bakery A", "Coffee", "giro"),
		mk(4, -350, "Bakery B", "Coffee", "giro"),
	}
	kept, removed := Dedup(txs)
	if removed != 0 || len(kept) != 2 {
		t.Errorf("distinct rows collapsed: kept=%d removed=%d", len(kept), removed)
	}
}
