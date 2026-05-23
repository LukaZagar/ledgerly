package profile

import (
	"fmt"
	"strings"
)

// Detect picks the built-in profile whose signature columns are all present in
// the given CSV header. When several match, the one with the most signature
// columns wins, so a more specific profile beats a generic one. It returns an
// error when nothing matches.
func Detect(header []string) (*Profile, error) {
	present := make(map[string]bool, len(header))
	for _, h := range header {
		present[strings.TrimSpace(h)] = true
	}

	profiles, err := Builtins()
	if err != nil {
		return nil, err
	}

	var best *Profile
	bestScore := 0
	for _, p := range profiles {
		if len(p.Signature) == 0 {
			continue
		}
		if matchesSignature(p.Signature, present) && len(p.Signature) > bestScore {
			best = p
			bestScore = len(p.Signature)
		}
	}
	if best == nil {
		return nil, fmt.Errorf("could not auto-detect a profile from header: %s", strings.Join(header, ", "))
	}
	return best, nil
}

func matchesSignature(signature []string, present map[string]bool) bool {
	for _, col := range signature {
		if !present[strings.TrimSpace(col)] {
			return false
		}
	}
	return true
}
