package export

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SciTee/ledgerly/internal/model"
)

var update = flag.Bool("update", false, "update golden files")

func sampleTransactions() []model.Transaction {
	return []model.Transaction{
		{
			Date:     model.NewDate(2026, time.January, 2),
			Amount:   -4312,
			Currency: "EUR",
			Payee:    "REWE SAGT DANKE",
			Purpose:  "Einkauf Lebensmittel",
			Account:  "ING-Giro",
			Category: "Groceries",
		},
		{
			Date:     model.NewDate(2026, time.January, 5),
			Amount:   245000,
			Currency: "EUR",
			Payee:    "Mustermann GmbH",
			Purpose:  "Gehalt Januar",
			Account:  "ING-Giro",
			Category: "Income",
		},
		{
			Date:      model.NewDate(2026, time.January, 7),
			ValueDate: model.NewDate(2026, time.January, 8),
			Amount:    -1299,
			Currency:  "EUR",
			Payee:     "NETFLIX.COM",
			Purpose:   "Abo Januar",
			Account:   "ING-Giro",
			Category:  "Subscriptions",
			Reference: "NETFLIX-2026-01",
		},
		{
			Date:     model.NewDate(2026, time.January, 9),
			Amount:   -899,
			Currency: "EUR",
			Payee:    "Unbekannt",
			Purpose:  "",
			Account:  "ING-Giro",
			// no category -> should export as Uncategorized
		},
	}
}

func checkGolden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", "golden", name)
	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s (run with -update first): %v", name, err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s mismatch.\n--- got ---\n%s\n--- want ---\n%s", name, got, want)
	}
}

func TestWriteCSVGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteCSV(&buf, sampleTransactions()); err != nil {
		t.Fatal(err)
	}
	checkGolden(t, "unified.csv", buf.Bytes())
}

func TestWriteJSONGolden(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, sampleTransactions()); err != nil {
		t.Fatal(err)
	}
	checkGolden(t, "unified.json", buf.Bytes())
}
