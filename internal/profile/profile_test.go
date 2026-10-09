package profile

import (
	"strings"
	"testing"

	"github.com/LukaZagar/ledgerly/internal/parse"
)

const messyHeaderCSV = "\xef\xbb\xbf Buchung ; Betrag ;Verwendungszweck\n" +
	"02.01.2026;-9,99;Test\n"

func TestApplyToleratesHeaderWhitespaceAndBOM(t *testing.T) {
	p := &Profile{
		Name:     "Whitespacey",
		Decimal:  "comma",
		Currency: "EUR",
		Columns: Columns{
			Date:    "Buchung",
			Amount:  "Betrag",
			Purpose: "Verwendungszweck",
		},
	}

	tab, err := parse.ReadCSV(strings.NewReader(messyHeaderCSV))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	txs, err := p.Apply(tab, "acc")
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(txs) != 1 {
		t.Fatalf("got %d txs, want 1", len(txs))
	}
	if txs[0].Amount != -999 {
		t.Errorf("amount = %d, want -999", int64(txs[0].Amount))
	}
	if txs[0].Purpose != "Test" {
		t.Errorf("purpose = %q, want Test", txs[0].Purpose)
	}
}

func TestValidateRejectsIncompleteProfile(t *testing.T) {
	cases := map[string]*Profile{
		"no name":   {Columns: Columns{Date: "d", Amount: "a"}},
		"no date":   {Name: "x", Columns: Columns{Amount: "a"}},
		"no amount": {Name: "x", Columns: Columns{Date: "d"}},
	}
	for name, p := range cases {
		if err := p.Validate(); err == nil {
			t.Errorf("%s: expected validation error", name)
		}
	}
}
