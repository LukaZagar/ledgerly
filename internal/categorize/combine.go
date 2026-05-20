package categorize

import "io"

// Combine merges several rule sets into one, preserving order. Rules from
// earlier sets are evaluated first, so passing user rules before the defaults
// lets a user override how a merchant is categorized without restating the
// whole list.
func Combine(sets ...*RuleSet) *RuleSet {
	merged := &RuleSet{}
	for _, s := range sets {
		if s == nil {
			continue
		}
		merged.rules = append(merged.rules, s.rules...)
	}
	return merged
}

// LoadWithDefaults loads a user rule file and stacks it on top of the built-in
// defaults: the user's rules win, anything they don't cover falls through to
// the defaults.
func LoadWithDefaults(r io.Reader) (*RuleSet, error) {
	user, err := Load(r)
	if err != nil {
		return nil, err
	}
	def, err := Default()
	if err != nil {
		return nil, err
	}
	return Combine(user, def), nil
}

// Len reports how many rules are in the set.
func (rs *RuleSet) Len() int { return len(rs.rules) }
