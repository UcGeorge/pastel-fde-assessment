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

const (
	defaultBaseURL    = "https://sigmaprod.sabipay.com/"
	defaultAMLBaseURL = "https://sigmaaml.sabipay.com/"
)

type ClientOption func(*Client)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = client }
}

func WithBaseURL(urlStr string) ClientOption {
	return func(c *Client) { c.baseURL, _ = url.Parse(urlStr) }
}

func WithAMLBaseURL(urlStr string) ClientOption {
	return func(c *Client) { c.amlBaseURL, _ = url.Parse(urlStr) }
}

type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    *url.URL
	amlBaseURL *url.URL
	httpClient *http.Client
}

// New initializes a new Sigma API client
func New(apiKey, apiSecret string, opts ...ClientOption) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	amlBaseURL, _ := url.Parse(defaultAMLBaseURL)

	c := &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		amlBaseURL: amlBaseURL,
		apiSecret:  apiSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// doSigmaRequest makes an authenticated request against the Transaction Monitoring base URL.
func (c *Client) doSigmaRequest(ctx context.Context, method, endpoint string, payload any, responseTarget any) error {
	return c.doRequest(ctx, c.baseURL, method, endpoint, payload, responseTarget)
}

// doAMLRequest makes an authenticated request against the AML base URL (PEP, Sanctions, Adverse Media).
func (c *Client) doAMLRequest(ctx context.Context, method, endpoint string, payload any, responseTarget any) error {
	return c.doRequest(ctx, c.amlBaseURL, method, endpoint, payload, responseTarget)
}

func (c *Client) doRequest(ctx context.Context, base *url.URL, method, endpoint string, payload any, responseTarget any) error {
	var bodyReader io.Reader

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
	fullURL := base.ResolveReference(parsedPath).String()

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
