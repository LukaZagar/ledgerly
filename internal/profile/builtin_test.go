package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/LukaZagar/ledgerly/internal/parse"
)

// TestBuiltinProfilesGolden applies each built-in profile to its sample export
// and checks the unified result against a golden file. Run with -update to
// regenerate the golden files after an intentional change.
func TestBuiltinProfilesGolden(t *testing.T) {
	cases := []struct {
		profile string
		fixture string
		golden  string
	}{
		{"DKB", "dkb_sample.csv", "dkb.json"},
		{"Sparkasse", "sparkasse_sample.csv", "sparkasse.json"},
		{"PayPal", "paypal_sample.csv", "paypal.json"},
	}

	for _, c := range cases {
		t.Run(c.profile, func(t *testing.T) {
			p, err := LoadBuiltin(c.profile)
			if err != nil {
				t.Fatalf("load builtin %s: %v", c.profile, err)
			}

			f, err := os.Open(filepath.Join("testdata", c.fixture))
			if err != nil {
				t.Fatalf("open fixture: %v", err)
			}
			defer f.Close()

			tab, err := parse.ReadCSVOpts(f, p.ReadOptions())
			if err != nil {
				t.Fatalf("read csv: %v", err)
			}
			txs, err := p.Apply(tab, c.profile)
			if err != nil {
				t.Fatalf("apply: %v", err)
			}

			got, err := json.MarshalIndent(txs, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			got = append(got, '\n')

			goldenPath := filepath.Join("testdata", "golden", c.golden)
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
				t.Errorf("%s mismatch.\n--- got ---\n%s\n--- want ---\n%s", c.profile, got, want)
			}
		})
	}
}
