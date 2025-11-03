package plainid

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
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
    EntityID               string          `json:"entityId"`
    EntityTypeID           string          `json:"entityTypeId"`
    ClientID               string          `json:"clientId"`
    ClientSecret           string          `json:"clientSecret"`
    ListOfResources        []ResourceGroup `json:"listOfResources"`
    UseCache               bool            `json:"useCache,omitempty"`
    IncludeAccessPolicy    bool            `json:"includeAccessPolicy,omitempty"`
    IncludeIdentity        bool            `json:"includeIdentity,omitempty"`
    IncludeAssetAttributes bool            `json:"includeAssetAttributes,omitempty"`
    IncludeDenyReason      bool            `json:"includeDenyReason,omitempty"`
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

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return "", err
    }

    // Log the response body regardless of status code
    if resp.StatusCode != http.StatusOK {
        log.Printf("Authorization failed with status %d: %s", resp.StatusCode, string(bodyBytes))
        return "", err
    }

    log.Printf("Authorization succeeded with status 200. Response body: %s", string(bodyBytes))

    var r Response
    if err := json.Unmarshal(bodyBytes, &r); err != nil {
        log.Printf("Error decoding response: %v", err)
        return "", err
    }

    log.Printf("Parsed PlainID result: %s", r.Result)
    return r.Result, nil
}
