package zeptomail

import "fmt"

// APIError is returned when the ZeptoMail API responds with a non-2xx status.
type APIError struct {
	HTTPStatusCode int
	Code           string
	Message        string
	Details        []ErrorDetail
	RequestID      string
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("zeptomail: HTTP %d: %s - %s (request_id: %s)", e.HTTPStatusCode, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("zeptomail: HTTP %d: %s - %s", e.HTTPStatusCode, e.Code, e.Message)
}
