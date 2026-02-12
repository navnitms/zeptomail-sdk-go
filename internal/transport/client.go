package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthType int

const (
	AuthTypeAPI AuthType = iota
	AuthTypeTemplates
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	oAuthToken string
	authType   AuthType
}

func NewEmailClient(apiKey string, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		apiKey:     apiKey,
		authType:   AuthTypeAPI,
	}
}

func NewTemplatesClient(oAuthToken string, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
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

func (c *Client) Upload(path, filename string, content []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", c.baseURL+path+"?name="+filename, bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
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

	return respBody, nil
}

func (c *Client) Request(method, path string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
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

	return respBody, nil
}
