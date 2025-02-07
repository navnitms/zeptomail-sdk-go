package zeptomail

import (
	"encoding/json"
	"fmt"
	"navnitms/zeptomail-sdk-go/internal/http"
)

const baseURL = "https://api.zeptomail.in/v1.1"
const filesEndpoint = "/files"

type EmailClient struct {
	httpClient *http.Client
}

type TemplatesClient struct {
	httpClient *http.Client
}

func NewEmailClient(apiKey string) *EmailClient {
	return &EmailClient{
		httpClient: http.NewEmailClient(apiKey, baseURL),
	}
}

func NewTemplatesClient(oAuthToken string) *TemplatesClient {
	return &TemplatesClient{
		httpClient: http.NewTemplatesClient(oAuthToken, baseURL),
	}
}

func (ec *EmailClient) handleResponse(resp []byte) (*SuccessResponse, error) {
	var successResp SuccessResponse
	if err := json.Unmarshal(resp, &successResp); err == nil {
		return &successResp, nil
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(resp, &errorResp); err == nil {
		return nil, fmt.Errorf("API error: %s - %s", errorResp.Error.Code, errorResp.Error.Message)
	}

	return nil, fmt.Errorf("failed to parse response: %s", string(resp))
}

func (ec *EmailClient) SendEmail(req *EmailRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request("POST", "/email", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) SendBatchEmail(req *EmailRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request("POST", "/email/batch", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) SendTemplateEmail(req *TemplateRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request("POST", "/email/template", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) SendBatchTemplateEmail(req *TemplateRequest) (*SuccessResponse, error) {
	resp, err := ec.httpClient.Request("POST", "/email/template/batch", req)
	if err != nil {
		return nil, err
	}
	return ec.handleResponse(resp)
}

func (ec *EmailClient) FileCacheUpload(filename string, content []byte) (*FileUploadResponse, error) {
	resp, err := ec.httpClient.Upload(filesEndpoint, filename,	 content)
	if err != nil {
		return nil, err
	}

	var fileResp FileUploadResponse
	if err := json.Unmarshal(resp, &fileResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &fileResp, nil
}

func (tc *TemplatesClient) CreateTemplate(mailagentAlias string, req *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates", mailagentAlias)
	resp, err := tc.httpClient.Request("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	var templateResp CreateTemplateResponse
	if err := json.Unmarshal(resp, &templateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &templateResp, nil
}

func (tc *TemplatesClient) GetTemplate(mailagentAlias, templateKey string) (*GetTemplateResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/%s", mailagentAlias, templateKey)
	resp, err := tc.httpClient.Request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var templateResp GetTemplateResponse
	if err := json.Unmarshal(resp, &templateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &templateResp, nil
}

func (tc *TemplatesClient) UpdateTemplate(mailagentAlias, templateKey string, req *UpdateTemplateRequest) (*CreateTemplateResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/%s", mailagentAlias, templateKey)
	resp, err := tc.httpClient.Request("PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	var templateResp CreateTemplateResponse
	if err := json.Unmarshal(resp, &templateResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &templateResp, nil
}

func (tc *TemplatesClient) ListTemplates(mailagentAlias string, params ListTemplatesParams) (*ListTemplatesResponse, error) {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/?offset=%d&limit=%d",
		mailagentAlias, params.Offset, params.Limit)

	resp, err := tc.httpClient.Request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var listResp ListTemplatesResponse
	if err := json.Unmarshal(resp, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &listResp, nil
}

func (tc *TemplatesClient) DeleteTemplate(mailagentAlias, templateKey string) error {
	endpoint := fmt.Sprintf("/mailagents/%s/templates/%s", mailagentAlias, templateKey)
	_, err := tc.httpClient.Request("DELETE", endpoint, nil)
	return err
}
