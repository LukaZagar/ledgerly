package export

import (
	"encoding/json"
	"io"

	"github.com/LukaZagar/ledgerly/internal/model"
)

// jsonTx is the user-facing JSON shape. It deliberately differs from the
// internal model: the amount is rendered as a decimal string (e.g. "-12.34")
// rather than the internal minor-unit integer, and dates are ISO strings. That
// keeps the exported file self-explanatory and free of float rounding.
type jsonTx struct {
	Date      string `json:"date"`
	ValueDate string `json:"value_date,omitempty"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
	Payee     string `json:"payee"`
	Purpose   string `json:"purpose"`
	Account   string `json:"account"`
	Category  string `json:"category"`
	Reference string `json:"reference,omitempty"`
}

// WriteJSON writes the transactions as a pretty-printed JSON array.
func WriteJSON(w io.Writer, txs []model.Transaction) error {
	out := make([]jsonTx, 0, len(txs))
	for _, tx := range txs {
		out = append(out, jsonTx{
			Date:      tx.Date.String(),
			ValueDate: tx.ValueDate.String(),
			Amount:    tx.Amount.String(),
			Currency:  tx.Currency,
			Payee:     tx.Payee,
			Purpose:   tx.Purpose,
			Account:   tx.Account,
			Category:  model.CategoryOrDefault(tx.Category),
			Reference: tx.Reference,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(out)
}
