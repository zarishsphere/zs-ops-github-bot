package main

import (
	"log"
	"net/http"

	"github.com/zarishsphere/zs-agent-platform-bot/config"
	"github.com/zarishsphere/zs-agent-platform-bot/internal/platform_bot"
)

func main() {
	cfg := config.Load()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", platform_bot.HandleWebhook(cfg))

	log.Println("Platform Bot listening on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
