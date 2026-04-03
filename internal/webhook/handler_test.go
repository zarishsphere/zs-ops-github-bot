package webhook

import (
	"testing"

	"github.com/google/go-github/v60/github"
	"github.com/zarishsphere/zs-ops-github-bot/internal/config"
)

func TestNewHandler(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	handler := NewHandler(cfg)
	if handler == nil {
		t.Fatal("NewHandler returned nil")
	}

	if handler.config != cfg {
		t.Error("Handler config not set correctly")
	}
}

func TestIsRFCIssue(t *testing.T) {
	cfg, _ := config.LoadConfig()
	handler := NewHandler(cfg)

	tests := []struct {
		name     string
		title    string
		labels   []string
		expected bool
	}{
		{
			name:     "RFC in title",
			title:    "RFC: Add new feature",
			labels:   []string{},
			expected: true,
		},
		{
			name:     "Request for comments in title",
			title:    "Request for comments: API changes",
			labels:   []string{},
			expected: true,
		},
		{
			name:     "RFC label",
			title:    "Add feature",
			labels:   []string{"rfc"},
			expected: true,
		},
		{
			name:     "Proposal label",
			title:    "Add feature",
			labels:   []string{"proposal"},
			expected: true,
		},
		{
			name:     "No RFC indicators",
			title:    "Fix bug",
			labels:   []string{"bug"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock issue
			issue := &github.Issue{
				Title: &tt.title,
			}

			// Add labels
			for _, label := range tt.labels {
				issue.Labels = append(issue.Labels, &github.Label{
					Name: &label,
				})
			}

			result := handler.isRFCIssue(issue)
			if result != tt.expected {
				t.Errorf("isRFCIssue() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
