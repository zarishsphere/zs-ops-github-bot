package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/zarishsphere/zs-ops-github-bot/internal/config"
	"github.com/zarishsphere/zs-ops-github-bot/internal/labeler"
	"github.com/zarishsphere/zs-ops-github-bot/internal/merger"
	"github.com/zarishsphere/zs-ops-github-bot/internal/reporter"
	"github.com/zarishsphere/zs-ops-github-bot/internal/rfc"
)

// Handler handles GitHub webhook events
type Handler struct {
	config   *config.Config
	labeler  *labeler.Labeler
	merger   *merger.Merger
	rfc      *rfc.StateMachine
	reporter *reporter.Reporter
}

// NewHandler creates a new webhook handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config:   cfg,
		labeler:  labeler.NewLabeler(cfg),
		merger:   merger.NewMerger(cfg),
		rfc:      rfc.NewStateMachine(cfg),
		reporter: reporter.NewReporter(cfg),
	}
}

// HandleWebhook processes incoming GitHub webhook events
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Verify webhook signature if secret is configured
	if h.config.WebhookSecret != "" {
		if !h.verifySignature(r, payload) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Parse webhook payload
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("Could not parse webhook payload: %v", err)
		http.Error(w, "Could not parse webhook payload", http.StatusBadRequest)
		return
	}

	// Route event to appropriate handler
	switch e := event.(type) {
	case *github.PullRequestEvent:
		h.handlePullRequest(r.Context(), e)
	case *github.IssuesEvent:
		h.handleIssue(r.Context(), e)
	case *github.PullRequestReviewEvent:
		h.handlePullRequestReview(r.Context(), e)
	default:
		log.Printf("Unhandled event type: %T", e)
	}

	w.WriteHeader(http.StatusOK)
}

// verifySignature verifies the GitHub webhook signature
func (h *Handler) verifySignature(r *http.Request, payload []byte) bool {
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return false
	}

	// Remove "sha256=" prefix
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	signature = strings.TrimPrefix(signature, "sha256=")

	// Compute expected signature
	mac := hmac.New(sha256.New, []byte(h.config.WebhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// handlePullRequest handles pull request events
func (h *Handler) handlePullRequest(ctx context.Context, event *github.PullRequestEvent) {
	log.Printf("Received PR event: %s (#%d)", event.GetAction(), event.PullRequest.GetNumber())

	switch event.GetAction() {
	case "opened", "reopened", "synchronize":
		// Auto-label PR based on changed files
		if err := h.labeler.LabelPullRequest(ctx, event.PullRequest); err != nil {
			log.Printf("Error labeling PR: %v", err)
		}

		// Check for auto-merge eligibility
		if err := h.merger.CheckAutoMerge(ctx, event.PullRequest); err != nil {
			log.Printf("Error checking auto-merge: %v", err)
		}
	case "closed":
		if event.PullRequest.GetMerged() {
			log.Printf("PR #%d merged successfully", event.PullRequest.GetNumber())
		}
	}
}

// handleIssue handles issue events
func (h *Handler) handleIssue(ctx context.Context, event *github.IssuesEvent) {
	log.Printf("Received Issue event: %s (#%d)", event.GetAction(), event.Issue.GetNumber())

	// Handle RFC state machine for issues labeled as RFC
	if h.isRFCIssue(event.Issue) {
		if err := h.rfc.ProcessIssueEvent(ctx, event.Repo.Owner.GetLogin(), event.Repo.GetName(), event.Issue); err != nil {
			log.Printf("Error processing RFC event: %v", err)
		}
	}
}

// handlePullRequestReview handles pull request review events
func (h *Handler) handlePullRequestReview(ctx context.Context, event *github.PullRequestReviewEvent) {
	log.Printf("Received PR review event: %s on PR #%d", event.GetAction(), event.PullRequest.GetNumber())

	// Re-check auto-merge after review
	if event.GetAction() == "submitted" {
		prEvent := &github.PullRequestEvent{
			Action:      github.String("synchronize"),
			PullRequest: event.PullRequest,
			Repo:        event.Repo,
		}

		if err := h.merger.CheckAutoMerge(ctx, prEvent.PullRequest); err != nil {
			log.Printf("Error re-checking auto-merge after review: %v", err)
		}
	}
}

// isRFCIssue checks if an issue is an RFC based on labels or title
func (h *Handler) isRFCIssue(issue *github.Issue) bool {
	// Check labels first
	for _, label := range issue.Labels {
		labelName := strings.ToLower(label.GetName())
		if labelName == "rfc" || labelName == "proposal" {
			return true
		}
	}

	// Check title
	title := strings.ToLower(issue.GetTitle())
	return strings.Contains(title, "rfc") || strings.Contains(title, "request for comments")
}
