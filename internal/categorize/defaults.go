package categorize

import (
	_ "embed"
	"bytes"
)

//go:embed rules.yaml
var defaultRulesYAML []byte

// Default returns the built-in rule set shipped with ledgerly.
func Default() (*RuleSet, error) {
	return Load(bytes.NewReader(defaultRulesYAML))
}

// DefaultYAML returns the raw default rules document, useful for `rules dump`
// so users can copy it as a starting point for their own file.
func DefaultYAML() []byte {
	out := make([]byte, len(defaultRulesYAML))
	copy(out, defaultRulesYAML)
	return out
}
