package porkbun

import (
	"encoding/json"
	"fmt"
)

type apiRequest interface{}

type authRequest struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
	apiRequest
}

func (f authRequest) MarshalJSON() ([]byte, error) {
	type clone authRequest
	cloned := clone(f)

	root, err := json.Marshal(cloned)
	if err != nil {
		return nil, err
	}

	if cloned.apiRequest == nil {
		return root, nil
	}

	embedded, err := json.Marshal(cloned.apiRequest)
	if err != nil {
		return nil, err
	}

	return []byte(string(root[:len(root)-1]) + ",   " + string(embedded[1:])), nil
}

// Status the API response status.
type Status struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// StatusError is a custom error type for easier handling of Porkbun API Errors.
type StatusError struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (a StatusError) Error() string {
	return fmt.Sprintf("status: %s message: %s", a.Status, a.Message)
}

// Record a DNS record.
type Record struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
	TTL     string `json:"ttl,omitempty"`
	Prio    string `json:"prio,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

type pingResponse struct {
	Status
	YourIP string `json:"yourIp"`
}

type createResponse struct {
	Status
	ID int `json:"id"`
}

type retrieveResponse struct {
	Status
	Records []Record `json:"records"`
}
