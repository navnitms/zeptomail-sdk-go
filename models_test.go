package zeptomail

import (
	"encoding/json"
	"testing"
)

func TestAttachment_JSONMarshal_FileCacheKey(t *testing.T) {
	a := Attachment{FileCacheKey: "abc-123"}
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `{"file_cache_key":"abc-123"}`
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAttachment_JSONMarshal_Base64Content(t *testing.T) {
	a := Attachment{
		Content:  "SGVsbG8=",
		MimeType: "text/plain",
		Name:     "hello.txt",
	}
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]string
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["content"] != "SGVsbG8=" {
		t.Errorf("content = %q", parsed["content"])
	}
	if parsed["mime_type"] != "text/plain" {
		t.Errorf("mime_type = %q", parsed["mime_type"])
	}
	if parsed["name"] != "hello.txt" {
		t.Errorf("name = %q", parsed["name"])
	}
	// file_cache_key should be omitted
	if _, ok := parsed["file_cache_key"]; ok {
		t.Error("file_cache_key should be omitted when empty")
	}
}

func TestInlineImage_JSONMarshal(t *testing.T) {
	img := InlineImage{
		CID:      "logo",
		Content:  "iVBOR...",
		MimeType: "image/png",
	}
	b, err := json.Marshal(img)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]string
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["cid"] != "logo" {
		t.Errorf("cid = %q", parsed["cid"])
	}
	if parsed["content"] != "iVBOR..." {
		t.Errorf("content = %q", parsed["content"])
	}
	if _, ok := parsed["file_cache_key"]; ok {
		t.Error("file_cache_key should be omitted when empty")
	}
}

func TestInlineImage_JSONMarshal_FileCacheKey(t *testing.T) {
	img := InlineImage{
		CID:          "logo",
		FileCacheKey: "fck-456",
	}
	b, err := json.Marshal(img)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]string
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["file_cache_key"] != "fck-456" {
		t.Errorf("file_cache_key = %q", parsed["file_cache_key"])
	}
}

func TestEmailRequest_JSONOmitsEmpty(t *testing.T) {
	req := EmailRequest{
		From:    EmailAddress{Address: "a@b.com"},
		To:      []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
		Subject: "hi",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	// These optional fields should not appear
	for _, field := range []string{"cc", "bcc", "reply_to", "textbody", "attachments", "inline_images", "bounce_address", "track_clicks", "track_opens"} {
		if contains(s, `"`+field+`"`) {
			t.Errorf("field %q should be omitted when empty", field)
		}
	}
}

func TestTemplateRequest_JSONMarshal(t *testing.T) {
	track := false
	req := TemplateRequest{
		TemplateKey:   "tpl-1",
		TemplateAlias: "welcome",
		BounceAddress: "bounce@test.com",
		From:          EmailAddress{Address: "a@b.com"},
		To:            []Recipient{{EmailAddress: EmailAddress{Address: "c@d.com"}}},
		TrackClicks:   &track,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["template_key"] != "tpl-1" {
		t.Errorf("template_key = %v", parsed["template_key"])
	}
	if parsed["template_alias"] != "welcome" {
		t.Errorf("template_alias = %v", parsed["template_alias"])
	}
	if parsed["bounce_address"] != "bounce@test.com" {
		t.Errorf("bounce_address = %v", parsed["bounce_address"])
	}
	if parsed["track_clicks"] != false {
		t.Errorf("track_clicks = %v", parsed["track_clicks"])
	}
}

func TestUpdateTemplateRequest_IsAlias(t *testing.T) {
	// Verify UpdateTemplateRequest is a type alias for CreateTemplateRequest
	var u UpdateTemplateRequest
	u.TemplateName = "test"
	u.Subject = "subj"

	var c CreateTemplateRequest = u
	if c.TemplateName != "test" || c.Subject != "subj" {
		t.Errorf("type alias assignment failed: %+v", c)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
