package plainid

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

// Resource represents an individual resource in a resource group
type Resource struct {
    Path   string `json:"path"`
    Action string `json:"action"`
}

// ResourceGroup groups multiple resources under a type
type ResourceGroup struct {
    ResourceType string     `json:"resourceType"`
    Resources    []Resource `json:"resources"`
}

// Request represents the payload sent to PlainID
type Request struct {
    EntityID       string          `json:"entityId"`
    EntityTypeID   string          `json:"entityTypeId"`
    ClientID       string          `json:"clientId"`
    ClientSecret   string          `json:"clientSecret"`
    ListOfResources []ResourceGroup `json:"listOfResources"`
}

// Response represents the PlainID authorization response
type Response struct {
    Result string `json:"result"`
}

// Authorize sends the request to PlainID and returns the decision
func Authorize(endpoint string, req Request) (string, error) {
    data, err := json.Marshal(req)
    if err != nil {
        log.Printf("Error marshalling request: %v", err)
        return "", err
    }

    log.Printf("Sending authorization request to PlainID endpoint: %s", endpoint)
    log.Printf("Request payload: %s", string(data))

    client := &http.Client{Timeout: 5 * time.Second}
    httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
    if err != nil {
        log.Printf("Error creating HTTP request: %v", err)
        return "", err
    }
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(httpReq)
    if err != nil {
        log.Printf("Error making HTTP request: %v", err)
        return "", err
    }
    defer resp.Body.Close()

    var r Response
    if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
        log.Printf("Error decoding response: %v", err)
        return "", err
    }

    log.Printf("Received response from PlainID: %s", r.Result)
    return r.Result, nil
}

