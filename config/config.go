package config

import (
	"os"
)

type Config struct {
	GitHubToken   string
	WebhookSecret string
}

func Load() *Config {
	return &Config{
		GitHubToken:   os.Getenv("GH_TOKEN"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
	}
}
