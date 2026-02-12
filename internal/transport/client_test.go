package transport

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequest_SetsHeaders(t *testing.T) {
	var gotHeaders http.Header
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeaders = r.Header
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := NewEmailClient("test-key", ts.URL, ts.Client())
	_, err := c.Request(context.Background(), "POST", "/email", map[string]string{"key": "val"})
	if err != nil {
		t.Fatal(err)
	}

	if got := gotHeaders.Get("Accept"); got != "application/json" {
		t.Errorf("Accept = %q, want %q", got, "application/json")
	}
	if got := gotHeaders.Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
	if got := gotHeaders.Get("Authorization"); got != "Zoho-enczapikey test-key" {
		t.Errorf("Authorization = %q, want %q", got, "Zoho-enczapikey test-key")
	}
}

func TestRequest_MarshalPayload(t *testing.T) {
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := NewEmailClient("key", ts.URL, ts.Client())
	_, err := c.Request(context.Background(), "POST", "/email", map[string]string{"subject": "hi"})
	if err != nil {
		t.Fatal(err)
	}

	if string(gotBody) != `{"subject":"hi"}` {
		t.Errorf("body = %q, want %q", string(gotBody), `{"subject":"hi"}`)
	}
}

func TestRequest_NilPayload(t *testing.T) {
	var gotBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := NewEmailClient("key", ts.URL, ts.Client())
	_, err := c.Request(context.Background(), "GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(gotBody) != 0 {
		t.Errorf("body = %q, want empty", string(gotBody))
	}
}

func TestRequest_ReturnsStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		w.Write([]byte(`{"error":{"code":"BAD"}}`))
	}))
	defer ts.Close()

	c := NewEmailClient("key", ts.URL, ts.Client())
	resp, err := c.Request(context.Background(), "POST", "/email", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 422 {
		t.Errorf("StatusCode = %d, want 422", resp.StatusCode)
	}
}

func TestSetAuthHeader_API(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := NewEmailClient("my-api-key", ts.URL, ts.Client())
	_, err := c.Request(context.Background(), "GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	if gotAuth != "Zoho-enczapikey my-api-key" {
		t.Errorf("auth = %q, want %q", gotAuth, "Zoho-enczapikey my-api-key")
	}
}

func TestSetAuthHeader_Templates(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	c := NewTemplatesClient("my-oauth-token", ts.URL, ts.Client())
	_, err := c.Request(context.Background(), "GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	if gotAuth != "Zoho-oauthtoken my-oauth-token" {
		t.Errorf("auth = %q, want %q", gotAuth, "Zoho-oauthtoken my-oauth-token")
	}
}

func TestUpload_ContentTypeDetection(t *testing.T) {
	// PNG header bytes
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	var gotContentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(200)
		w.Write([]byte(`{"file_cache_key":"abc"}`))
	}))
	defer ts.Close()

	c := NewEmailClient("key", ts.URL, ts.Client())
	_, err := c.Upload(context.Background(), "/files", "image.png", pngHeader)
	if err != nil {
		t.Fatal(err)
	}

	if gotContentType != "image/png" {
		t.Errorf("Content-Type = %q, want %q", gotContentType, "image/png")
	}
}

func TestUpload_TextContent(t *testing.T) {
	var gotContentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(200)
		w.Write([]byte(`{"file_cache_key":"abc"}`))
	}))
	defer ts.Close()

	c := NewEmailClient("key", ts.URL, ts.Client())
	_, err := c.Upload(context.Background(), "/files", "file.txt", []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}

	if gotContentType != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", gotContentType, "text/plain; charset=utf-8")
	}
}
