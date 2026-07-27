package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const helpText = `prepbot — interview-prep tracker

Tap ☰ (left of the input) for every command, or /menu to pin the quick buttons.
AI commands emit a prompt you run in YOUR Claude, then paste the JSON reply
back here (like /newcompany → /importcompany). /cancel drops a pending paste.

Daily loop
  /today            3 topics + 1 drill + resources (never more than 3)
  /done <slug>      advance a topic a stage (or name the missing prerequisite)
  /drill            emit a drill of your weakest kind → paste the score JSON back
  /debrief <text>   emit an extraction prompt → paste gaps JSON back

Practice
  /quiz <slug>      emit a 10-item recognition quiz (do it in Claude; 8/10 = gate)
  /boss             system-design mock brief (once the gate is met)
  /boss behavioral  behavioral mock, targeting your uncovered competencies

Stories & language
  /story            mine a STAR story → paste the story JSON back
  /stories          competency matrix (who still has no story)
  /glossary         review English terms you couldn't produce

Companies
  /ready            readiness % per company + named blocker
  /newcompany <n>   emit a research prompt to run in Claude with web search
  /importcompany    import the returned JSON as a company definition
  /interview <slug> <YYYY-MM-DD>   set a date; /today reorders around it

Overview & setup
  /progress         compact per-track summary
  /stats            XP, streak, mastery (needs GAMIFICATION=true)
  /start            emit a calibration prompt → paste the scores JSON back
  /push on|off      toggle your daily plan push
  /whoami           show your chat id
  /cancel           drop a pending paste-back
  /reload           re-read YAML from disk
  /ping             liveness check`

func (b *Bot) handleHelp(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	// Install the quick-bar alongside the full list.
	b.replyKb(ctx, update.Message.Chat.ID, helpText)
}

func (b *Bot) handleWhoami(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID := update.Message.Chat.ID
	username := ""
	if update.Message.From != nil {
		username = update.Message.From.Username
	}
	if username != "" {
		b.reply(ctx, chatID, fmt.Sprintf("chat id: %d\nusername: @%s", chatID, username))
	} else {
		b.reply(ctx, chatID, fmt.Sprintf("chat id: %d", chatID))
	}
}
