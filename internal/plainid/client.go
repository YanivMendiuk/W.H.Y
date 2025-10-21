package plainid

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type Resource struct {
    Path   string `json:"path"`
    Action string `json:"action"`
}

type ResourceGroup struct {
    ResourceType string     `json:"resourceType"`
    Resources    []Resource `json:"resources"`
}

type Request struct {
    EntityID       string          `json:"entityId"`
    EntityTypeID   string          `json:"entityTypeId"`
    ClientID       string          `json:"clientId"`
    ClientSecret   string          `json:"clientSecret"`
    ListOfResources []ResourceGroup `json:"listOfResources"`
}

type Response struct {
    Result string `json:"result"`
}

func Authorize(endpoint string, req Request) (string, error) {
    data, _ := json.Marshal(req)
    client := &http.Client{Timeout: 5 * time.Second}
    httpReq, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(httpReq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var r Response
    if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
        return "", err
    }
    return r.Result, nil
}

