package labeler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/zarishsphere/zs-ops-github-bot/internal/config"
)

// Labeler handles automatic labeling of pull requests
type Labeler struct {
	config *config.Config
	rules  map[string][]string
}

// NewLabeler creates a new labeler
func NewLabeler(cfg *config.Config) *Labeler {
	return &Labeler{
		config: cfg,
		rules: map[string][]string{
			"zs-core-":     {"core", "backend"},
			"zs-svc-":      {"service", "backend"},
			"zs-ui-":       {"frontend", "ui"},
			"zs-mobile-":   {"mobile", "flutter"},
			"zs-desktop-":  {"desktop", "tauri"},
			"zs-iac-":      {"infrastructure", "terraform", "kubernetes"},
			"zs-data-":     {"data", "terminology", "content"},
			"zs-content-":  {"content", "forms", "clinical"},
			"zs-ops-":      {"operations", "automation"},
			"zs-agent-":    {"automation", "ai"},
			"zs-distro-":   {"distribution", "deployment"},
			"zs-docs-":     {"documentation"},
			"zs-int-":      {"integration", "interoperability"},
			".github/":     {"github", "automation"},
			"*.md":         {"documentation"},
			"*.yml":        {"ci-cd", "automation"},
			"*.yaml":       {"ci-cd", "automation"},
			"go.mod":       {"golang", "dependencies"},
			"package.json": {"javascript", "dependencies"},
			"Dockerfile":   {"docker", "containerization"},
		},
	}
}

// LabelPullRequest automatically labels a pull request based on changed files
func (l *Labeler) LabelPullRequest(ctx context.Context, pr *github.PullRequest) error {
	if pr == nil || pr.Base == nil || pr.Base.Repo == nil {
		return fmt.Errorf("invalid pull request")
	}

	owner := pr.Base.Repo.GetOwner().GetLogin()
	repo := pr.Base.Repo.GetName()
	prNumber := pr.GetNumber()

	log.Printf("Labeling PR %d in %s/%s", prNumber, owner, repo)

	// Get list of changed files
	files, _, err := l.config.GitHubClient.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return fmt.Errorf("failed to get PR files: %w", err)
	}

	// Determine labels based on file patterns
	labels := l.determineLabels(files)

	if len(labels) == 0 {
		log.Printf("No labels to apply for PR %d", prNumber)
		return nil
	}

	// Apply labels
	_, _, err = l.config.GitHubClient.Issues.AddLabelsToIssue(ctx, owner, repo, prNumber, labels)
	if err != nil {
		return fmt.Errorf("failed to add labels: %w", err)
	}

	log.Printf("Applied labels %v to PR %d", labels, prNumber)
	return nil
}

// determineLabels determines which labels to apply based on changed files
func (l *Labeler) determineLabels(files []*github.CommitFile) []string {
	labelSet := make(map[string]bool)

	for _, file := range files {
		filename := file.GetFilename()

		// Check file path patterns
		for pattern, labels := range l.rules {
			if l.matchesPattern(filename, pattern) {
				for _, label := range labels {
					labelSet[label] = true
				}
			}
		}
	}

	// Convert set to slice
	var labels []string
	for label := range labelSet {
		labels = append(labels, label)
	}

	return labels
}

// matchesPattern checks if a filename matches a pattern
func (l *Labeler) matchesPattern(filename, pattern string) bool {
	if strings.HasPrefix(pattern, "*.") {
		ext := strings.TrimPrefix(pattern, "*.")
		return strings.HasSuffix(filename, "."+ext)
	}
	return strings.Contains(filename, pattern)
}
