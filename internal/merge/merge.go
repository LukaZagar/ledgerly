// Package merge combines transactions from several accounts into a single
// chronological timeline and removes duplicate rows that appear when export
// date ranges overlap.
package merge

import (
	"sort"

	"github.com/LukaZagar/ledgerly/internal/model"
)

// Merge concatenates one or more transaction slices and returns them sorted
// into a single chronological timeline. The input slices are not modified.
//
// Ordering is by booking date, then account, amount and payee, so the result
// is fully deterministic regardless of the order the files were given in.
func Merge(groups ...[]model.Transaction) []model.Transaction {
	var n int
	for _, g := range groups {
		n += len(g)
	}
	out := make([]model.Transaction, 0, n)
	for _, g := range groups {
		out = append(out, g...)
	}

	sort.SliceStable(out, func(i, j int) bool {
		return less(out[i], out[j])
	})
	return out
}

func less(a, b model.Transaction) bool {
	if !a.Date.Equal(b.Date.Time) {
		return a.Date.Before(b.Date.Time)
	}
	if a.Account != b.Account {
		return a.Account < b.Account
	}
	if a.Amount != b.Amount {
		return a.Amount < b.Amount
	}
	return a.Payee < b.Payee
}
