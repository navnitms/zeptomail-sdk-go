package zeptomail

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{
		HTTPStatusCode: 400,
		Code:           "INVALID_DATA",
		Message:        "The email address is invalid",
		RequestID:      "abc123",
	}

	got := err.Error()
	want := `zeptomail: HTTP 400: INVALID_DATA - The email address is invalid (request_id: abc123)`
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAPIError_Error_NoRequestID(t *testing.T) {
	err := &APIError{
		HTTPStatusCode: 500,
		Code:           "GENERAL_ERROR",
		Message:        "Something went wrong",
	}

	got := err.Error()
	want := `zeptomail: HTTP 500: GENERAL_ERROR - Something went wrong`
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAPIError_ErrorsAs(t *testing.T) {
	original := &APIError{
		HTTPStatusCode: 422,
		Code:           "INVALID_DATA",
		Message:        "Missing required field",
		Details: []ErrorDetail{
			{Code: "REQUIRED", Message: "to is required", Target: "to"},
		},
		RequestID: "req-1",
	}

	// Wrap it
	wrapped := errors.New("outer: " + original.Error())
	_ = wrapped

	// errors.As should work with the original
	var apiErr *APIError
	err := error(original)
	if !errors.As(err, &apiErr) {
		t.Fatal("errors.As failed to match *APIError")
	}

	if apiErr.HTTPStatusCode != 422 {
		t.Errorf("HTTPStatusCode = %d, want 422", apiErr.HTTPStatusCode)
	}
	if apiErr.Code != "INVALID_DATA" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "INVALID_DATA")
	}
	if len(apiErr.Details) != 1 {
		t.Fatalf("Details length = %d, want 1", len(apiErr.Details))
	}
	if apiErr.Details[0].Target != "to" {
		t.Errorf("Details[0].Target = %q, want %q", apiErr.Details[0].Target, "to")
	}
}
