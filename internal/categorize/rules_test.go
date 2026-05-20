package categorize

import (
	"strings"
	"testing"

	"github.com/SciTee/ledgerly/internal/model"
)

func tx(payee, purpose string) model.Transaction {
	return model.Transaction{Payee: payee, Purpose: purpose}
}

func TestDefaultRulesCompile(t *testing.T) {
	rs, err := Default()
	if err != nil {
		t.Fatalf("default rules failed to compile: %v", err)
	}
	if rs.Len() == 0 {
		t.Fatal("default rule set is empty")
	}
}

func TestDefaultCategorization(t *testing.T) {
	rs, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		tx   model.Transaction
		want string
	}{
		{"groceries", tx("REWE SAGT DANKE", "Einkauf"), "Groceries"},
		{"subscription netflix", tx("NETFLIX.COM", "Abo"), "Subscriptions"},
		{"amazon prime beats shopping", tx("AMAZON PRIME VIDEO", "Abo"), "Subscriptions"},
		{"bare amazon is shopping", tx("AMAZON MKTPLACE", "Bestellung"), "Shopping"},
		{"income", tx("Mustermann GmbH", "Gehalt Januar"), "Income"},
		{"unknown falls back", tx("Tante Emma Laden", "Sonstiges"), model.Uncategorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rs.Categorize(tt.tx); got != tt.want {
				t.Errorf("Categorize(%q/%q) = %q, want %q", tt.tx.Payee, tt.tx.Purpose, got, tt.want)
			}
		})
	}
}

func TestFirstMatchWins(t *testing.T) {
	rs, err := Compile([]Rule{
		{Match: "REWE", Category: "First"},
		{Match: "REWE", Category: "Second"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := rs.Categorize(tx("REWE City", "")); got != "First" {
		t.Errorf("got %q, want First (order should decide)", got)
	}
}

func TestUserRulesOverrideDefaults(t *testing.T) {
	userYAML := "rules:\n  - { match: \"REWE\", category: \"Food\" }\n"
	rs, err := LoadWithDefaults(strings.NewReader(userYAML))
	if err != nil {
		t.Fatal(err)
	}
	// User rule wins over the default Groceries mapping...
	if got := rs.Categorize(tx("REWE SAGT DANKE", "")); got != "Food" {
		t.Errorf("user override = %q, want Food", got)
	}
	// ...but defaults still apply to everything else.
	if got := rs.Categorize(tx("NETFLIX.COM", "")); got != "Subscriptions" {
		t.Errorf("default fallthrough = %q, want Subscriptions", got)
	}
}

func TestFieldScoping(t *testing.T) {
	rs, err := Compile([]Rule{
		{Match: "Gehalt", Field: FieldPurpose, Category: "Income"},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Match only against purpose: a payee literally named "Gehalt GmbH" with an
	// unrelated purpose must not match.
	if got := rs.Categorize(tx("Gehalt GmbH", "Rechnung")); got != model.Uncategorized {
		t.Errorf("payee-only text matched a purpose-scoped rule: %q", got)
	}
	if got := rs.Categorize(tx("Arbeitgeber", "Gehalt Mai")); got != "Income" {
		t.Errorf("purpose match = %q, want Income", got)
	}
}

func TestApplySetsCategory(t *testing.T) {
	rs, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	txs := []model.Transaction{tx("REWE", "x"), tx("Nobody", "nothing")}
	rs.Apply(txs)
	if txs[0].Category != "Groceries" {
		t.Errorf("tx[0] category = %q", txs[0].Category)
	}
	if txs[1].Category != model.Uncategorized {
		t.Errorf("tx[1] category = %q", txs[1].Category)
	}
}
