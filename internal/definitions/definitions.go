// Package definitions loads the immutable content of the bot — topics, drills,
// books, resources and companies — from YAML on disk into in-memory structs.
//
// Definitions live in YAML in the repo; only mutable state lives in Postgres
// (spec §3, §5). Load is called at startup and on /reload. Any structural
// problem (malformed YAML, unknown field, dangling cross-reference) returns an
// error so an invalid data set prevents startup with a clear message.
package definitions

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Track is the study track a topic belongs to.
type Track string

const (
	TrackAlgorithms   Track = "algorithms"
	TrackSystemDesign Track = "system-design"
)

// DrillKind is one of the four daily process-drill kinds (spec §10).
type DrillKind string

const (
	DrillEstimation    DrillKind = "estimation"
	DrillContradiction DrillKind = "contradiction"
	DrillClarify       DrillKind = "clarify"
	DrillNextStep      DrillKind = "next-step"
)

// Topic is one rung on the five-stage ladder (spec §5, §8, §9).
type Topic struct {
	Slug      string   `yaml:"slug"`
	Track     Track    `yaml:"track"`
	Name      string   `yaml:"name"`
	Priority  int      `yaml:"priority"`
	DependsOn []string `yaml:"depends_on"`
	Gate      string   `yaml:"gate"`
	EstHours  int      `yaml:"est_hours"`
}

// Drill is a daily process drill (spec §5, §10).
type Drill struct {
	Slug        string    `yaml:"slug"`
	Kind        DrillKind `yaml:"kind"`
	Name        string    `yaml:"name"`
	DurationMin int       `yaml:"duration_min"`
	Prompt      string    `yaml:"prompt"` // path to a prompt file under prompts/
}

// Book is a study book. Edition is load-bearing: translations paginate
// differently from the originals (spec §5).
type Book struct {
	Title   string `yaml:"title"`
	Author  string `yaml:"author"`
	Edition string `yaml:"edition"`
}

// Resource attaches a study resource to a topic AND a ladder stage (spec §11).
type Resource struct {
	Topic   string  `yaml:"topic"`
	Stage   int     `yaml:"stage"`
	Type    string  `yaml:"type"` // book | video | article
	Source  string  `yaml:"source"`
	Chapter int     `yaml:"chapter"`
	Section string  `yaml:"section"`
	Pages   *string `yaml:"pages"` // null until filled by hand from the reader's edition
	EstMin  int     `yaml:"est_min"`
}

// InterviewStage is one round in a company's interview process (spec §5).
type InterviewStage struct {
	Stage       int    `yaml:"stage"`
	Name        string `yaml:"name"`
	Format      string `yaml:"format"`
	DurationMin int    `yaml:"duration_min"`
}

// Company is a target employer (spec §5, §12).
type Company struct {
	Slug             string           `yaml:"slug"`
	Name             string           `yaml:"name"`
	Locations        []string         `yaml:"locations"`
	Stack            []string         `yaml:"stack"`
	InterviewProcess []InterviewStage `yaml:"interview_process"`
	RequiredTopics   []string         `yaml:"required_topics"`
	Values           []string         `yaml:"values"`
	Referral         bool             `yaml:"referral"`
	ResearchedAt     string           `yaml:"researched_at"`
	Confidence       string           `yaml:"confidence"`
}

// Definitions is the complete in-memory content set.
type Definitions struct {
	Topics    map[string]Topic
	Drills    map[string]Drill
	Books     map[string]Book
	Resources []Resource
	Companies map[string]Company
	Theory    map[string]string // topic slug -> seeded theory markdown
	Designs   map[string]string // design slug -> seeded "design X" walkthrough
	Prep      map[string]string // company slug -> interview prep card
}

// Load reads and validates every definition file under dataDir. Missing
// directories are treated as "nothing to load" (valid); malformed or
// inconsistent content is an error.
func Load(dataDir string) (*Definitions, error) {
	defs := &Definitions{
		Topics:    make(map[string]Topic),
		Drills:    make(map[string]Drill),
		Books:     make(map[string]Book),
		Companies: make(map[string]Company),
		Theory:    make(map[string]string),
		Designs:   make(map[string]string),
		Prep:      make(map[string]string),
	}

	if err := defs.loadTopics(filepath.Join(dataDir, "topics")); err != nil {
		return nil, err
	}
	if err := defs.loadTheory(filepath.Join(dataDir, "theory")); err != nil {
		return nil, err
	}
	if err := defs.loadMarkdown(filepath.Join(dataDir, "designs"), defs.Designs); err != nil {
		return nil, err
	}
	if err := defs.loadMarkdown(filepath.Join(dataDir, "prep"), defs.Prep); err != nil {
		return nil, err
	}
	if err := defs.loadDrills(filepath.Join(dataDir, "drills")); err != nil {
		return nil, err
	}
	if err := defs.loadBooks(filepath.Join(dataDir, "books.yaml")); err != nil {
		return nil, err
	}
	if err := defs.loadResources(filepath.Join(dataDir, "resources.yaml")); err != nil {
		return nil, err
	}
	if err := defs.loadCompanies(filepath.Join(dataDir, "companies")); err != nil {
		return nil, err
	}

	if err := defs.validate(); err != nil {
		return nil, err
	}
	return defs, nil
}

