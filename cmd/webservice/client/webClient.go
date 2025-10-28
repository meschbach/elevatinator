package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type webClient struct {
	baseURL       string
	httpClient    *http.Client
	defaultHeader http.Header
}

func (c *webClient) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	fullURL, err := c.resolveURL(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, buf)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Default headers
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Merge custom default headers (do not override explicitly set headers)
	for k, vals := range c.defaultHeader {
		// Only set if not already present
		if req.Header.Get(k) == "" {
			for _, v := range vals {
				req.Header.Add(k, v)
			}
		}
	}

	return req, nil
}

func (c *webClient) resolveURL(path string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}

	base, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse baseURL: %w", err)
	}

	// Normalize path joining
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	base.Path = strings.TrimRight(base.Path, "/") + path

	return base.String(), nil
}

// Do execute the HTTP request and, if v is not nil, decodes the JSON response into v.
// For non-2xx responses it returns an *APIError containing the status and body text.
func (c *webClient) Do(req *http.Request, v any) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer func() { _ = drainAndClose(res.Body) }()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10)) // cap to 64KB
		return errors.New(fmt.Sprintf("unexpected status code: %d, body: %s", res.StatusCode, body))
	}

	if v == nil {
		return nil
	}

	//todo: flag to dump otherwise directly to JSON unmarshall
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, v); err != nil && err != io.EOF {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

const DefaultBaseURL = "http://localhost:8999"

func newWebClient(baseURL string) *webClient {
	return &webClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		defaultHeader: make(http.Header),
	}
}

// drainAndClose ensures the response body is fully read and closed.
func drainAndClose(rc io.ReadCloser) error {
	defer rc.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(rc, 8<<10)) // best-effort drain up to 8KB
	return nil
}
