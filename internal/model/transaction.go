// Package model defines the unified transaction schema that every bank CSV is
// normalized into. Keeping it in one place means parsers, exporters and the
// categorizer all agree on the same shape.
package model

import "fmt"

// Money is a monetary amount stored in minor units (e.g. cents for EUR) to
// avoid floating-point rounding when summing many transactions. Negative
// values are outflows, positive values are inflows.
type Money int64

// String renders the amount with a dot decimal separator and two fraction
// digits, e.g. -1234 -> "-12.34". This is the neutral form used in exports;
// locale-specific formatting happens elsewhere.
func (m Money) String() string {
	neg := m < 0
	v := int64(m)
	if neg {
		v = -v
	}
	s := fmt.Sprintf("%d.%02d", v/100, v%100)
	if neg {
		s = "-" + s
	}
	return s
}

// Abs returns the absolute value of the amount.
func (m Money) Abs() Money {
	if m < 0 {
		return -m
	}
	return m
}

// Transaction is one normalized booking line. Zero values are valid for the
// optional fields (ValueDate, Reference) so partially-populated bank exports
// still map cleanly.
type Transaction struct {
	Date      Date   `json:"date"`                 // booking date
	ValueDate Date   `json:"value_date,omitempty"` // Wertstellung, optional
	Amount    Money  `json:"amount"`               // signed, minor units
	Currency  string `json:"currency"`             // ISO 4217, e.g. "EUR"
	Payee     string `json:"payee"`                // counterparty
	Purpose   string `json:"purpose"`              // description / Verwendungszweck
	Account   string `json:"account"`              // source account label
	Category  string `json:"category"`             // assigned category
	Reference string `json:"reference,omitempty"`  // bank/end-to-end reference
}

// IsExpense reports whether the transaction moves money out of the account.
func (t Transaction) IsExpense() bool { return t.Amount < 0 }

// IsIncome reports whether the transaction moves money into the account.
func (t Transaction) IsIncome() bool { return t.Amount > 0 }
