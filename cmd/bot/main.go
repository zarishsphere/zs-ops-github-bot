package main

import (
	"log"
	"net/http"

	"github.com/zarishsphere/zs-ops-github-bot/internal/config"
	"github.com/zarishsphere/zs-ops-github-bot/internal/webhook"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	handler := webhook.NewHandler(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", handler.HandleWebhook)
	mux.HandleFunc("/health", healthCheck)

	log.Println("🤖 ZarishSphere GitHub Bot listening on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
