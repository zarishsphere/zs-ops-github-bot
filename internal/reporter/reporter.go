package reporter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/zarishsphere/zs-ops-github-bot/internal/config"
)

// Reporter handles periodic reporting on repository health
type Reporter struct {
	config *config.Config
}

// NewReporter creates a new reporter
func NewReporter(cfg *config.Config) *Reporter {
	return &Reporter{
		config: cfg,
	}
}

// GenerateWeeklyStalePRReport generates a report of stale PRs and creates an issue
func (r *Reporter) GenerateWeeklyStalePRReport(ctx context.Context, owner, repo string) error {
	log.Printf("Generating weekly stale PR report for %s/%s", owner, repo)

	// Get all open PRs
	prs, _, err := r.config.GitHubClient.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		return fmt.Errorf("failed to list PRs: %w", err)
	}

	// Find stale PRs
	var stalePRs []*github.PullRequest
	for _, pr := range prs {
		if pr.UpdatedAt != nil {
			daysSinceUpdate := time.Since(pr.UpdatedAt.Time).Hours() / 24
			if daysSinceUpdate >= 30 { // 30 days threshold
				stalePRs = append(stalePRs, pr)
			}
		}
	}

	if len(stalePRs) == 0 {
		log.Printf("No stale PRs found for %s/%s", owner, repo)
		return nil
	}

	// Create report
	report := r.formatStalePRReport(stalePRs)

	// Create GitHub issue with report
	issue := &github.IssueRequest{
		Title:  github.String(fmt.Sprintf("🤖 Weekly Stale PR Report (%d PRs)", len(stalePRs))),
		Body:   github.String(report),
		Labels: &[]string{"automated", "maintenance", "stale-prs"},
	}

	_, _, err = r.config.GitHubClient.Issues.Create(ctx, owner, repo, issue)
	if err != nil {
		return fmt.Errorf("failed to create stale PR report issue: %w", err)
	}

	log.Printf("Created stale PR report issue for %s/%s with %d stale PRs", owner, repo, len(stalePRs))
	return nil
}

// formatStalePRReport formats the stale PR report as markdown
func (r *Reporter) formatStalePRReport(stalePRs []*github.PullRequest) string {
	report := `# 🤖 Weekly Stale PR Report

This report identifies pull requests that haven't been updated in 30+ days.

## Stale Pull Requests

| PR | Title | Author | Last Updated | Days Stale |
|----|-------|--------|--------------|------------|
`

	for _, pr := range stalePRs {
		daysStale := int(time.Since(pr.UpdatedAt.Time).Hours() / 24)
		report += fmt.Sprintf("| #%d | %s | @%s | %s | %d days |\n",
			pr.GetNumber(),
			pr.GetTitle(),
			pr.User.GetLogin(),
			pr.UpdatedAt.Format("2006-01-02"),
			daysStale)
	}

	report += `

## Recommendations

- Review and either merge, close, or update these PRs
- Consider adding "stale" labels to old PRs
- Set up automated stale PR management

---
*This report is generated automatically by the ZarishSphere GitHub Bot.*
`

	return report
}
