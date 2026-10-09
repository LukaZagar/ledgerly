// Package categorize assigns a category to each transaction using an ordered
// list of rules that match on the payee and/or purpose text.
package categorize

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/LukaZagar/ledgerly/internal/model"
	"gopkg.in/yaml.v3"
)

// Field selects which text a rule matches against.
const (
	FieldAny     = "any"     // payee and purpose combined (default)
	FieldPayee   = "payee"   // payee only
	FieldPurpose = "purpose" // purpose only
)

// Rule is one categorization rule. It matches either a case-insensitive
// substring (Match) or a regular expression (Regex) against the chosen field.
// The first rule that matches a transaction, in list order, wins.
type Rule struct {
	Match    string `yaml:"match"`
	Regex    string `yaml:"regex"`
	Field    string `yaml:"field"`
	Category string `yaml:"category"`
}

// RuleSet is an ordered collection of compiled rules.
type RuleSet struct {
	rules []compiledRule
}

type compiledRule struct {
	field    string
	category string
	needle   string         // lower-cased substring, for Match rules
	re       *regexp.Regexp // non-nil for Regex rules
}

// Load decodes a rule set from YAML and compiles it.
func Load(r io.Reader) (*RuleSet, error) {
	var doc struct {
		Rules []Rule `yaml:"rules"`
	}
	if err := yaml.NewDecoder(r).Decode(&doc); err != nil {
		return nil, fmt.Errorf("decode rules: %w", err)
	}
	return Compile(doc.Rules)
}

// Compile turns raw rules into a ready-to-use RuleSet.
func Compile(rules []Rule) (*RuleSet, error) {
	rs := &RuleSet{rules: make([]compiledRule, 0, len(rules))}
	for i, r := range rules {
		if r.Category == "" {
			return nil, fmt.Errorf("rule %d: category is required", i+1)
		}
		cr := compiledRule{field: normalizeField(r.Field), category: r.Category}
		switch {
		case r.Regex != "":
			re, err := regexp.Compile("(?i)" + r.Regex)
			if err != nil {
				return nil, fmt.Errorf("rule %d: bad regex %q: %w", i+1, r.Regex, err)
			}
			cr.re = re
		case r.Match != "":
			cr.needle = strings.ToLower(r.Match)
		default:
			return nil, fmt.Errorf("rule %d: needs either match or regex", i+1)
		}
		rs.rules = append(rs.rules, cr)
	}
	return rs, nil
}

// Categorize returns the category for a single transaction, or
// model.Uncategorized when no rule matches.
func (rs *RuleSet) Categorize(tx model.Transaction) string {
	for _, r := range rs.rules {
		if r.matches(tx) {
			return r.category
		}
	}
	return model.Uncategorized
}

// Apply categorizes every transaction in place.
func (rs *RuleSet) Apply(txs []model.Transaction) {
	for i := range txs {
		txs[i].Category = rs.Categorize(txs[i])
	}
}

func (r compiledRule) matches(tx model.Transaction) bool {
	hay := haystack(tx, r.field)
	if r.re != nil {
		return r.re.MatchString(hay)
	}
	return strings.Contains(strings.ToLower(hay), r.needle)
}

func haystack(tx model.Transaction, field string) string {
	switch field {
	case FieldPayee:
		return tx.Payee
	case FieldPurpose:
		return tx.Purpose
	default:
		return tx.Payee + " " + tx.Purpose
	}
}

func normalizeField(f string) string {
	switch strings.ToLower(strings.TrimSpace(f)) {
	case FieldPayee:
		return FieldPayee
	case FieldPurpose:
		return FieldPurpose
	default:
		return FieldAny
	}
}
