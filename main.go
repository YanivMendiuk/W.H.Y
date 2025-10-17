package main

import (
    "log"
    "net/http"

    "github.com/yanivmendiuk/W.H.Y/internal/webhook"
)

func main() {
    http.HandleFunc("/authorize", webhook.AuthorizeHandler)
    http.HandleFunc("/healthz", webhook.HealthHandler)

    log.Println("Starting webhook server on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %s", err)
    }
}

