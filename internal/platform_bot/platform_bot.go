package platform_bot

import (
	"io"
	"log"
	"net/http"

	"github.com/google/go-github/v60/github"
	"github.com/zarishsphere/zs-agent-platform-bot/config"
)

func HandleWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			http.Error(w, "Could not parse webhook payload", http.StatusBadRequest)
			return
		}

		switch e := event.(type) {
		case *github.PullRequestEvent:
			log.Printf("Received PR event: %s (#%d)", *e.Action, *e.PullRequest.Number)
			// TODO: Implement advanced PR validation routing
		case *github.IssuesEvent:
			log.Printf("Received Issue event: %s (#%d)", *e.Action, *e.Issue.Number)
			// TODO: Implement RFC State Machine
		default:
			log.Printf("Unhandled event: %T", e)
		}

		w.WriteHeader(http.StatusOK)
	}
}
