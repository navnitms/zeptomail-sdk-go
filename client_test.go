package zeptomail

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func successJSON() string {
	return `{"data":[{"code":"SUCCESS","message":"OK","additional_info":[]}],"message":"OK","request_id":"req-1","object":"email"}`
}

func errorJSON() string {
	return `{"error":{"code":"INVALID_DATA","message":"bad request","details":[{"code":"REQUIRED","message":"to is required","target":"to"}],"request_id":"req-err"}}`
}

func newTestEmailClient(url string) *EmailClient {
	return NewEmailClient("test-api-key", WithBaseURL(url))
}

func newTestTemplatesClient(url string) *TemplatesClient {
	return NewTemplatesClient("test-oauth-token", WithBaseURL(url))
}

// --- SendEmail ---

func TestSendEmail_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	resp, err := client.SendEmail(context.Background(), &EmailRequest{
		From:    EmailAddress{Address: "a@b.com"},
		To:      []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
		Subject: "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.RequestID != "req-1" {
		t.Errorf("RequestID = %q, want %q", resp.RequestID, "req-1")
	}
}

func TestSendEmail_RequestBody(t *testing.T) {
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	track := true
	_, err := client.SendEmail(context.Background(), &EmailRequest{
		From:        EmailAddress{Address: "a@b.com", Name: "Sender"},
		To:          []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
		Subject:     "Test",
		HTMLBody:    "<b>hi</b>",
		TrackClicks: &track,
	})
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(gotBody, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["subject"] != "Test" {
		t.Errorf("subject = %v, want %q", parsed["subject"], "Test")
	}
	if parsed["track_clicks"] != true {
		t.Errorf("track_clicks = %v, want true", parsed["track_clicks"])
	}
}

func TestSendEmail_AuthHeader(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendEmail(context.Background(), &EmailRequest{
		From: EmailAddress{Address: "a@b.com"},
		To:   []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if gotAuth != "Zoho-enczapikey test-api-key" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Zoho-enczapikey test-api-key")
	}
}

func TestSendEmail_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(errorJSON()))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendEmail(context.Background(), &EmailRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.HTTPStatusCode != 400 {
		t.Errorf("HTTPStatusCode = %d, want 400", apiErr.HTTPStatusCode)
	}
	if apiErr.Code != "INVALID_DATA" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "INVALID_DATA")
	}
	if apiErr.RequestID != "req-err" {
		t.Errorf("RequestID = %q, want %q", apiErr.RequestID, "req-err")
	}
	if len(apiErr.Details) != 1 {
		t.Fatalf("Details length = %d, want 1", len(apiErr.Details))
	}
}

func TestSendEmail_NonJSONError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(502)
		w.Write([]byte("<html>Bad Gateway</html>"))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendEmail(context.Background(), &EmailRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.HTTPStatusCode != 502 {
		t.Errorf("HTTPStatusCode = %d, want 502", apiErr.HTTPStatusCode)
	}
	if apiErr.Code != "Bad Gateway" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "Bad Gateway")
	}
	if !strings.Contains(apiErr.Message, "<html>") {
		t.Errorf("Message = %q, want to contain HTML body", apiErr.Message)
	}
}

func TestSendEmail_NetworkError(t *testing.T) {
	// Use a closed server to trigger a network error.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendEmail(context.Background(), &EmailRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Should NOT be an APIError (it's a transport error).
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		t.Errorf("expected non-APIError, got %v", apiErr)
	}
}

func TestSendEmail_ContextCancelled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	client := newTestEmailClient(ts.URL)
	_, err := client.SendEmail(ctx, &EmailRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context canceled error, got: %v", err)
	}
}

