package profile

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed builtin/*.yaml
var builtinFS embed.FS

// Builtins parses and returns every profile shipped with ledgerly, sorted by
// name.
func Builtins() ([]*Profile, error) {
	entries, err := fs.ReadDir(builtinFS, "builtin")
	if err != nil {
		return nil, err
	}
	profiles := make([]*Profile, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		f, err := builtinFS.Open("builtin/" + e.Name())
		if err != nil {
			return nil, err
		}
		p, err := Load(f)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("builtin %s: %w", e.Name(), err)
		}
		profiles = append(profiles, p)
	}
	sort.Slice(profiles, func(i, j int) bool { return profiles[i].Name < profiles[j].Name })
	return profiles, nil
}

// BuiltinNames returns the names of the shipped profiles.
func BuiltinNames() ([]string, error) {
	ps, err := Builtins()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(ps))
	for i, p := range ps {
		names[i] = p.Name
	}
	return names, nil
}

// Resolve loads a profile from either a built-in name or a YAML file path. A
// value ending in .yaml/.yml or containing a path separator is treated as a
// file; otherwise it is looked up among the built-ins.
func Resolve(nameOrPath string) (*Profile, error) {
	if strings.HasSuffix(nameOrPath, ".yaml") || strings.HasSuffix(nameOrPath, ".yml") ||
		strings.ContainsAny(nameOrPath, "/\\") {
		return LoadFile(nameOrPath)
	}
	return LoadBuiltin(nameOrPath)
}

// LoadBuiltin returns the shipped profile with the given name, matched
// case-insensitively.
func LoadBuiltin(name string) (*Profile, error) {
	ps, err := Builtins()
	if err != nil {
		return nil, err
	}
	for _, p := range ps {
		if strings.EqualFold(p.Name, name) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no built-in profile named %q", name)
}
