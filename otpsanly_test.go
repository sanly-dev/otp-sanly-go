package otpsanly

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &Client{APIKey: "otpsanly_test", BaseURL: srv.URL}
}

func TestSendOtp_RequiresPhoneOrEmail(t *testing.T) {
	c := &Client{APIKey: "otpsanly_test"}
	if _, err := c.SendOtp(context.Background(), SendOtpParams{}); err == nil {
		t.Fatal("expected error when neither Phone nor Email is set")
	}
}

func TestSendOtp_Success(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/send-otp" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["apiKey"] != "otpsanly_test" {
			t.Errorf("apiKey not sent correctly: %v", body["apiKey"])
		}
		if body["phone"] != "+99361234567" {
			t.Errorf("phone not sent correctly: %v", body["phone"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"success": true, "otpId": 42, "target": "+99361234567", "channel": "sms",
		})
	})

	res, err := c.SendOtp(context.Background(), SendOtpParams{Phone: "+99361234567", Lang: "ru"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success || res.OtpID != 42 || res.Channel != "sms" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestSendOtp_ErrorResponse(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "Limit exceeded"})
	})

	_, err := c.SendOtp(context.Background(), SendOtpParams{Email: "a@b.com"})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Status != 400 || apiErr.Message != "Limit exceeded" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
}

func TestVerifyOtp_RequiresCode(t *testing.T) {
	c := &Client{APIKey: "otpsanly_test"}
	if _, err := c.VerifyOtp(context.Background(), VerifyOtpParams{Phone: "+99361234567"}); err == nil {
		t.Fatal("expected error when Code is missing")
	}
}

func TestVerifyOtp_RequiresTarget(t *testing.T) {
	c := &Client{APIKey: "otpsanly_test"}
	if _, err := c.VerifyOtp(context.Background(), VerifyOtpParams{Code: "123456"}); err == nil {
		t.Fatal("expected error when neither Phone nor Email is set")
	}
}

func TestVerifyOtp_Success(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"success": true, "otpId": 42, "target": "+99361234567", "verifiedAt": "2026-01-01T00:00:00Z",
		})
	})

	res, err := c.VerifyOtp(context.Background(), VerifyOtpParams{Phone: "+99361234567", Code: "123456"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success || res.OtpID != 42 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("otpsanly_test")
	if c.APIKey != "otpsanly_test" || c.BaseURL != DefaultBaseURL {
		t.Fatalf("unexpected client: %+v", c)
	}
}
