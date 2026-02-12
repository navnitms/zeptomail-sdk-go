package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/navnitms/zeptomail-sdk-go"
)

func main() {
	emailOperations()
	templateOperations()
}

func emailOperations() {
	emailClient := zeptomail.NewEmailClient("YOUR-API-KEY",
		zeptomail.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
	)

	ctx := context.Background()
	sendSingleEmail(ctx, emailClient)
	sendBatchEmail(ctx, emailClient)
	sendTemplateEmail(ctx, emailClient)
	sendBatchTemplateEmail(ctx, emailClient)
	fileUploadCache(ctx, emailClient)
}

func templateOperations() {
	templatesClient := zeptomail.NewTemplatesClient("YOUR-OAUTH-TOKEN")

	ctx := context.Background()
	createTemplate(ctx, templatesClient)
	getTemplate(ctx, templatesClient)
	updateTemplate(ctx, templatesClient)
	listTemplates(ctx, templatesClient)
	deleteTemplate(ctx, templatesClient)
}

func sendSingleEmail(ctx context.Context, client *zeptomail.EmailClient) {
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
				Content:  "SGVsbG8gV29ybGQ=", // base64 of "Hello World"
				MimeType: "text/plain",
				Name:     "hello.txt",
			},
		},
		InlineImages: []zeptomail.InlineImage{
			{
				CID:      "logo",
				Content:  "iVBORw0KGgo=", // truncated base64 PNG
				MimeType: "image/png",
			},
		},
	}

	resp, err := client.SendEmail(ctx, emailReq)
	if err != nil {
		var apiErr *zeptomail.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API error %d: %s — %s\n", apiErr.HTTPStatusCode, apiErr.Code, apiErr.Message)
			return
		}
		log.Fatal("Failed to send email:", err)
	}
	fmt.Printf("Email sent successfully: %+v\n", resp)
}

func sendBatchEmail(ctx context.Context, client *zeptomail.EmailClient) {
	batchEmailReq := &zeptomail.EmailRequest{
		From: zeptomail.EmailAddress{
			Address: "sender@example.com",
			Name:    "Invoice",
		},
		To: []zeptomail.Recipient{
			{
				EmailAddress: zeptomail.EmailAddress{
					Address: "recipient@example.com",
					Name:    "Paul",
				},
				MergeInfo: map[string]string{
					"contact": "98********",
					"company": "Test Company",
				},
			},
			{
				EmailAddress: zeptomail.EmailAddress{
					Address: "recipient@example2.com",
					Name:    "Rebecca",
				},
				MergeInfo: map[string]string{
					"contact": "87********",
					"company": "Test Company 2",
				},
			},
		},
		HTMLBody: "<div><b>This is a sample email.{{contact}} {{company}}</b></div>",
		TextBody: "This is a sample email",
	}

	resp, err := client.SendBatchEmail(ctx, batchEmailReq)
	if err != nil {
		log.Fatal("Failed to send batch email:", err)
	}
	fmt.Printf("Batch email sent successfully: %+v\n", resp)
}

func sendTemplateEmail(ctx context.Context, client *zeptomail.EmailClient) {
	templateReq := &zeptomail.TemplateRequest{
		TemplateKey:   "template-key",
		BounceAddress: "bounce@bounce.zylker.com",
		From: zeptomail.EmailAddress{
			Address: "rebecca@zylker.com",
			Name:    "Rebecca",
		},
		To: []zeptomail.Recipient{
			{
				EmailAddress: zeptomail.EmailAddress{
					Address: "paula@zylker.com",
					Name:    "Paula M",
				},
			},
		},
		MergeInfo: map[string]string{
			"meeting_link": "https://meeting.zoho.com/join?key=1234",
		},
	}

	resp, err := client.SendTemplateEmail(ctx, templateReq)
	if err != nil {
		log.Fatal("Failed to send template email:", err)
	}
	fmt.Printf("Template email sent successfully: %+v\n", resp)
}

func sendBatchTemplateEmail(ctx context.Context, client *zeptomail.EmailClient) {
	batchTemplateReq := &zeptomail.TemplateRequest{
		TemplateKey:   "template-key",
		BounceAddress: "bounce@bounce.zylker.com",
		From: zeptomail.EmailAddress{
			Address: "rebecca@zylker.com",
			Name:    "Rebecca",
		},
		To: []zeptomail.Recipient{
			{
				EmailAddress: zeptomail.EmailAddress{
					Address: "paula@zylker.com",
					Name:    "Paula",
				},
				MergeInfo: map[string]string{
					"contact": "960*******23",
					"company": "Zylker",
				},
			},
			{
				EmailAddress: zeptomail.EmailAddress{
					Address: "charles@zylker.com",
					Name:    "Charles",
				},
				MergeInfo: map[string]string{
					"contact": "860*******13",
					"company": "Zillum",
				},
			},
		},
	}

	resp, err := client.SendBatchTemplateEmail(ctx, batchTemplateReq)
	if err != nil {
		log.Fatal("Failed to send batch template email:", err)
	}
	fmt.Printf("Batch template email sent successfully: %+v\n", resp)
}

func fileUploadCache(ctx context.Context, client *zeptomail.EmailClient) {
	fileContent, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatal("Failed to read file:", err)
	}

	resp, err := client.FileCacheUpload(ctx, "test.txt", fileContent)
	if err != nil {
		log.Fatal("Failed to upload file:", err)
	}
	fmt.Printf("File uploaded: %s\n", resp.FileCacheKey)
}

func createTemplate(ctx context.Context, client *zeptomail.TemplatesClient) {
	createReq := &zeptomail.CreateTemplateRequest{
		TemplateName: "Welcome Email",
		Subject:      "Welcome to our service",
		HTMLBody:     "<h1>Welcome {{name}}!</h1>",
	}

	resp, err := client.CreateTemplate(ctx, "your-mailagent-alias", createReq)
	if err != nil {
		log.Fatal("Failed to create template:", err)
	}
	fmt.Printf("Template created: %+v\n", resp)
}

func getTemplate(ctx context.Context, client *zeptomail.TemplatesClient) {
	resp, err := client.GetTemplate(ctx, "your-mailagent-alias", "template-key")
	if err != nil {
		log.Fatal("Failed to get template:", err)
	}
	fmt.Printf("Template details: %+v\n", resp)
}

func updateTemplate(ctx context.Context, client *zeptomail.TemplatesClient) {
	updateReq := &zeptomail.UpdateTemplateRequest{
		TemplateName: "Welcome Email (Updated)",
		Subject:      "Welcome to our service - New Version",
		HTMLBody:     "<h1>Welcome {{name}}!</h1><p>We're glad to have you.</p>",
	}

	resp, err := client.UpdateTemplate(ctx, "your-mailagent-alias", "template-key", updateReq)
	if err != nil {
		log.Fatal("Failed to update template:", err)
	}
	fmt.Printf("Template updated: %+v\n", resp)
}

func listTemplates(ctx context.Context, client *zeptomail.TemplatesClient) {
	resp, err := client.ListTemplates(ctx, "your-mailagent-alias", zeptomail.ListTemplatesParams{
		Offset: 0,
		Limit:  10,
	})
	if err != nil {
		log.Fatal("Failed to list templates:", err)
	}
	fmt.Printf("Templates: %+v\n", resp)
}

func deleteTemplate(ctx context.Context, client *zeptomail.TemplatesClient) {
	err := client.DeleteTemplate(ctx, "your-mailagent-alias", "template-key")
	if err != nil {
		log.Fatal("Failed to delete template:", err)
	}
	fmt.Println("Template deleted successfully")
}
