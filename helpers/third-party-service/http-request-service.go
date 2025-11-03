package thirdpartyservice

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	HTTPRequester interface {
		GetURL() string
		GetHTTPRequestType() string
		GetBodyForRequest() *bytes.Buffer
		GetParameters() map[string]string
		NeededResponse() bool
		SetBodyResponse([]byte) (err error)
		AddAuthorization(req *http.Request)
	}
)

func ExecuteHttpRequest(requester HTTPRequester) error {
	// Build the HTTP request based on interface methods
	url := requester.GetURL()
	requestType := requester.GetHTTPRequestType()
	body := requester.GetBodyForRequest()
	params := requester.GetParameters()

	var req *http.Request
	var err error

	switch requestType {
	case http.MethodGet:
		req, err = http.NewRequest(requestType, url, nil)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if body == nil {
			body = bytes.NewBuffer([]byte{})
		}
		req, err = http.NewRequest(requestType, url, body)
	default:
		return errors.New("unsupported HTTP request type")
	}

	if err != nil {
		return err
	}

	// Add parameters to the request URL (optional)
	if params != nil {
		query := req.URL.Query()
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	requester.AddAuthorization(req)

	// Execute the request with the default client
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Set body response if needed
	if requester.NeededResponse() {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		if err := requester.SetBodyResponse(responseBody); err != nil {
			return fmt.Errorf("error setting response body: %w", err)
		}
	}

	return err
}
