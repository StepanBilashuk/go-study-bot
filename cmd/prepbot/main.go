// Command prepbot is the single static binary that runs the interview-prep
// Telegram bot: a long-polling loop plus (from Step 3) a scheduler goroutine.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"prepbot/internal/config"
	"prepbot/internal/db"
	"prepbot/internal/definitions"
	"prepbot/internal/prompts"
	"prepbot/internal/tgbot"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Load definitions first: invalid YAML must prevent startup with a clear
	// message before we touch the database or Telegram (spec §3, §5).
	defs, err := definitions.Load(cfg.DataDir)
	if err != nil {
		return fmt.Errorf("load definitions: %w", err)
	}
	slog.Info("definitions loaded",
		"topics", len(defs.Topics),
		"drills", len(defs.Drills),
		"books", len(defs.Books),
		"resources", len(defs.Resources),
		"companies", len(defs.Companies),
	)

	ps, err := prompts.Load(cfg.PromptsDir)
	if err != nil {
		return fmt.Errorf("load prompts: %w", err)
	}
	slog.Info("prompts loaded", "count", ps.Count())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	database, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer database.Close()

	if err := database.Migrate(ctx); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	slog.Info("migrations applied")

	b, err := tgbot.New(cfg, database, defs, ps)
	if err != nil {
		return err
	}

	slog.Info("prepbot starting; telegram long polling")
	b.Start(ctx) // blocks until ctx is cancelled by SIGINT/SIGTERM
	slog.Info("prepbot stopped cleanly")
	return nil
}
