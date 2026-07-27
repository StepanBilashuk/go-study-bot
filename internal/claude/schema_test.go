package claude

import "testing"

func TestExtractJSON(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"plain", `{"a":1}`, `{"a":1}`},
		{"fenced json", "```json\n{\"a\":1}\n```", `{"a":1}`},
		{"fenced bare", "```\n{\"a\":1}\n```", `{"a":1}`},
		{"surrounding whitespace", "  {\"a\":1}  ", `{"a":1}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractJSON(tt.in); got != tt.want {
				t.Errorf("ExtractJSON(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseCalibration(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		r, err := ParseCalibration([]byte(`{"scores":[{"topic":"sliding-window","confidence":4,"note":"used it"}]}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(r.Scores) != 1 || r.Scores[0].Topic != "sliding-window" || r.Scores[0].Confidence != 4 {
			t.Fatalf("unexpected result: %+v", r)
		}
	})

	bad := map[string]string{
		"empty scores":        `{"scores":[]}`,
		"confidence too high":  `{"scores":[{"topic":"x","confidence":6}]}`,
		"confidence too low":   `{"scores":[{"topic":"x","confidence":0}]}`,
		"empty topic":          `{"scores":[{"topic":"","confidence":3}]}`,
		"malformed json":       `{"scores":`,
	}
	for name, in := range bad {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseCalibration([]byte(in)); err == nil {
				t.Errorf("expected error for %s", name)
			}
		})
	}
}

func TestParseDebrief(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		in := `{
			"topics_touched":["sliding-window"],
			"gaps":[{"topic":"sliding-window","detail":"variable window","severity":"high"}],
			"confidence_updates":[{"topic":"sliding-window","new_confidence":2}],
			"new_glossary_terms":[{"term":"amortized","context":"complexity"}]
		}`
		r, err := ParseDebrief([]byte(in))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(r.Gaps) != 1 || r.Gaps[0].Severity != "high" {
			t.Fatalf("unexpected gaps: %+v", r.Gaps)
		}
		if len(r.ConfidenceUpdates) != 1 || r.ConfidenceUpdates[0].NewConfidence != 2 {
			t.Fatalf("unexpected updates: %+v", r.ConfidenceUpdates)
		}
	})

	t.Run("empty payload is valid", func(t *testing.T) {
		if _, err := ParseDebrief([]byte(`{}`)); err != nil {
			t.Errorf("empty debrief should be valid: %v", err)
		}
	})

	bad := map[string]string{
		"bad severity":       `{"gaps":[{"topic":"x","detail":"y","severity":"critical"}]}`,
		"confidence out":     `{"confidence_updates":[{"topic":"x","new_confidence":9}]}`,
		"malformed":          `not json`,
	}
	for name, in := range bad {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseDebrief([]byte(in)); err == nil {
				t.Errorf("expected error for %s", name)
			}
		})
	}
}

func TestParseDrillScore(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		r, err := ParseDrillScore([]byte(`{"score":3,"outcome":"stated assumptions"}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.Score != 3 {
			t.Fatalf("got score %d", r.Score)
		}
	})
	for name, in := range map[string]string{
		"too high": `{"score":6}`,
		"too low":  `{"score":0}`,
		"bad json": `{`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseDrillScore([]byte(in)); err == nil {
				t.Errorf("expected error for %s", name)
			}
		})
	}
}

func TestParseQuiz(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		in := `{"items":[{"question":"Longest substring without repeats?","options":["two pointers","sliding window","dfs","dp"],"answer":1,"explanation":"variable window"}]}`
		q, err := ParseQuiz([]byte(in))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(q.Items) != 1 || q.Items[0].Answer != 1 {
			t.Fatalf("unexpected quiz: %+v", q)
		}
	})
	for name, in := range map[string]string{
		"no items":        `{"items":[]}`,
		"answer OOB":      `{"items":[{"question":"q","options":["a","b"],"answer":5}]}`,
		"too few options": `{"items":[{"question":"q","options":["a"],"answer":0}]}`,
		"empty question":  `{"items":[{"question":"","options":["a","b"],"answer":0}]}`,
		"bad json":        `{`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseQuiz([]byte(in)); err == nil {
				t.Errorf("expected error for %s", name)
			}
		})
	}
}

func TestParseStory(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		in := `{"title":"Notification delivery","situation":"deliverability was 85%","task":"raise it","action":"added outbox + retries","result":"99%+ delivery","metrics":["85% -> 99%"],"competencies":["ownership"]}`
		s, err := ParseStory([]byte(in))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(s.Metrics) != 1 {
			t.Fatalf("unexpected metrics: %+v", s.Metrics)
		}
	})
	for name, in := range map[string]string{
		"no metric":     `{"situation":"s","action":"a","result":"r","metrics":[]}`,
		"missing action": `{"situation":"s","result":"r","metrics":["x"]}`,
		"bad json":      `{`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseStory([]byte(in)); err == nil {
				t.Errorf("expected error for %s", name)
			}
		})
	}
}

func TestParseCompanyResearch(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		in := `{"name":"Aiven","locations":["Helsinki"],"stack":["Go"],"required_topics":["kafka-log-based"],"confidence":"medium"}`
		c, err := ParseCompanyResearch([]byte(in))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Name != "Aiven" {
			t.Fatalf("got %q", c.Name)
		}
	})
	for name, in := range map[string]string{
		"empty name": `{"name":""}`,
		"bad json":   `{`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseCompanyResearch([]byte(in)); err == nil {
				t.Errorf("expected error for %s", name)
			}
		})
	}
}
