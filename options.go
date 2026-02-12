package zeptomail

import "net/http"

// Option tweaks client behaviour. Pass to NewEmailClient or NewTemplatesClient.
type Option func(*clientConfig)

type clientConfig struct {
	httpClient *http.Client
	baseURL    string
}

// WithHTTPClient replaces the default http.Client (which has a 30 s timeout).
func WithHTTPClient(c *http.Client) Option {
	return func(cfg *clientConfig) {
		cfg.httpClient = c
	}
}

// WithBaseURL overrides the default API base URL. Mainly useful for tests.
func WithBaseURL(url string) Option {
	return func(cfg *clientConfig) {
		cfg.baseURL = url
	}
}
