# zeptomail-sdk-go

A zero-dependency Go SDK for interacting with the ZeptoMail API. This SDK provides a simple and efficient way to integrate ZeptoMail services into your Go applications. Features include sending emails (single and batch), managing email templates, and handling file uploads.

## Installation

```bash
go get github.com/navnitms/zeptomail-sdk-go
```

## Quick Start Guide

```go
package main

import (
    "context"
    "net/http"
    "time"

    "github.com/navnitms/zeptomail-sdk-go"
)

func main() {
    // Initialize email client for sending emails and file uploads.
    // Pass only the API key (without the "Zoho-enczapikey" prefix).
    emailClient := zeptomail.NewEmailClient("YOUR-API-KEY",
        zeptomail.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
    )

    // Initialize templates client for managing email templates
    templatesClient := zeptomail.NewTemplatesClient("YOUR-OAUTH-TOKEN")

    ctx := context.Background()
    _ = emailClient
    _ = templatesClient
    _ = ctx
}
```

## Usage Examples

### Send Single Email
```go
ctx := context.Background()

emailReq := &zeptomail.EmailRequest{
    From: zeptomail.EmailAddress{
        Address: "sender@example.com",
        Name:    "Sender",
    },
    To: []zeptomail.Recipient{
        {
            EmailAddress: zeptomail.EmailAddress{
                Address: "recipient@example.com",
                Name:    "Recipient",
            },
        },
    },
    Subject:  "Test Email",
    HTMLBody: "<h1>Hello from ZeptoMail SDK!</h1>",
    Attachments: []zeptomail.Attachment{
        {
            Content:  "SGVsbG8gV29ybGQ=", // base64-encoded content
            MimeType: "text/plain",
            Name:     "hello.txt",
        },
    },
    InlineImages: []zeptomail.InlineImage{
        {
            CID:      "logo",
            Content:  "iVBORw0KGgo=", // base64-encoded image
            MimeType: "image/png",
        },
    },
}

resp, err := emailClient.SendEmail(ctx, emailReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Email sent: %+v\n", resp)
```

### Send Batch Email with Merge Fields
```go
ctx := context.Background()

batchEmailReq := &zeptomail.EmailRequest{
    From: zeptomail.EmailAddress{
        Address: "sender@example.com",
        Name:    "Sender",
    },
    To: []zeptomail.Recipient{
        {
            EmailAddress: zeptomail.EmailAddress{
                Address: "recipient1@example.com",
                Name:    "Recipient 1",
            },
            MergeInfo: map[string]string{
                "name": "John",
                "company": "Example Corp",
            },
        },
        {
            EmailAddress: zeptomail.EmailAddress{
                Address: "recipient2@example.com",
                Name:    "Recipient 2",
            },
            MergeInfo: map[string]string{
                "name": "Jane",
                "company": "Sample Inc",
            },
        },
    },
    HTMLBody: "<h1>Hello {{name}}</h1><p>Welcome to {{company}}!</p>",
}

resp, err := emailClient.SendBatchEmail(ctx, batchEmailReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Batch email sent: %+v\n", resp)
```

### Send Template Email
```go
ctx := context.Background()

templateReq := &zeptomail.TemplateRequest{
    TemplateKey:   "template-key",
    BounceAddress: "bounce@bounce.example.com",
    From: zeptomail.EmailAddress{
        Address: "sender@example.com",
        Name:    "Sender",
    },
    To: []zeptomail.Recipient{
        {
            EmailAddress: zeptomail.EmailAddress{
                Address: "recipient@example.com",
                Name:    "Recipient",
            },
        },
    },
    MergeInfo: map[string]string{
        "name": "John Doe",
    },
}

resp, err := emailClient.SendTemplateEmail(ctx, templateReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Template email sent: %+v\n", resp)
```

### Template Management
```go
ctx := context.Background()

// Create a new template
createReq := &zeptomail.CreateTemplateRequest{
    TemplateName: "Welcome Email",
    Subject:      "Welcome to our service",
    HTMLBody:     "<h1>Welcome {{name}}!</h1>",
}

template, err := templatesClient.CreateTemplate(ctx, "your-mailagent-alias", createReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Template created: %+v\n", template)
```

### File Upload
```go
ctx := context.Background()

// Read file content
fileContent, err := os.ReadFile("attachment.pdf")
if err != nil {
    log.Fatal(err)
}

// Upload file to ZeptoMail
fileResp, err := emailClient.FileCacheUpload(ctx, "attachment.pdf", fileContent)
if err != nil {
    log.Fatal(err)
}

// Use the file cache key in an email
emailReq := &zeptomail.EmailRequest{
    // ... other email fields ...
    Attachments: []zeptomail.Attachment{
        {
            FileCacheKey: fileResp.FileCacheKey,
        },
    },
}
```

## Error Handling

All API errors are returned as `*zeptomail.APIError`, which you can inspect with `errors.As`:

```go
resp, err := emailClient.SendEmail(ctx, emailReq)
if err != nil {
    var apiErr *zeptomail.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s — %s\n", apiErr.HTTPStatusCode, apiErr.Code, apiErr.Message)
        return
    }
    log.Fatal("unexpected error:", err)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
