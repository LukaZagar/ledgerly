package profile

import (
	"fmt"
	"strings"

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

// mapRow turns a single CSV row into a Transaction.
func (p *Profile) mapRow(row []string, idx columnIndex, account string) (model.Transaction, error) {
	dateStr := get(row, idx, p.Columns.Date)
	date, err := parse.ParseDate(dateStr, p.DateLayouts)
	if err != nil {
		return model.Transaction{}, err
	}

	tx := model.Transaction{
		Date:      date,
		Currency:  p.Currency,
		Payee:     get(row, idx, p.Columns.Payee),
		Purpose:   get(row, idx, p.Columns.Purpose),
		Reference: get(row, idx, p.Columns.Reference),
		Account:   account,
	}

	if vd := get(row, idx, p.Columns.ValueDate); vd != "" {
		if d, err := parse.ParseDate(vd, p.DateLayouts); err == nil {
			tx.ValueDate = d
		}
	}

	amount, err := parse.ParseAmount(get(row, idx, p.Columns.Amount), p.decimalSeparator())
	if err != nil {
		return model.Transaction{}, err
	}
	tx.Amount = amount

	return tx, nil
}

// columnIndex maps a header column name to its position in a row.
type columnIndex map[string]int

func buildIndex(header []string) columnIndex {
	idx := make(columnIndex, len(header))
	for i, h := range header {
		idx[h] = i
	}
	return idx
}

// checkColumns verifies the columns the profile references actually exist in
// the file's header.
func (p *Profile) checkColumns(idx columnIndex) error {
	required := []string{p.Columns.Date}
	if p.Columns.Amount != "" {
		required = append(required, p.Columns.Amount)
	}
	for _, col := range required {
		if _, ok := idx[col]; !ok {
			return fmt.Errorf("profile %q: column %q not found in header", p.Name, col)
		}
	}
	return nil
}

// get returns the trimmed value of the named column for a row, or "" when the
// column is unmapped or missing.
func get(row []string, idx columnIndex, col string) string {
	if col == "" {
		return ""
	}
	if i, ok := idx[col]; ok && i < len(row) {
		return strings.TrimSpace(row[i])
	}
	return ""
}

func isBlank(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}
