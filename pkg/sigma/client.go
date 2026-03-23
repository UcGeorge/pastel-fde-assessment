package sigma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://sigmaprod.sabipay.com/"

type ClientOption func(*Client)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = client }
}

func WithBaseURL(urlStr string) ClientOption {
	return func(c *Client) { c.baseURL, _ = url.Parse(urlStr) }
}

type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    *url.URL
	httpClient *http.Client
}

// New initializes a new Sigma API client
func New(apiKey, apiSecret string, opts ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		apiSecret:  apiSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) doSigmaRequest(ctx context.Context, method, endpoint string, payload any, responseTarget any) error {
	var bodyReader io.Reader
	var err error

	// Prepare the payload
	if payload != nil {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	// Resolve the request URL
	parsedPath, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	fullURL := c.baseURL.ResolveReference(parsedPath).String()

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Required Headers
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header["apiKey"] = []string{c.apiKey}
	req.Header["apiSecret"] = []string{c.apiSecret}

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error during API call: %w", err)
	}

	// Discard the remaining body bytes before closing.
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	// Handle Edge Cases (Non-2xx Status Codes)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &SigmaAPIError{
			StatusCode: resp.StatusCode,
		}

		// Try to decode the Sigma API's JSON error response
		// If decoding fails, we just keep the status code and add a generic message
		if err := json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			apiErr.Message = "failed to parse error response from sigma"
		}

		// Map the HTTP status to our Sentinel Errors for Unwrap()
		switch resp.StatusCode {
		case http.StatusBadRequest:
			apiErr.err = ErrBadRequest
		case http.StatusUnauthorized:
			apiErr.err = ErrUnauthorized
		case http.StatusNotFound:
			apiErr.err = ErrResourceNotFound
		case http.StatusTooManyRequests:
			apiErr.err = ErrRateLimited
		case http.StatusInternalServerError, http.StatusBadGateway:
			apiErr.err = ErrInternalServer
		default:
			apiErr.err = ErrUnknown
		}

		// Return the beautifully wrapped custom error
		return apiErr
	}

	// Decode the successful response
	if responseTarget != nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(responseTarget); err != nil {
			return fmt.Errorf("failed to decode API response: %w", err)
		}
	}

	return nil
}
