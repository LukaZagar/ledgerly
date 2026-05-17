package profile

import (
	"fmt"
	"strings"

	"github.com/SciTee/ledgerly/internal/model"
	"github.com/SciTee/ledgerly/internal/parse"
)

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

	amount, err := p.resolveAmount(row, idx)
	if err != nil {
		return model.Transaction{}, err
	}
	tx.Amount = amount

	return tx, nil
}

// resolveAmount computes the signed amount, either from a single signed column
// or from a debit/credit pair. With a pair, the debit side is treated as an
// outflow (made negative) and the credit side as an inflow; exactly one is
// normally populated per row.
func (p *Profile) resolveAmount(row []string, idx columnIndex) (model.Money, error) {
	dec := p.decimalSeparator()

	if p.Columns.Amount != "" {
		return parse.ParseAmount(get(row, idx, p.Columns.Amount), dec)
	}

	var total model.Money
	if raw := get(row, idx, p.Columns.Debit); raw != "" {
		debit, err := parse.ParseAmount(raw, dec)
		if err != nil {
			return 0, err
		}
		total -= debit.Abs()
	}
	if raw := get(row, idx, p.Columns.Credit); raw != "" {
		credit, err := parse.ParseAmount(raw, dec)
		if err != nil {
			return 0, err
		}
		total += credit.Abs()
	}
	return total, nil
}

// columnIndex maps a header column name to its position in a row.
type columnIndex map[string]int

func buildIndex(header []string) columnIndex {
	idx := make(columnIndex, len(header))
	for i, h := range header {
		// Trim so headers like "Betrag " or a column left dirty by a stray
		// BOM still match the profile's clean column names.
		idx[strings.TrimSpace(h)] = i
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
		if _, ok := idx[strings.TrimSpace(col)]; !ok {
			return fmt.Errorf("profile %q: column %q not found in header", p.Name, col)
		}
	}
	return nil
}

// get returns the trimmed value of the named column for a row, or "" when the
// column is unmapped or missing.
func get(row []string, idx columnIndex, col string) string {
	col = strings.TrimSpace(col)
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
