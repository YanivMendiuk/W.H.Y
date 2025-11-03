package main

import (
    "bytes"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/yanivmendiuk/W.H.Y/internal/config"
    "github.com/yanivmendiuk/W.H.Y/internal/plainid"
)

// authorizeHandler handles authorization requests
func authorizeHandler(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Log the incoming request body
        bodyBytes, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Failed to read request body: %v", err)
            http.Error(w, "Failed to read request", http.StatusBadRequest)
            return
        }
        log.Printf("Incoming request body: %s", string(bodyBytes))
        // Reset the request body so it can be read again if needed
        r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

        // Extract dynamic values from ArgoCD pre-sync request
        committer := r.Header.Get("ARGOCD_APP_SOURCE_USER")
        if committer == "" {
            committer = r.Header.Get("Argocd-App-Source-User")
        }

        path := r.URL.Query().Get("path")

        if committer == "" || path == "" {
            http.Error(w, "Missing committer or path", http.StatusBadRequest)
            return
        }

        req := plainid.Request{
            EntityID:               committer,
            EntityTypeID:           cfg.EntityTypeID,
            ClientID:               os.Getenv("CLIENT_ID"),
            ClientSecret:           os.Getenv("CLIENT_SECRET"),
            ListOfResources: []plainid.ResourceGroup{
                {
                    ResourceType: cfg.ResourceType,
                    Resources: []plainid.Resource{
                        {Path: path, Action: cfg.Action}, // dynamic action from config
                    },
                },
            },
            UseCache:               false,
            IncludeAccessPolicy:    true,
            IncludeIdentity:        true,
            IncludeAssetAttributes: true,
            IncludeDenyReason:      true,
        }

        // Call PlainID Runtime
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
    // Determine config path
    cfgPath := os.Getenv("CONFIG_PATH")
    if cfgPath == "" {
        cwd, err := os.Getwd()
        if err != nil {
            log.Fatalf("Failed to get current directory: %v", err)
        }
        cfgPath = filepath.Join(cwd, "config", "application.yaml")
    }

    cfg, err := config.LoadConfig(cfgPath)
    if err != nil {
        log.Fatalf("Failed to read config: %v", err)
    }

    // Allow overriding PlainIDEndpoint via env var (useful in K8s)
    if envEndpoint := os.Getenv("PLAINID_ENDPOINT"); envEndpoint != "" {
        cfg.PlainIDEndpoint = envEndpoint
    }

    log.Printf("Using PlainID Runtime endpoint: %s", cfg.PlainIDEndpoint)

    http.HandleFunc("/authorize", authorizeHandler(cfg))
    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    log.Println("Starting webhook server on :8181...")
    if err := http.ListenAndServe(":8181", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
