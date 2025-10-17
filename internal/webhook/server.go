package webhook

import (
    "encoding/json"
    "net/http"

    "github.com/yanivmendiuk/W.H.Y/internal/plainid"
)

// AuthorizeHandler handles the webhook requests from ArgoCD
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: replace with actual ArgoCD payload parsing
    req := plainid.PlainIDRequest{
        EntityId:     "ABC",
        EntityTypeId: "Leumi_Array_Identity",
        ClientId:     "PJMBO9FCTONFJ8EDYWT7",
        ClientSecret: "JZ6gerPR3SWCugXOmHhtsuMZZnkkclS26o43XSDN",
        ListOfResources: []plainid.ResourceGroup{
            {
                ResourceType: "EP_ROLE_TEST",
                Resources: []plainid.Resource{
                    {Path: "/v1/plainId/tester", Action: "POST"},
                },
            },
        },
    }

    result, err := plainid.CallPlainID(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(plainid.PlainIDResponse{Result: result})
}

// HealthHandler simple health endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}

