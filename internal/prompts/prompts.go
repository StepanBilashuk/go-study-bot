// Package prompts loads Claude prompt templates from prompts/*.yaml at runtime,
// so they can be tuned without redeploying (spec §13). Each file holds a single
// `template:` string; {{placeholders}} are substituted at render time.
package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Set is the loaded collection of prompt templates, keyed by file base name
// (e.g. "calibration", "drill-estimation").
type Set struct {
	templates map[string]string
}

type promptFile struct {
	Template string `yaml:"template"`
}

// Load reads every *.yaml under dir into a Set. A malformed file, an unexpected
// key, or an empty template is an error, so bad prompts fail at startup/reload.
func Load(dir string) (*Set, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("glob prompts: %w", err)
	}
	set := &Set{templates: make(map[string]string)}
	for _, path := range matches {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", path, err)
		}
		dec := yaml.NewDecoder(f)
		dec.KnownFields(true)
		var pf promptFile
		decErr := dec.Decode(&pf)
		f.Close()
		if decErr != nil {
			return nil, fmt.Errorf("parse %s: %w", path, decErr)
		}
		if strings.TrimSpace(pf.Template) == "" {
			return nil, fmt.Errorf("%s: empty template", path)
		}
		name := strings.TrimSuffix(filepath.Base(path), ".yaml")
		set.templates[name] = pf.Template
	}
	return set, nil
}

// Count returns how many templates are loaded.
func (s *Set) Count() int { return len(s.templates) }

// Render substitutes {{key}} placeholders and returns the prompt text. It is an
// error if the named prompt is missing.
func (s *Set) Render(name string, vars map[string]string) (string, error) {
	t, ok := s.templates[name]
	if !ok {
		return "", fmt.Errorf("prompt %q not found", name)
	}
	for k, v := range vars {
		t = strings.ReplaceAll(t, "{{"+k+"}}", v)
	}
	return t, nil
}
