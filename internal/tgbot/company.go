package tgbot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gopkg.in/yaml.v3"

	"prepbot/internal/claude"
	"prepbot/internal/db"
	"prepbot/internal/definitions"
)

// handleNewCompany (Phase 4) emits a research prompt to paste into Claude with
// web search. It fills {{company}} and {{topic_slugs}} — no API call here.
func (b *Bot) handleNewCompany(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	name := commandArg(update.Message.Text, "/newcompany")
	if name == "" {
		b.reply(ctx, chatID, "Usage: /newcompany <company name>")
		return
	}

	defs := b.getDefs()
	slugs := make([]string, 0, len(defs.Topics))
	for s := range defs.Topics {
		slugs = append(slugs, s)
	}
	sort.Strings(slugs)

	prompt, err := b.getPrompts().Render("company-research", map[string]string{
		"company":     name,
		"topic_slugs": strings.Join(slugs, ", "),
	})
	if err != nil {
		b.reply(ctx, chatID, "Prompt error: "+err.Error())
		return
	}
	b.reply(ctx, chatID, "Paste this into Claude WITH web search, then send the JSON back with /importcompany:\n\n"+prompt)
}

// handleImportCompany parks the chat to receive the researched JSON.
func (b *Bot) handleImportCompany(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	b.setConv(chatID, &conversation{mode: modeImportCompany})
	b.reply(ctx, chatID, "Send the JSON returned by the research prompt. I'll validate it and write data/companies/<slug>.yaml.")
}

// handleImportCompanyJSON validates the research JSON, writes the company
// definition to disk, and reloads (spec §7).
func (b *Bot) handleImportCompanyJSON(ctx context.Context, chatID int64, text string) {
	research, err := claude.ParseCompanyResearch([]byte(claude.ExtractJSON(text)))
	if err != nil {
		b.reply(ctx, chatID, "Invalid company JSON: "+err.Error()+"\nPaste valid JSON, or /today to cancel.")
		return
	}
	b.clearConv(chatID)

	defs := b.getDefs()
	slug := slugify(research.Name)

	// Keep only required_topics that map to known slugs; report the rest.
	var known, unknown []string
	for _, rt := range research.RequiredTopics {
		if _, ok := defs.Topics[rt]; ok {
			known = append(known, rt)
		} else {
			unknown = append(unknown, rt)
		}
	}

	company := definitions.Company{
		Slug:           slug,
		Name:           research.Name,
		Locations:      research.Locations,
		Stack:          research.Stack,
		RequiredTopics: known,
		Values:         research.Values,
		Referral:       false,
		ResearchedAt:   time.Now().Format("2006-01-02"),
		Confidence:     research.Confidence,
	}
	for _, s := range research.InterviewProcess {
		company.InterviewProcess = append(company.InterviewProcess, definitions.InterviewStage{
			Stage:       s.Stage,
			Name:        s.Name,
			Format:      s.Format,
			DurationMin: s.DurationMin,
		})
	}

	out, err := yaml.Marshal(company)
	if err != nil {
		b.reply(ctx, chatID, "YAML marshal error: "+err.Error())
		return
	}
	path := filepath.Join(b.cfg.DataDir, "companies", slug+".yaml")
	if err := os.WriteFile(path, out, 0o644); err != nil {
		b.reply(ctx, chatID, "Write error: "+err.Error())
		return
	}

	// Reload so the new company is live; if it fails validation, say so.
	newDefs, err := definitions.Load(b.cfg.DataDir)
	if err != nil {
		b.reply(ctx, chatID, "Wrote "+path+" but reload failed:\n"+err.Error())
		return
	}
	b.mu.Lock()
	b.defs = newDefs
	b.mu.Unlock()

	msg := fmt.Sprintf("✅ Imported %s → %s (%d required topics).", research.Name, path, len(known))
	if len(unknown) > 0 {
		msg += "\nDropped non-slug topics (add them to data/topics if they matter): " + strings.Join(unknown, ", ")
	}
	b.reply(ctx, chatID, msg)
}

// handleReady shows a readiness score per company with its named blocker
// (spec §7, §12).
func (b *Bot) handleReady(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	prog, err := b.db.ListProgress(ctx, chatID)
	if err != nil {
		b.reply(ctx, chatID, "DB error: "+err.Error())
		return
	}
	defs := b.getDefs()

	type row struct {
		name    string
		score   int
		blocker string
	}
	var rows []row
	for _, c := range defs.Companies {
		score, blocker := readiness(c, prog)
		rows = append(rows, row{c.Name, score, blocker})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].score > rows[j].score })

	var sb strings.Builder
	sb.WriteString("🎯 Readiness\n")
	for _, r := range rows {
		if r.blocker == "" {
			fmt.Fprintf(&sb, "%3d%%  %s → ready\n", r.score, r.name)
		} else {
			fmt.Fprintf(&sb, "%3d%%  %s → blocker: %s\n", r.score, r.name, r.blocker)
		}
	}
	b.reply(ctx, chatID, strings.TrimRight(sb.String(), "\n"))
}

// readiness is the weighted coverage of required_topics by current confidence
// (spec §12). Returns a 0-100 score and the lowest-confidence required topic as
// the named blocker (empty when everything is at confidence ≥4).
func readiness(c definitions.Company, prog map[string]db.Progress) (int, string) {
	if len(c.RequiredTopics) == 0 {
		return 0, "no required topics"
	}
	sum := 0.0
	blocker, blockerConf := "", 6
	for _, rt := range c.RequiredTopics {
		conf := prog[rt].Confidence
		sum += float64(conf) / 5.0
		if conf < blockerConf {
			blocker, blockerConf = rt, conf
		}
	}
	score := int(sum / float64(len(c.RequiredTopics)) * 100)
	if blockerConf >= 4 {
		blocker = "" // ready
	}
	return score, blocker
}

// slugify turns a company name into a filesystem/slug-safe token.
func slugify(name string) string {
	var sb strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			sb.WriteRune(r)
			prevDash = false
		default:
			if !prevDash && sb.Len() > 0 {
				sb.WriteByte('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(sb.String(), "-")
}
