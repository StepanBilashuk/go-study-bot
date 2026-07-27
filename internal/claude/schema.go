package claude

import (
	"encoding/json"
	"fmt"
)

// The types and parsers below mirror the strict JSON schemas the bot asks
// Claude to return (spec §13). Each Parse* function unmarshals AND validates
// ranges/enums, so a malformed or out-of-range reply is rejected and triggers
// the one-shot retry in CompleteJSON. These parsers are the unit-tested core.

// --- calibration (/start) ---

// CalibrationScore is one topic's calibrated confidence.
type CalibrationScore struct {
	Topic      string `json:"topic"`
	Confidence int    `json:"confidence"`
	Note       string `json:"note"`
}

// CalibrationResult is the payload of the calibration prompt.
type CalibrationResult struct {
	Scores []CalibrationScore `json:"scores"`
}

// ParseCalibration validates {"scores":[{topic,confidence 1-5,note}]}.
func ParseCalibration(b []byte) (CalibrationResult, error) {
	var r CalibrationResult
	if err := json.Unmarshal(b, &r); err != nil {
		return CalibrationResult{}, fmt.Errorf("calibration json: %w", err)
	}
	if len(r.Scores) == 0 {
		return CalibrationResult{}, fmt.Errorf("calibration: empty scores")
	}
	for i, s := range r.Scores {
		if s.Topic == "" {
			return CalibrationResult{}, fmt.Errorf("calibration score[%d]: empty topic", i)
		}
		if s.Confidence < 1 || s.Confidence > 5 {
			return CalibrationResult{}, fmt.Errorf("calibration score[%d]: confidence %d out of 1-5", i, s.Confidence)
		}
	}
	return r, nil
}

// --- debrief (/debrief) ---

// Gap is a specific weakness extracted from a debrief.
type Gap struct {
	Topic    string `json:"topic"`
	Detail   string `json:"detail"`
	Severity string `json:"severity"` // low | medium | high
}

// ConfidenceUpdate revises a topic's confidence after a debrief.
type ConfidenceUpdate struct {
	Topic         string `json:"topic"`
	NewConfidence int    `json:"new_confidence"`
}

// GlossaryTerm is an English term the candidate could not produce.
type GlossaryTerm struct {
	Term    string `json:"term"`
	Context string `json:"context"`
}

// DebriefResult is the payload of the debrief prompt.
type DebriefResult struct {
	TopicsTouched     []string           `json:"topics_touched"`
	Gaps              []Gap              `json:"gaps"`
	ConfidenceUpdates []ConfidenceUpdate `json:"confidence_updates"`
	NewGlossaryTerms  []GlossaryTerm     `json:"new_glossary_terms"`
}

// ParseDebrief validates the debrief JSON, including severity enums and
// confidence ranges.
func ParseDebrief(b []byte) (DebriefResult, error) {
	var r DebriefResult
	if err := json.Unmarshal(b, &r); err != nil {
		return DebriefResult{}, fmt.Errorf("debrief json: %w", err)
	}
	for i, g := range r.Gaps {
		switch g.Severity {
		case "low", "medium", "high":
		default:
			return DebriefResult{}, fmt.Errorf("debrief gap[%d]: bad severity %q", i, g.Severity)
		}
	}
	for i, cu := range r.ConfidenceUpdates {
		if cu.NewConfidence < 1 || cu.NewConfidence > 5 {
			return DebriefResult{}, fmt.Errorf("debrief confidence_update[%d]: %d out of 1-5", i, cu.NewConfidence)
		}
	}
	return r, nil
}

// --- drill scoring (/drill) ---

// DrillScore is the payload of the drill-score prompt.
type DrillScore struct {
	Score   int    `json:"score"`
	Outcome string `json:"outcome"`
}

// ParseDrillScore validates {"score":1-5,"outcome":...}.
func ParseDrillScore(b []byte) (DrillScore, error) {
	var r DrillScore
	if err := json.Unmarshal(b, &r); err != nil {
		return DrillScore{}, fmt.Errorf("drill score json: %w", err)
	}
	if r.Score < 1 || r.Score > 5 {
		return DrillScore{}, fmt.Errorf("drill score %d out of 1-5", r.Score)
	}
	return r, nil
}

