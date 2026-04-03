package rfc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/zarishsphere/zs-ops-github-bot/internal/config"
)

// StateMachine manages RFC lifecycle states
type StateMachine struct {
	config *config.Config
}

// RFC states
const (
	StateDraft     = "draft"
	StateReview    = "review"
	StateAccepted  = "accepted"
	StateRejected  = "rejected"
	StateWithdrawn = "withdrawn"
)

// NewStateMachine creates a new RFC state machine
func NewStateMachine(cfg *config.Config) *StateMachine {
	return &StateMachine{
		config: cfg,
	}
}

// ProcessIssueEvent handles RFC-related issue events
func (sm *StateMachine) ProcessIssueEvent(ctx context.Context, owner, repo string, issue *github.Issue) error {
	if issue == nil {
		return fmt.Errorf("invalid issue")
	}

	log.Printf("Processing RFC issue event: %s in %s/%s", issue.GetTitle(), owner, repo)

	// Check current labels and transition state
	currentState := sm.getCurrentState(issue)
	newState := sm.determineNewState(issue)

	if currentState != newState {
		return sm.transitionState(ctx, owner, repo, issue, currentState, newState)
	}

	return nil
}

// ProcessCommentEvent handles RFC-related comment events
func (sm *StateMachine) ProcessCommentEvent(ctx context.Context, owner, repo string, issue *github.Issue, comment *github.IssueComment) error {
	if issue == nil || comment == nil {
		return fmt.Errorf("invalid issue or comment")
	}

	log.Printf("Processing RFC comment event on issue %d in %s/%s", issue.GetNumber(), owner, repo)

	// Check for state transition commands in comments
	body := strings.ToLower(comment.GetBody())

	if strings.Contains(body, "/rfc review") {
		return sm.transitionToReview(ctx, owner, repo, issue)
	} else if strings.Contains(body, "/rfc accept") {
		return sm.transitionToAccepted(ctx, owner, repo, issue)
	} else if strings.Contains(body, "/rfc reject") {
		return sm.transitionToRejected(ctx, owner, repo, issue)
	}

	return nil
}

// getCurrentState determines the current state from issue labels
func (sm *StateMachine) getCurrentState(issue *github.Issue) string {
	for _, label := range issue.Labels {
		labelName := strings.ToLower(label.GetName())
		switch labelName {
		case "rfc-draft":
			return StateDraft
		case "rfc-review":
			return StateReview
		case "rfc-accepted":
			return StateAccepted
		case "rfc-rejected":
			return StateRejected
		case "rfc-withdrawn":
			return StateWithdrawn
		}
	}
	return StateDraft // Default state
}

// determineNewState determines the new state based on current labels
func (sm *StateMachine) determineNewState(issue *github.Issue) string {
	return sm.getCurrentState(issue) // For now, keep current state
}

// transitionState transitions the RFC to a new state
func (sm *StateMachine) transitionState(ctx context.Context, owner, repo string, issue *github.Issue, from, to string) error {
	issueNumber := issue.GetNumber()

	log.Printf("Transitioning RFC %d from %s to %s in %s/%s", issueNumber, from, to, owner, repo)

	// Remove old state label
	oldLabel := fmt.Sprintf("rfc-%s", from)
	if err := sm.removeLabel(ctx, owner, repo, issueNumber, oldLabel); err != nil {
		log.Printf("Failed to remove old label %s: %v", oldLabel, err)
	}

	// Add new state label
	newLabel := fmt.Sprintf("rfc-%s", to)
	if err := sm.addLabel(ctx, owner, repo, issueNumber, newLabel); err != nil {
		return fmt.Errorf("failed to add new label %s: %w", newLabel, err)
	}

	// Add state comment
	comment := sm.generateStateComment(to)
	if err := sm.addComment(ctx, owner, repo, issueNumber, comment); err != nil {
		log.Printf("Failed to add state comment: %v", err)
	}

	return nil
}

// transitionToReview transitions RFC to review state
func (sm *StateMachine) transitionToReview(ctx context.Context, owner, repo string, issue *github.Issue) error {
	return sm.transitionState(ctx, owner, repo, issue, sm.getCurrentState(issue), StateReview)
}

// transitionToAccepted transitions RFC to accepted state
func (sm *StateMachine) transitionToAccepted(ctx context.Context, owner, repo string, issue *github.Issue) error {
	return sm.transitionState(ctx, owner, repo, issue, sm.getCurrentState(issue), StateAccepted)
}

// transitionToRejected transitions RFC to rejected state
func (sm *StateMachine) transitionToRejected(ctx context.Context, owner, repo string, issue *github.Issue) error {
	return sm.transitionState(ctx, owner, repo, issue, sm.getCurrentState(issue), StateRejected)
}

// generateStateComment generates a comment for state transitions
func (sm *StateMachine) generateStateComment(state string) string {
	switch state {
	case StateDraft:
		return `🤖 **RFC State: Draft**

This RFC has been marked as a draft. Community feedback is welcome!

**Next steps:**
- Add detailed proposal in the issue description
- Label with ` + "`rfc-review`" + ` when ready for review
- Follow the [RFC template](https://github.com/zarishsphere/zs-docs-rfc/blob/main/RFC-TEMPLATE.md)`
	case StateReview:
		return `🤖 **RFC State: Under Review**

This RFC is now under active review by the community and maintainers.

**Review checklist:**
- [ ] Technical feasibility assessed
- [ ] Impact on existing systems evaluated
- [ ] Implementation plan discussed
- [ ] Timeline and resources estimated

Use ` + "`rfc-accepted`" + ` or ` + "`rfc-rejected`" + ` labels to complete the review.`
	case StateAccepted:
		return `🎉 **RFC State: Accepted**

This RFC has been accepted! The proposed changes will be implemented.

**Next steps:**
- Create implementation issues
- Assign to appropriate team members
- Schedule for upcoming sprint`
	case StateRejected:
		return `❌ **RFC State: Rejected**

This RFC has been rejected after review.

**Common reasons for rejection:**
- Technical constraints
- Conflicts with existing architecture
- Insufficient community support
- Better alternatives available

You can reopen this RFC if new information becomes available.`
	default:
		return fmt.Sprintf("🤖 RFC state changed to: %s", state)
	}
}

// addLabel adds a label to an issue
func (sm *StateMachine) addLabel(ctx context.Context, owner, repo string, issueNumber int, label string) error {
	_, _, err := sm.config.GitHubClient.Issues.AddLabelsToIssue(ctx, owner, repo, issueNumber, []string{label})
	return err
}

// removeLabel removes a label from an issue
func (sm *StateMachine) removeLabel(ctx context.Context, owner, repo string, issueNumber int, label string) error {
	_, err := sm.config.GitHubClient.Issues.RemoveLabelForIssue(ctx, owner, repo, issueNumber, label)
	return err
}

// addComment adds a comment to an issue
func (sm *StateMachine) addComment(ctx context.Context, owner, repo string, issueNumber int, comment string) error {
	commentReq := &github.IssueComment{
		Body: &comment,
	}
	_, _, err := sm.config.GitHubClient.Issues.CreateComment(ctx, owner, repo, issueNumber, commentReq)
	return err
}
