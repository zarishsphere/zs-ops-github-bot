package config

import (
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/v60/github"
)

// Config holds application configuration
type Config struct {
	GitHubToken      string
	WebhookSecret    string
	Organization     string
	AutoMergeEnabled bool
	StalePRThreshold time.Duration
	GitHubClient     *github.Client
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		GitHubToken:      getEnv("GH_TOKEN", ""),
		WebhookSecret:    getEnv("WEBHOOK_SECRET", ""),
		Organization:     getEnv("GITHUB_ORG", "zarishsphere"),
		AutoMergeEnabled: getEnvAsBool("AUTO_MERGE_ENABLED", true),
		StalePRThreshold: getEnvAsDuration("STALE_PR_THRESHOLD", 720*time.Hour),
	}

	// Initialize GitHub client
	if cfg.GitHubToken != "" {
		cfg.GitHubClient = github.NewClient(nil).WithAuthToken(cfg.GitHubToken)
	} else {
		cfg.GitHubClient = github.NewClient(nil)
	}

	return cfg, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as a boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as a duration
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