// --- pattern-recognition quiz (/quiz, Phase 2) ---

// QuizItem is one multiple-choice item. Answer is the 0-based index of the
// correct option.
type QuizItem struct {
	Question    string   `json:"question"`
	Options     []string `json:"options"`
	Answer      int      `json:"answer"`
	Explanation string   `json:"explanation"`
}

// Quiz is the payload of the quiz prompt.
type Quiz struct {
	Items []QuizItem `json:"items"`
}

// ParseQuiz validates a quiz: at least one item, each with 2-4 options and an
// in-range answer index.
func ParseQuiz(b []byte) (Quiz, error) {
	var q Quiz
	if err := json.Unmarshal(b, &q); err != nil {
		return Quiz{}, fmt.Errorf("quiz json: %w", err)
	}
	if len(q.Items) == 0 {
		return Quiz{}, fmt.Errorf("quiz: no items")
	}
	for i, it := range q.Items {
		if it.Question == "" {
			return Quiz{}, fmt.Errorf("quiz item[%d]: empty question", i)
		}
		if len(it.Options) < 2 || len(it.Options) > 4 {
			return Quiz{}, fmt.Errorf("quiz item[%d]: %d options (want 2-4)", i, len(it.Options))
		}
		if it.Answer < 0 || it.Answer >= len(it.Options) {
			return Quiz{}, fmt.Errorf("quiz item[%d]: answer %d out of range", i, it.Answer)
		}
	}
	return q, nil
}

// --- STAR story extraction (/story, Phase 3) ---

// Story is a mined STAR story. A measurable result is REQUIRED (spec §4, §2):
// a story with no number is rejected.
type Story struct {
	Title        string   `json:"title"`
	Situation    string   `json:"situation"`
	Task         string   `json:"task"`
	Action       string   `json:"action"`
	Result       string   `json:"result"`
	Metrics      []string `json:"metrics"`
	TechTags     []string `json:"tech_tags"`
	Competencies []string `json:"competencies"`
}

// ParseStory validates an extracted story and enforces the metric requirement.
func ParseStory(b []byte) (Story, error) {
	var s Story
	if err := json.Unmarshal(b, &s); err != nil {
		return Story{}, fmt.Errorf("story json: %w", err)
	}
	if s.Situation == "" || s.Action == "" || s.Result == "" {
		return Story{}, fmt.Errorf("story: situation, action and result are all required")
	}
	if len(s.Metrics) == 0 {
		return Story{}, fmt.Errorf("story has no measurable result — every story needs a number")
	}
	return s, nil
}

// --- company research import (/importcompany, Phase 4) ---

// ResearchStage is one round of a researched interview process.
type ResearchStage struct {
	Stage       int    `json:"stage"`
	Name        string `json:"name"`
	Format      string `json:"format"`
	DurationMin int    `json:"duration_min"`
}

// CompanyResearch is the JSON returned by the company-research prompt (spec
// §13). Only the fields the bot maps into a company definition are typed;
// unknown fields are ignored.
type CompanyResearch struct {
	Name             string          `json:"name"`
	Locations        []string        `json:"locations"`
	Stack            []string        `json:"stack"`
	InterviewProcess []ResearchStage `json:"interview_process"`
	Values           []string        `json:"values"`
	RequiredTopics   []string        `json:"required_topics"`
	NewTopics        []string        `json:"new_topics"`
	Confidence       string          `json:"confidence"`
}

// ParseCompanyResearch validates the researched company JSON.
func ParseCompanyResearch(b []byte) (CompanyResearch, error) {
	var c CompanyResearch
	if err := json.Unmarshal(b, &c); err != nil {
		return CompanyResearch{}, fmt.Errorf("company research json: %w", err)
	}
	if c.Name == "" {
		return CompanyResearch{}, fmt.Errorf("company research: empty name")
	}
	return c, nil
}
