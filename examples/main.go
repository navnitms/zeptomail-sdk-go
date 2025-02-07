package main

import (
	"fmt"
	"log"
	"os"

	"github.com/navnitms/zeptomail-sdk-go/zeptomail"
)

func main() {
	emailOperations()
	templateOperations()
}

func emailOperations() {
	emailClient := zeptomail.NewEmailClient("YOUR-API-KEY")
	sendSingleEmail(emailClient)
	sendBatchEmail(emailClient)
	sendTemplateEmail(emailClient)
	sendBatchTemplateEmail(emailClient)
	fileUploadCache(emailClient)
}

func templateOperations() {
	templatesClient := zeptomail.NewTemplatesClient("YOUR-OAUTH-TOKEN")
	createTemplate(templatesClient)
	getTemplate(templatesClient)
	updateTemplate(templatesClient)
	listTemplates(templatesClient)
	deleteTemplate(templatesClient)
}

func sendSingleEmail(client *zeptomail.EmailClient) {
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

	resp, err := client.SendEmail(emailReq)
	if err != nil {
		log.Fatal("Failed to send email:", err)
	}
	fmt.Printf("Email sent successfully: %+v\n", resp)
}

func sendBatchEmail(client *zeptomail.EmailClient) {
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

	resp, err := client.SendBatchEmail(batchEmailReq)
	if err != nil {
		log.Fatal("Failed to send batch email:", err)
	}
	fmt.Printf("Batch email sent successfully: %+v\n", resp)
}

func sendTemplateEmail(client *zeptomail.EmailClient) {
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

	resp, err := client.SendTemplateEmail(templateReq)
	if err != nil {
		log.Fatal("Failed to send template email:", err)
	}
	fmt.Printf("Template email sent successfully: %+v\n", resp)
}

func sendBatchTemplateEmail(client *zeptomail.EmailClient) {
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

	resp, err := client.SendBatchTemplateEmail(batchTemplateReq)
	if err != nil {
		log.Fatal("Failed to send batch template email:", err)
	}
	fmt.Printf("Batch template email sent successfully: %+v\n", resp)
}

func fileUploadCache(client *zeptomail.EmailClient) {
	fileContent, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatal("Failed to read file:", err)
	}

	resp, err := client.FileCacheUpload("test.txt", fileContent)
	if err != nil {
		log.Fatal("Failed to upload file:", err)
	}
	fmt.Printf("File uploaded: %s\n", resp.FileCacheKey)
}

func createTemplate(client *zeptomail.TemplatesClient) {
	createReq := &zeptomail.CreateTemplateRequest{
		TemplateName: "Welcome Email",
		Subject:      "Welcome to our service",
		HTMLBody:     "<h1>Welcome {{name}}!</h1>",
	}

	resp, err := client.CreateTemplate("your-mailagent-alias", createReq)
	if err != nil {
		log.Fatal("Failed to create template:", err)
	}
	fmt.Printf("Template created: %+v\n", resp)
}

func getTemplate(client *zeptomail.TemplatesClient) {
	resp, err := client.GetTemplate("your-mailagent-alias", "template-key")
	if err != nil {
		log.Fatal("Failed to get template:", err)
	}
	fmt.Printf("Template details: %+v\n", resp)
}

func updateTemplate(client *zeptomail.TemplatesClient) {
	updateReq := &zeptomail.UpdateTemplateRequest{
		TemplateName: "Welcome Email (Updated)",
		Subject:      "Welcome to our service - New Version",
		HTMLBody:     "<h1>Welcome {{name}}!</h1><p>We're glad to have you.</p>",
	}

	resp, err := client.UpdateTemplate("your-mailagent-alias", "template-key", updateReq)
	if err != nil {
		log.Fatal("Failed to update template:", err)
	}
	fmt.Printf("Template updated: %+v\n", resp)
}

func listTemplates(client *zeptomail.TemplatesClient) {
	resp, err := client.ListTemplates("your-mailagent-alias", zeptomail.ListTemplatesParams{
		Offset: 0,
		Limit:  10,
	})
	if err != nil {
		log.Fatal("Failed to list templates:", err)
	}
	fmt.Printf("Templates: %+v\n", resp)
}

func deleteTemplate(client *zeptomail.TemplatesClient) {
	err := client.DeleteTemplate("your-mailagent-alias", "template-key")
	if err != nil {
		log.Fatal("Failed to delete template:", err)
	}
	fmt.Println("Template deleted successfully")
}
