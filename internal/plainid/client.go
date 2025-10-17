package plainid

import (
    "bytes"
    "encoding/json"
    "net/http"
)

// Request structs
type PlainIDRequest struct {
    EntityId        string          `json:"entityId"`
    EntityTypeId    string          `json:"entityTypeId"`
    ClientId        string          `json:"clientId"`
    ClientSecret    string          `json:"clientSecret"`
    ListOfResources []ResourceGroup `json:"listOfResources"`
}

type ResourceGroup struct {
    ResourceType string     `json:"resourceType"`
    Resources    []Resource `json:"resources"`
}

type Resource struct {
    Path   string `json:"path"`
    Action string `json:"action"`
}

// Response struct
type PlainIDResponse struct {
    Result string `json:"result"` // PERMIT or DENY
}

// CallPlainID sends the request to PlainID Runtime and returns the result
func CallPlainID(reqData PlainIDRequest) (string, error) {
    url := "https://a93c9e08076354cdbaf7e0ffe6ca8ef5-1536407533.us-east-1.elb.amazonaws.com/api/runtime/permit-deny/v3"

    jsonData, err := json.Marshal(reqData)
    if err != nil {
        return "", err
    }

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result PlainIDResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    return result.Result, nil
}