func (d *Definitions) loadTopics(dir string) error {
	files, err := yamlFiles(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		var topics []Topic
		if err := decodeFile(f, &topics); err != nil {
			return err
		}
		for _, t := range topics {
			if _, dup := d.Topics[t.Slug]; dup {
				return fmt.Errorf("%s: duplicate topic slug %q", f, t.Slug)
			}
			d.Topics[t.Slug] = t
		}
	}
	return nil
}

// loadTheory reads data/theory/<slug>.md into the Theory map (seeded in-bot
// theory shown when a topic is tapped).
func (d *Definitions) loadTheory(dir string) error {
	return d.loadMarkdown(dir, d.Theory)
}

// loadMarkdown reads every <name>.md under dir into dest, keyed by base name.
// A missing directory is fine (nothing to load).
func (d *Definitions) loadMarkdown(dir string, dest map[string]string) error {
	if !fileExists(dir) {
		return nil
	}
	matches, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return fmt.Errorf("glob %s: %w", dir, err)
	}
	for _, f := range matches {
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		dest[strings.TrimSuffix(filepath.Base(f), ".md")] = string(content)
	}
	return nil
}

func (d *Definitions) loadDrills(dir string) error {
	files, err := yamlFiles(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		var drills []Drill
		if err := decodeFile(f, &drills); err != nil {
			return err
		}
		for _, dr := range drills {
			if _, dup := d.Drills[dr.Slug]; dup {
				return fmt.Errorf("%s: duplicate drill slug %q", f, dr.Slug)
			}
			d.Drills[dr.Slug] = dr
		}
	}
	return nil
}

func (d *Definitions) loadBooks(path string) error {
	if !fileExists(path) {
		return nil
	}
	return decodeFile(path, &d.Books)
}

func (d *Definitions) loadResources(path string) error {
	if !fileExists(path) {
		return nil
	}
	return decodeFile(path, &d.Resources)
}

func (d *Definitions) loadCompanies(dir string) error {
	files, err := yamlFiles(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		var c Company
		if err := decodeFile(f, &c); err != nil {
			return err
		}
		if _, dup := d.Companies[c.Slug]; dup {
			return fmt.Errorf("%s: duplicate company slug %q", f, c.Slug)
		}
		d.Companies[c.Slug] = c
	}
	return nil
}

// validate enforces cross-references and enums. Each failure names the file or
// entity so the operator can fix the YAML without guessing.
func (d *Definitions) validate() error {
	for slug, t := range d.Topics {
		if t.Slug == "" {
			return fmt.Errorf("topic with empty slug (name %q)", t.Name)
		}
		if t.Track != TrackAlgorithms && t.Track != TrackSystemDesign {
			return fmt.Errorf("topic %q: invalid track %q (want algorithms|system-design)", slug, t.Track)
		}
		for _, dep := range t.DependsOn {
			if _, ok := d.Topics[dep]; !ok {
				return fmt.Errorf("topic %q depends_on unknown topic %q", slug, dep)
			}
		}
	}

	for slug, dr := range d.Drills {
		switch dr.Kind {
		case DrillEstimation, DrillContradiction, DrillClarify, DrillNextStep:
		default:
			return fmt.Errorf("drill %q: invalid kind %q", slug, dr.Kind)
		}
	}

	for i, r := range d.Resources {
		if _, ok := d.Topics[r.Topic]; !ok {
			return fmt.Errorf("resource[%d]: references unknown topic %q", i, r.Topic)
		}
		if r.Type == "book" && r.Source != "" {
			if _, ok := d.Books[r.Source]; !ok {
				return fmt.Errorf("resource[%d] (topic %q): references unknown book %q", i, r.Topic, r.Source)
			}
		}
	}

	for slug, c := range d.Companies {
		for _, rt := range c.RequiredTopics {
			if _, ok := d.Topics[rt]; !ok {
				return fmt.Errorf("company %q: required_topic %q is not a known topic slug", slug, rt)
			}
		}
	}
	return nil
}

// --- helpers ---

// yamlFiles returns the sorted .yaml/.yml files directly under dir. A missing
// directory yields no files and no error.
func yamlFiles(dir string) ([]string, error) {
	if !fileExists(dir) {
		return nil, nil
	}
	var files []string
	for _, ext := range []string{"*.yaml", "*.yml"} {
		matches, err := filepath.Glob(filepath.Join(dir, ext))
		if err != nil {
			return nil, fmt.Errorf("glob %s: %w", dir, err)
		}
		files = append(files, matches...)
	}
	sort.Strings(files)
	return files, nil
}

// decodeFile strictly decodes a single YAML document into out. KnownFields
// makes an unexpected key a hard error, so typos surface at startup.
func decodeFile(path string, out any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
