package main

import (
    "log"
    "net/http"
    "os"

    "github.com/yanivmendiuk/W.H.Y/internal/config"
    "github.com/yanivmendiuk/W.H.Y/internal/plainid"
)

func authorizeHandler(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract dynamic values from ArgoCD pre-sync request
        committer := r.Header.Get("ARGOCD_APP_SOURCE_USER") // Example
        path := r.URL.Query().Get("path")                  // e.g., /v1/plainId/tester

        if committer == "" || path == "" {
            http.Error(w, "Missing committer or path", http.StatusBadRequest)
            return
        }

        req := plainid.Request{
            EntityID:     committer,
            EntityTypeID: cfg.EntityTypeID,
            ClientID:     cfg.ClientID,
            ClientSecret: cfg.ClientSecret,
            ListOfResources: []plainid.ResourceGroup{
                {
                    ResourceType: cfg.ResourceType,
                    Resources: []plainid.Resource{
                        {Path: path, Action: "POST"},
                    },
                },
            },
        }

        result, err := plainid.Authorize(cfg.PlainIDEndpoint, req)
        if err != nil {
            log.Printf("Authorization failed: %v", err)
            http.Error(w, "Authorization error", http.StatusInternalServerError)
            return
        }

        if result == "PERMIT" {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("ALLOW"))
        } else {
            w.WriteHeader(http.StatusForbidden)
            w.Write([]byte("DENY"))
        }
    }
}

func main() {
    cfg := config.LoadConfig("application.yaml")

    http.HandleFunc("/authorize", authorizeHandler(cfg))
    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    log.Println("Starting webhook server on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

