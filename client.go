package zeptomail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/navnitms/zeptomail-sdk-go/internal/transport"
)

const baseURL = "https://api.zeptomail.com/v1.1"
const filesEndpoint = "/files"

// EmailClient talks to the ZeptoMail transactional email endpoints.
type EmailClient struct {
	httpClient *transport.Client
}

// TemplatesClient talks to the ZeptoMail template CRUD endpoints.
type TemplatesClient struct {
	httpClient *transport.Client
}

func defaultHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

// NewEmailClient returns a client that authenticates with the given API key.
func NewEmailClient(apiKey string, opts ...Option) *EmailClient {
	cfg := &clientConfig{
		httpClient: defaultHTTPClient(),
		baseURL:    baseURL,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &EmailClient{
		httpClient: transport.NewEmailClient(apiKey, cfg.baseURL, cfg.httpClient),
	}
}

// NewTemplatesClient returns a client that authenticates with the given OAuth token.
func NewTemplatesClient(oAuthToken string, opts ...Option) *TemplatesClient {
	cfg := &clientConfig{
		httpClient: defaultHTTPClient(),
		baseURL:    baseURL,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &TemplatesClient{
		httpClient: transport.NewTemplatesClient(oAuthToken, cfg.baseURL, cfg.httpClient),
	}
}

const maxErrorBodyLen = 512

func checkForError(resp *transport.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	apiErr := &APIError{HTTPStatusCode: resp.StatusCode}

	var errorResp ErrorResponse
	if err := json.Unmarshal(resp.Body, &errorResp); err == nil && errorResp.Error.Code != "" {
		apiErr.Code = errorResp.Error.Code
		apiErr.Message = errorResp.Error.Message
		apiErr.Details = errorResp.Error.Details
		apiErr.RequestID = errorResp.Error.RequestID
		return apiErr
	}

	// Non-JSON or unexpected error body
	body := string(resp.Body)
	if len(body) > maxErrorBodyLen {
		body = body[:maxErrorBodyLen]
	}
	apiErr.Code = http.StatusText(resp.StatusCode)
	apiErr.Message = body
	return apiErr
}

func (ec *EmailClient) handleResponse(resp *transport.Response) (*SuccessResponse, error) {
	if err := checkForError(resp); err != nil {
		return nil, err
	}

	var successResp SuccessResponse
	if err := json.Unmarshal(resp.Body, &successResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &successResp, nil
}

func (ec *EmailClient) SendEmail(ctx context.Context, req *EmailRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request(ctx, "POST", "/email", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) SendBatchEmail(ctx context.Context, req *EmailRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request(ctx, "POST", "/email/batch", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) SendTemplateEmail(ctx context.Context, req *TemplateRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request(ctx, "POST", "/email/template", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) SendBatchTemplateEmail(ctx context.Context, req *TemplateRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request(ctx, "POST", "/email/template/batch", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) FileCacheUpload(ctx context.Context, filename string, content []byte) (*FileUploadResponse, error) {
	resp, err := ec.httpClient.Upload(ctx, filesEndpoint, filename, content)
	if err != nil {
		return nil, err
	}

	if err := checkForError(resp); err != nil {
		return nil, err
	}

	var fileResp FileUploadResponse
	if err := json.Unmarshal(resp.Body, &fileResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &fileResp, nil
}

func (tc *TemplatesClient) CreateTemplate(ctx context.Context, mailagentAlias string, req *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates", url.PathEscape(mailagentAlias))
	resp, err := tc.httpClient.Request(ctx, "POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	if err := checkForError(resp); err != nil {
		return nil, err
	}

	var templateResp CreateTemplateResponse
	if err := json.Unmarshal(resp.Body, &templateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &templateResp, nil
}

func (tc *TemplatesClient) GetTemplate(ctx context.Context, mailagentAlias, templateKey string) (*GetTemplateResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/%s", url.PathEscape(mailagentAlias), url.PathEscape(templateKey))
	resp, err := tc.httpClient.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	if err := checkForError(resp); err != nil {
		return nil, err
	}

	var templateResp GetTemplateResponse
	if err := json.Unmarshal(resp.Body, &templateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &templateResp, nil
}

func (tc *TemplatesClient) UpdateTemplate(ctx context.Context, mailagentAlias, templateKey string, req *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/%s", url.PathEscape(mailagentAlias), url.PathEscape(templateKey))
	resp, err := tc.httpClient.Request(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	if err := checkForError(resp); err != nil {
		return nil, err
	}

	var templateResp CreateTemplateResponse
	if err := json.Unmarshal(resp.Body, &templateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &templateResp, nil
}

func (tc *TemplatesClient) ListTemplates(ctx context.Context, mailagentAlias string, params ListTemplatesParams) (*ListTemplatesResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/?offset=%d&limit=%d",
		url.PathEscape(mailagentAlias), params.Offset, params.Limit)

	resp, err := tc.httpClient.Request(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	if err := checkForError(resp); err != nil {
		return nil, err
	}

	var listResp ListTemplatesResponse
	if err := json.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &listResp, nil
}

func (tc *TemplatesClient) DeleteTemplate(ctx context.Context, mailagentAlias, templateKey string) error {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/%s", url.PathEscape(mailagentAlias), url.PathEscape(templateKey))
	resp, err := tc.httpClient.Request(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return checkForError(resp)
}
