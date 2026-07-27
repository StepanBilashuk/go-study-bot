// Package config loads all runtime configuration from environment variables.
// Nothing is hardcoded and no config files are read: secrets arrive via the
// systemd EnvironmentFile in production (spec §15).
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds every value the bot needs at runtime.
type Config struct {
	TelegramToken   string // TELEGRAM_BOT_TOKEN (required)
	DatabaseURL     string // DATABASE_URL, pgx-style DSN (required)
	AnthropicAPIKey string // ANTHROPIC_API_KEY (required from Step 3 onward)
	AnthropicModel  string // ANTHROPIC_MODEL, tunable without a rebuild
	DataDir         string // directory holding data/** YAML definitions
	PromptsDir      string // directory holding prompts/*.yaml
	PushHour        int    // PUSH_HOUR, local hour (0-23) for the daily /today push
	Gamification    bool   // GAMIFICATION feature flag (Phase 5): XP, levels, streak
}

// Load reads configuration from the environment. Missing required variables
// produce a single clear error so the process fails fast at startup.
func Load() (Config, error) {
	c := Config{
		TelegramToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		AnthropicModel:  getenvDefault("ANTHROPIC_MODEL", "claude-opus-5"),
		DataDir:         getenvDefault("DATA_DIR", "data"),
		PromptsDir:      getenvDefault("PROMPTS_DIR", "prompts"),
		PushHour:        8,
	}

	if v := os.Getenv("PUSH_HOUR"); v != "" {
		h, err := strconv.Atoi(v)
		if err != nil || h < 0 || h > 23 {
			return Config{}, fmt.Errorf("PUSH_HOUR must be an integer 0-23, got %q", v)
		}
		c.PushHour = h
	}
	switch os.Getenv("GAMIFICATION") {
	case "", "0", "false", "off":
		c.Gamification = false
	default:
		c.Gamification = true
	}

	// ANTHROPIC_API_KEY is only exercised by Claude-backed commands (Step 3+),
	// so the bot is allowed to boot without it during Step 1.
	var missing []string
	if c.TelegramToken == "" {
		missing = append(missing, "TELEGRAM_BOT_TOKEN")
	}
	if c.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}
	return c, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
