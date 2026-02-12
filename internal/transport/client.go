package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AuthType int

const (
	AuthTypeAPI AuthType = iota
	AuthTypeTemplates
)

type Response struct {
	StatusCode int
	Body       []byte
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	oAuthToken string
	authType   AuthType
}

func NewEmailClient(apiKey string, baseURL string, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
		authType:   AuthTypeAPI,
	}
}

func NewTemplatesClient(oAuthToken string, baseURL string, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		oAuthToken: oAuthToken,
		authType:   AuthTypeTemplates,
	}
}

func (c *Client) setAuthHeader(req *http.Request) {
	switch c.authType {
	case AuthTypeAPI:
		req.Header.Set("Authorization", "Zoho-enczapikey "+c.apiKey)
	case AuthTypeTemplates:
		req.Header.Set("Authorization", "Zoho-oauthtoken "+c.oAuthToken)
	}
}

func (c *Client) Upload(ctx context.Context, path, filename string, content []byte) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path+"?name="+url.QueryEscape(filename), bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", http.DetectContentType(content))
	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return &Response{StatusCode: resp.StatusCode, Body: respBody}, nil
}

func (c *Client) Request(ctx context.Context, method, path string, payload interface{}) (*Response, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return &Response{StatusCode: resp.StatusCode, Body: respBody}, nil
}
