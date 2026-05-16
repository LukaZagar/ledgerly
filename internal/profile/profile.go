// Package profile describes how a particular bank's CSV columns map onto the
// unified Transaction schema, and applies that mapping to parsed rows.
package profile

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/SciTee/ledgerly/internal/parse"
	"gopkg.in/yaml.v3"
)

// Profile is a single bank's CSV layout. It is deliberately plain YAML so that
// users can drop their own profile next to the built-in ones without touching
// Go code.
type Profile struct {
	Name string `yaml:"name"`

	// Delimiter is an optional single-character field separator. When empty the
	// reader auto-detects it.
	Delimiter string `yaml:"delimiter"`

	// DateLayouts are Go time layouts tried in order. Empty falls back to the
	// parser defaults.
	DateLayouts []string `yaml:"date_layouts"`

	// Decimal selects how amounts are read: "comma" (1.234,56), "dot"
	// (1,234.56) or "auto".
	Decimal string `yaml:"decimal"`

	// Currency is the ISO 4217 code applied to every row, e.g. "EUR".
	Currency string `yaml:"currency"`

	// SkipRows drops a number of preamble lines some banks print above the
	// real header.
	SkipRows int `yaml:"skip_rows"`

	// Signature lists header column names that uniquely identify this bank,
	// used by auto-detection.
	Signature []string `yaml:"signature"`

	// Columns maps unified fields to this bank's column names.
	Columns Columns `yaml:"columns"`
}

// Columns maps each unified field to the source column header. An amount can
// arrive either as a single signed column (Amount) or as a Debit/Credit pair.
type Columns struct {
	Date      string `yaml:"date"`
	ValueDate string `yaml:"value_date"`
	Amount    string `yaml:"amount"`
	Debit     string `yaml:"debit"`
	Credit    string `yaml:"credit"`
	Payee     string `yaml:"payee"`
	Purpose   string `yaml:"purpose"`
	Reference string `yaml:"reference"`
}

// Load decodes a profile from YAML and validates it.
func Load(r io.Reader) (*Profile, error) {
	var p Profile
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

// LoadFile reads and decodes a profile from a YAML file on disk.
func LoadFile(path string) (*Profile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Load(f)
}

// Validate checks that the profile carries the minimum information needed to
// turn rows into transactions.
func (p *Profile) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("profile: name is required")
	}
	if p.Columns.Date == "" {
		return fmt.Errorf("profile %q: a date column is required", p.Name)
	}
	hasSigned := p.Columns.Amount != ""
	hasSplit := p.Columns.Debit != "" || p.Columns.Credit != ""
	if !hasSigned && !hasSplit {
		return fmt.Errorf("profile %q: need either an amount column or debit/credit columns", p.Name)
	}
	return nil
}

// decimalSeparator resolves the configured decimal mode.
func (p *Profile) decimalSeparator() parse.DecimalSeparator {
	switch strings.ToLower(strings.TrimSpace(p.Decimal)) {
	case "comma":
		return parse.DecimalComma
	case "dot":
		return parse.DecimalDot
	default:
		return parse.DecimalAuto
	}
}

// delimiterRune returns the configured delimiter and whether one was set.
func (p *Profile) delimiterRune() (rune, bool) {
	d := strings.TrimSpace(p.Delimiter)
	if d == "" {
		return 0, false
	}
	return []rune(d)[0], true
}
