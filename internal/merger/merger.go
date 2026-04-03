package merger

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v60/github"
	"github.com/zarishsphere/zs-ops-github-bot/internal/config"
)

// Merger handles automatic merging of pull requests
type Merger struct {
	config *config.Config
}

// NewMerger creates a new merger
func NewMerger(cfg *config.Config) *Merger {
	return &Merger{
		config: cfg,
	}
}

// CheckAutoMerge checks if a pull request is eligible for auto-merge
func (m *Merger) CheckAutoMerge(ctx context.Context, pr *github.PullRequest) error {
	if pr == nil || pr.Base == nil || pr.Base.Repo == nil {
		return fmt.Errorf("invalid pull request")
	}

	if !m.config.AutoMergeEnabled {
		return nil // Auto-merge disabled
	}

	owner := pr.Base.Repo.GetOwner().GetLogin()
	repo := pr.Base.Repo.GetName()
	prNumber := pr.GetNumber()

	log.Printf("Checking auto-merge eligibility for PR %d in %s/%s", prNumber, owner, repo)

	// Check if PR is from trusted bot
	if !m.isTrustedBot(pr) {
		log.Printf("PR %d not from trusted bot, skipping auto-merge", prNumber)
		return nil
	}

	// Check CI status
	status, err := m.checkCIStatus(ctx, owner, repo, pr)
	if err != nil {
		return fmt.Errorf("failed to check CI status: %w", err)
	}

	if !status {
		log.Printf("CI not passing for PR %d, skipping auto-merge", prNumber)
		return nil
	}

	// Check review requirements
	reviews, err := m.checkReviews(ctx, owner, repo, pr)
	if err != nil {
		return fmt.Errorf("failed to check reviews: %w", err)
	}

	if !reviews {
		log.Printf("Review requirements not met for PR %d, skipping auto-merge", prNumber)
		return nil
	}

	// Perform auto-merge
	return m.performAutoMerge(ctx, owner, repo, pr)
}

// isTrustedBot checks if the PR author is a trusted bot
func (m *Merger) isTrustedBot(pr *github.PullRequest) bool {
	author := pr.GetUser().GetLogin()
	trustedBots := []string{
		"dependabot",
		"renovate",
		"dependabot[bot]",
		"renovate[bot]",
	}

	for _, bot := range trustedBots {
		if author == bot {
			return true
		}
	}

	return false
}

// checkCIStatus checks if all required CI checks are passing
func (m *Merger) checkCIStatus(ctx context.Context, owner, repo string, pr *github.PullRequest) (bool, error) {
	// Get combined status for the PR head commit
	combined, _, err := m.config.GitHubClient.Repositories.GetCombinedStatus(ctx, owner, repo, pr.Head.GetSHA(), nil)
	if err != nil {
		return false, err
	}

	// Check if overall status is success
	return combined.GetState() == "success", nil
}

// checkReviews checks if the PR has required reviews
func (m *Merger) checkReviews(ctx context.Context, owner, repo string, pr *github.PullRequest) (bool, error) {
	reviews, _, err := m.config.GitHubClient.PullRequests.ListReviews(ctx, owner, repo, pr.GetNumber(), nil)
	if err != nil {
		return false, err
	}

	approved := 0
	for _, review := range reviews {
		if review.GetState() == "APPROVED" {
			approved++
		}
	}

	// Require at least 1 approval for bot PRs
	return approved >= 1, nil
}

// performAutoMerge performs the actual merge operation
func (m *Merger) performAutoMerge(ctx context.Context, owner, repo string, pr *github.PullRequest) error {
	prNumber := pr.GetNumber()
	sha := pr.Head.GetSHA()

	log.Printf("Performing auto-merge for PR %d", prNumber)

	// Merge with squash strategy
	opts := &github.PullRequestOptions{
		MergeMethod: "squash",
		SHA:         sha,
	}

	result, _, err := m.config.GitHubClient.PullRequests.Merge(ctx, owner, repo, prNumber, "", opts)
	if err != nil {
		return fmt.Errorf("failed to merge PR: %w", err)
	}

	if !result.GetMerged() {
		return fmt.Errorf("PR merge failed: %s", result.GetMessage())
	}

	log.Printf("Successfully auto-merged PR %d", prNumber)
	return nil
}
