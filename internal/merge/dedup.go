package merge

import (
	"strconv"
	"strings"

	"github.com/SciTee/ledgerly/internal/model"
)

// Dedup removes duplicate transactions, keeping the first occurrence of each.
// Two rows are considered the same when their booking date, amount, payee and
// purpose match — the signature of a row that shows up twice because two
// exports covered an overlapping period. It returns the deduplicated slice and
// the number of rows dropped.
func Dedup(txs []model.Transaction) (kept []model.Transaction, removed int) {
	seen := make(map[string]struct{}, len(txs))
	kept = make([]model.Transaction, 0, len(txs))
	for _, tx := range txs {
		k := dedupKey(tx)
		if _, ok := seen[k]; ok {
			removed++
			continue
		}
		seen[k] = struct{}{}
		kept = append(kept, tx)
	}
	return kept, removed
}

// dedupKey builds the identity string used to spot duplicates. The account is
// deliberately excluded so the same transaction exported from two overlapping
// statements of one account collapses, while normalisation (lower-case,
// collapsed whitespace) absorbs cosmetic differences between exports.
func dedupKey(tx model.Transaction) string {
	var b strings.Builder
	b.WriteString(tx.Date.String())
	b.WriteByte('|')
	b.WriteString(strconv.FormatInt(int64(tx.Amount), 10))
	b.WriteByte('|')
	b.WriteString(normalize(tx.Payee))
	b.WriteByte('|')
	b.WriteString(normalize(tx.Purpose))
	return b.String()
}

func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}
