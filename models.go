package zeptomail

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name,omitempty"`
}

type Recipient struct {
	EmailAddress `json:"email_address"`
	MergeInfo    map[string]string `json:"merge_info,omitempty"`
}

// Attachment can reference a previously uploaded file (FileCacheKey) or carry
// base64 content directly (Content + MimeType + Name).
type Attachment struct {
	FileCacheKey string `json:"file_cache_key,omitempty"`
	Content      string `json:"content,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	Name         string `json:"name,omitempty"`
}

// InlineImage is an image embedded via a CID reference in the HTML body.
type InlineImage struct {
	CID          string `json:"cid"`
	Content      string `json:"content,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileCacheKey string `json:"file_cache_key,omitempty"`
}

// ListTemplatesParams controls pagination for ListTemplates.
type ListTemplatesParams struct {
	Offset int
	Limit  int
}

// CreateTemplateRequest is the body for creating (or updating) a template.
type CreateTemplateRequest struct {
	TemplateName  string `json:"template_name"`
	TemplateAlias string `json:"template_alias,omitempty"`
	Subject       string `json:"subject"`
	HTMLBody      string `json:"htmlbody,omitempty"`
	TextBody      string `json:"textbody,omitempty"`
}

// UpdateTemplateRequest is an alias — the create and update payloads are identical.
type UpdateTemplateRequest = CreateTemplateRequest

// EmailRequest is used by SendEmail and SendBatchEmail.
type EmailRequest struct {
	From            EmailAddress      `json:"from"`
	To              []Recipient       `json:"to"`
	Cc              []Recipient       `json:"cc,omitempty"`
	Bcc             []Recipient       `json:"bcc,omitempty"`
	ReplyTo         []EmailAddress    `json:"reply_to,omitempty"`
	Subject         string            `json:"subject"`
	HTMLBody        string            `json:"htmlbody,omitempty"`
	TextBody        string            `json:"textbody,omitempty"`
	Attachments     []Attachment      `json:"attachments,omitempty"`
	InlineImages    []InlineImage     `json:"inline_images,omitempty"`
	BounceAddress   string            `json:"bounce_address,omitempty"`
	TrackClicks     *bool             `json:"track_clicks,omitempty"`
	TrackOpens      *bool             `json:"track_opens,omitempty"`
	ClientReference string            `json:"client_reference,omitempty"`
	MimeHeaders     map[string]string `json:"mime_headers,omitempty"`
	MergeInfo       map[string]string `json:"merge_info,omitempty"`
}

// TemplateRequest is used by SendTemplateEmail and SendBatchTemplateEmail.
type TemplateRequest struct {
	TemplateKey     string            `json:"template_key,omitempty"`
	TemplateAlias   string            `json:"template_alias,omitempty"`
	BounceAddress   string            `json:"bounce_address,omitempty"`
	From            EmailAddress      `json:"from"`
	To              []Recipient       `json:"to"`
	Cc              []Recipient       `json:"cc,omitempty"`
	Bcc             []Recipient       `json:"bcc,omitempty"`
	ReplyTo         []EmailAddress    `json:"reply_to,omitempty"`
	Subject         string            `json:"subject,omitempty"`
	HTMLBody        string            `json:"htmlbody,omitempty"`
	TextBody        string            `json:"textbody,omitempty"`
	Attachments     []Attachment      `json:"attachments,omitempty"`
	InlineImages    []InlineImage     `json:"inline_images,omitempty"`
	TrackClicks     *bool             `json:"track_clicks,omitempty"`
	TrackOpens      *bool             `json:"track_opens,omitempty"`
	ClientReference string            `json:"client_reference,omitempty"`
	MimeHeaders     map[string]string `json:"mime_headers,omitempty"`
	MergeInfo       map[string]string `json:"merge_info,omitempty"`
}

type AdditionalInfo struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResponseData struct {
	Code           string           `json:"code"`
	AdditionalInfo []AdditionalInfo `json:"additional_info"`
	Message        string           `json:"message"`
}

type SuccessResponse struct {
	Data      []ResponseData `json:"data"`
	Message   string         `json:"message"`
	RequestID string         `json:"request_id"`
	Object    string         `json:"object"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target"`
}

type ErrorResponse struct {
	Error struct {
		Code      string        `json:"code"`
		Details   []ErrorDetail `json:"details"`
		Message   string        `json:"message"`
		RequestID string        `json:"request_id"`
	} `json:"error"`
}

type FileUploadResponse struct {
	FileCacheKey string         `json:"file_cache_key"`
	Data         []ResponseData `json:"data"`
	Message      string         `json:"message"`
	Object       string         `json:"object"`
}

type TemplateAttachment struct {
	FileCacheKey string `json:"file_cache_key"`
	ContentType  string `json:"content_type"`
	FileName     string `json:"file_name"`
}

type TemplateData struct {
	HTMLBody        string               `json:"htmlbody"`
	CreatedTime     string               `json:"created_time"`
	ModifiedTime    string               `json:"modified_time"`
	TemplateName    string               `json:"template_name"`
	TemplateKey     string               `json:"template_key"`
	TemplateAlias   string               `json:"template_alias,omitempty"`
	Subject         string               `json:"subject"`
	Attachments     []TemplateAttachment `json:"attachments,omitempty"`
	SampleMergeInfo map[string]string    `json:"sample_merge_info,omitempty"`
}

type CreateTemplateResponse struct {
	Data    []TemplateData `json:"data"`
	Message string         `json:"message"`
	Object  string         `json:"object"`
}

type GetTemplateResponse struct {
	Data    TemplateData `json:"data"`
	Message string       `json:"message"`
	Object  string       `json:"object"`
}

type ListTemplatesMetadata struct {
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Limit  int `json:"limit"`
}

type TemplateListItem struct {
	CreatedTime   string `json:"created_time"`
	TemplateName  string `json:"template_name"`
	TemplateKey   string `json:"template_key"`
	ModifiedTime  string `json:"modified_time"`
	Subject       string `json:"subject"`
	TemplateAlias string `json:"template_alias,omitempty"`
}

type ListTemplatesResponse struct {
	Metadata ListTemplatesMetadata `json:"metadata"`
	Data     []TemplateListItem    `json:"data"`
	Message  string                `json:"message"`
}