func TestSendEmail_ContextTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendEmail(ctx, &EmailRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- SendBatchEmail ---

func TestSendBatchEmail_Success(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendBatchEmail(context.Background(), &EmailRequest{
		From: EmailAddress{Address: "a@b.com"},
		To:   []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/email/batch" {
		t.Errorf("path = %q, want %q", gotPath, "/email/batch")
	}
}

// --- SendTemplateEmail ---

func TestSendTemplateEmail_Success(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendTemplateEmail(context.Background(), &TemplateRequest{
		TemplateKey: "tpl-1",
		From:        EmailAddress{Address: "a@b.com"},
		To:          []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/email/template" {
		t.Errorf("path = %q, want %q", gotPath, "/email/template")
	}
}

// --- SendBatchTemplateEmail ---

func TestSendBatchTemplateEmail_Success(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.SendBatchTemplateEmail(context.Background(), &TemplateRequest{
		TemplateKey: "tpl-1",
		From:        EmailAddress{Address: "a@b.com"},
		To:          []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/email/template/batch" {
		t.Errorf("path = %q, want %q", gotPath, "/email/template/batch")
	}
}

// --- FileCacheUpload ---

func TestFileCacheUpload_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"file_cache_key":"fck-123","data":[],"message":"OK","object":"file"}`))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	resp, err := client.FileCacheUpload(context.Background(), "test.txt", []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.FileCacheKey != "fck-123" {
		t.Errorf("FileCacheKey = %q, want %q", resp.FileCacheKey, "fck-123")
	}
}

func TestFileCacheUpload_FilenameEncoding(t *testing.T) {
	var gotQuery string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(200)
		w.Write([]byte(`{"file_cache_key":"fck-123"}`))
	}))
	defer ts.Close()

	client := newTestEmailClient(ts.URL)
	_, err := client.FileCacheUpload(context.Background(), "my file (1).txt", []byte("data"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotQuery, "name=my+file+%281%29.txt") {
		t.Errorf("query = %q, want encoded filename", gotQuery)
	}
}

// --- CreateTemplate ---

func TestCreateTemplate_Success(t *testing.T) {
	var gotPath, gotMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(200)
		w.Write([]byte(`{"data":[{"template_name":"tpl","template_key":"k1"}],"message":"OK","object":"template"}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	resp, err := client.CreateTemplate(context.Background(), "agent-1", &CreateTemplateRequest{
		TemplateName: "tpl",
		Subject:      "subj",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != "POST" {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/mailagents/agent-1/templates" {
		t.Errorf("path = %q, want %q", gotPath, "/mailagents/agent-1/templates")
	}
	if len(resp.Data) != 1 || resp.Data[0].TemplateName != "tpl" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestCreateTemplate_AuthHeader(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
		w.Write([]byte(`{"data":[],"message":"OK","object":"template"}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	_, err := client.CreateTemplate(context.Background(), "agent", &CreateTemplateRequest{TemplateName: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Zoho-oauthtoken test-oauth-token" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Zoho-oauthtoken test-oauth-token")
	}
}

// --- GetTemplate ---

func TestGetTemplate_Success(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(200)
		w.Write([]byte(`{"data":{"template_name":"tpl","template_key":"k1","subject":"s","htmlbody":"<b>hi</b>","created_time":"t","modified_time":"t"},"message":"OK","object":"template"}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	resp, err := client.GetTemplate(context.Background(), "agent-1", "key-1")
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/mailagents/agent-1/templates/key-1" {
		t.Errorf("path = %q, want %q", gotPath, "/mailagents/agent-1/templates/key-1")
	}
	if resp.Data.TemplateName != "tpl" {
		t.Errorf("TemplateName = %q, want %q", resp.Data.TemplateName, "tpl")
	}
}

// --- UpdateTemplate ---

func TestUpdateTemplate_Success(t *testing.T) {
	var gotMethod string
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write([]byte(`{"data":[{"template_name":"updated"}],"message":"OK","object":"template"}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	_, err := client.UpdateTemplate(context.Background(), "agent", "key", &CreateTemplateRequest{
		TemplateName: "updated",
		Subject:      "new subj",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != "PUT" {
		t.Errorf("method = %q, want PUT", gotMethod)
	}
	if !strings.Contains(string(gotBody), `"template_name":"updated"`) {
		t.Errorf("body = %q, want to contain template_name", string(gotBody))
	}
}

// --- ListTemplates ---

func TestListTemplates_Success(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.RequestURI()
		w.WriteHeader(200)
		w.Write([]byte(`{"metadata":{"offset":0,"count":1,"limit":10},"data":[{"template_name":"t1","template_key":"k1","created_time":"t","modified_time":"t","subject":"s"}],"message":"OK"}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	resp, err := client.ListTemplates(context.Background(), "agent", ListTemplatesParams{Offset: 0, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(gotPath, "offset=0") || !strings.Contains(gotPath, "limit=10") {
		t.Errorf("path = %q, want offset and limit params", gotPath)
	}
	if len(resp.Data) != 1 {
		t.Errorf("data length = %d, want 1", len(resp.Data))
	}
}

// --- DeleteTemplate ---

func TestDeleteTemplate_Success(t *testing.T) {
	var gotMethod string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(200)
		w.Write([]byte(`{"message":"OK"}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	err := client.DeleteTemplate(context.Background(), "agent", "key")
	if err != nil {
		t.Fatal(err)
	}
	if gotMethod != "DELETE" {
		t.Errorf("method = %q, want DELETE", gotMethod)
	}
}

func TestDeleteTemplate_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error":{"code":"NOT_FOUND","message":"Template not found","details":[],"request_id":"req-404"}}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	err := client.DeleteTemplate(context.Background(), "agent", "missing-key")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.Code != "NOT_FOUND" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "NOT_FOUND")
	}
}

// --- Path Escaping ---

func TestTemplatesClient_PathEscaping(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.RawPath
		if gotPath == "" {
			gotPath = r.URL.Path
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"data":{"template_name":"t","template_key":"k","subject":"s","htmlbody":"","created_time":"t","modified_time":"t"},"message":"OK","object":"template"}`))
	}))
	defer ts.Close()

	client := newTestTemplatesClient(ts.URL)
	_, err := client.GetTemplate(context.Background(), "agent/with/slashes", "key with spaces")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(gotPath, "agent%2Fwith%2Fslashes") {
		t.Errorf("path = %q, want agent alias to be escaped", gotPath)
	}
	if !strings.Contains(gotPath, "key%20with%20spaces") {
		t.Errorf("path = %q, want template key to be escaped", gotPath)
	}
}

// --- Options ---

func TestNewEmailClient_DefaultTimeout(t *testing.T) {
	client := NewEmailClient("key")
	if client.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestNewEmailClient_WithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 99 * time.Second}
	client := NewEmailClient("key", WithHTTPClient(custom))
	if client.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestNewEmailClient_WithBaseURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(successJSON()))
	}))
	defer ts.Close()

	client := NewEmailClient("key", WithBaseURL(ts.URL))
	_, err := client.SendEmail(context.Background(), &EmailRequest{
		From: EmailAddress{Address: "a@b.com"},
		To:   []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
	})
	if err != nil {
		t.Fatalf("expected request to custom base URL to succeed, got: %v", err)
	}
}
