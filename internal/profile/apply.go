package profile

import (
	"fmt"

	"github.com/SciTee/ledgerly/internal/model"
	"github.com/SciTee/ledgerly/internal/parse"
)

// readOptions builds the CSV reader options implied by the profile.
func (p *Profile) readOptions() parse.ReadOptions {
	opts := parse.ReadOptions{SkipRows: p.SkipRows}
	if d, ok := p.delimiterRune(); ok {
		opts.Delimiter = d
	}
	return opts
}

// Apply maps every data row of t onto a unified Transaction. The account label
// is stamped on each transaction so merged timelines can tell sources apart.
func (p *Profile) Apply(t *parse.Table, account string) ([]model.Transaction, error) {
	idx := buildIndex(t.Header)
	if err := p.checkColumns(idx); err != nil {
		return nil, err
	}

	txs := make([]model.Transaction, 0, len(t.Rows))
	for i, row := range t.Rows {
		if isBlank(row) {
			continue
		}
		tx, err := p.mapRow(row, idx, account)
		if err != nil {
			return nil, fmt.Errorf("%s row %d: %w", p.Name, i+1, err)
		}
		txs = append(txs, tx)
	}
	return txs, nil
}
