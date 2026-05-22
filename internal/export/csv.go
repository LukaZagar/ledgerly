// Package export writes a unified transaction list out as CSV or JSON.
package export

import (
	"encoding/csv"
	"io"

	"github.com/SciTee/ledgerly/internal/model"
)

// csvHeader is the fixed column order of the unified CSV output.
var csvHeader = []string{
	"date", "value_date", "amount", "currency",
	"payee", "purpose", "account", "category", "reference",
}

// WriteCSV writes the transactions as a comma-separated file with a fixed
// header. Amounts use a dot decimal separator (e.g. -12.34) and dates are
// ISO-8601, so the output is portable into spreadsheets and budgeting tools
// regardless of the source bank's quirks.
func WriteCSV(w io.Writer, txs []model.Transaction) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(csvHeader); err != nil {
		return err
	}
	for _, tx := range txs {
		rec := []string{
			tx.Date.String(),
			tx.ValueDate.String(),
			tx.Amount.String(),
			tx.Currency,
			tx.Payee,
			tx.Purpose,
			tx.Account,
			model.CategoryOrDefault(tx.Category),
			tx.Reference,
		}
		if err := cw.Write(rec); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
