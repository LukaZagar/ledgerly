package profile

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/SciTee/ledgerly/internal/parse"
)

// update regenerates the golden files: go test ./... -update
var update = flag.Bool("update", false, "update golden files")

func TestApplyINGProfile(t *testing.T) {
	p, err := LoadFile(filepath.Join("builtin", "ing.yaml"))
	if err != nil {
		t.Fatalf("load profile: %v", err)
	}

	f, err := os.Open(filepath.Join("testdata", "ing_sample.csv"))
	if err != nil {
		t.Fatalf("open sample: %v", err)
	}
	defer f.Close()

	tab, err := parse.ReadCSVOpts(f, p.ReadOptions())
	if err != nil {
		t.Fatalf("read csv: %v", err)
	}

	txs, err := p.Apply(tab, "ING-Giro")
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(txs) != 3 {
		t.Fatalf("got %d transactions, want 3", len(txs))
	}

	got, err := json.MarshalIndent(txs, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')

	goldenPath := filepath.Join("testdata", "golden", "ing.json")
	if *update {
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden (run with -update first): %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("output mismatch.\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
