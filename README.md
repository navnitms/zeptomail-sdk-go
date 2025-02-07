# zeptomail-sdk-go

A Go SDK for interacting with the ZeptoMail API. This SDK provides a simple and efficient way to integrate ZeptoMail services into your Go applications. Features include sending emails (single and batch), managing email templates, and handling file uploads.

## Installation

```bash
go get github.com/navnitms/zeptomail-sdk-go
```

## Quick Start Guide

```go
package main

import (
    "github.com/navnitms/zeptomail-sdk-go/zeptomail"
)

func main() {
    // Initialize email client for sending emails and file uploads
    // Make sure you only pass the API KEY excluing the "Zoho-enczapikey" prefix 
    emailClient := zeptomail.NewEmailClient("YOUR-API-KEY")

    // Initialize templates client for managing email templates
    templatesClient := zeptomail.NewTemplatesClient("YOUR-TEMPLATES-TOKEN")
}
```

## Usage Examples

### Send Single Email
```go
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
}

resp, err := emailClient.SendEmail(emailReq)
if err != nil {
    log.Fatal(err)
}
```

### Send Batch Email with Merge Fields
```go
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

resp, err := emailClient.SendBatchEmail(batchEmailReq)
if err != nil {
    log.Fatal(err)
}
```

### Template Management
```go
// Create a new template
createReq := &zeptomail.CreateTemplateRequest{
    TemplateName: "Welcome Email",
    Subject:     "Welcome to our service",
    HTMLBody:    "<h1>Welcome {{name}}!</h1>",
}

template, err := templatesClient.CreateTemplate("your-mailagent-alias", createReq)
if err != nil {
    log.Fatal(err)
}

// Send email using template
templateReq := &zeptomail.TemplateRequest{
    TemplateKey: "template-key",
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
            MergeInfo: map[string]string{
                "name": "John Doe",
            },
        },
    },
}

resp, err := emailClient.SendTemplateEmail(templateReq)
if err != nil {
    log.Fatal(err)
}
```

### File Upload
```go
// Read file content
fileContent, err := os.ReadFile("attachment.pdf")
if err != nil {
    log.Fatal(err)
}

// Upload file to ZeptoMail
fileResp, err := emailClient.FileCacheUpload("attachment.pdf", fileContent)
if err != nil {
    log.Fatal(err)
}

// Use the file in an email
emailReq := &zeptomail.EmailRequest{
    // ... other email fields ...
    Attachments: []zeptomail.Attachment{
        {
            FileCacheKey: fileResp.FileCacheKey,
        },
    },
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
