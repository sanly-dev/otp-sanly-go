// Package otpsanly is the official Go SDK for OTP Sanly (https://otp.sanly.dev) —
// SMS & Email OTP authentication for Turkmenistan.
package otpsanly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DefaultBaseURL is used when Client.BaseURL is empty.
const DefaultBaseURL = "https://otp.sanly.dev"

// Client is an OTP Sanly API client.
type Client struct {
	// APIKey is your OTP Sanly API key (starts with "otpsanly_"). Required.
	// Get one at https://otp.sanly.dev/dashboard/api-keys
	APIKey string
	// BaseURL overrides the API base URL. Defaults to DefaultBaseURL.
	BaseURL string
	// HTTPClient overrides the underlying http.Client. Defaults to a client
	// with a 15 second timeout.
	HTTPClient *http.Client
}

// NewClient creates a new OTP Sanly client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		BaseURL:    DefaultBaseURL,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// APIError is returned when the OTP Sanly API responds with an error.
type APIError struct {
	Status  int
	Message string
	Body    map[string]any
}

func (e *APIError) Error() string {
	return fmt.Sprintf("otpsanly: %s (status %d)", e.Message, e.Status)
}

// SendOtpParams are the parameters for SendOtp. Provide Phone OR Email, not both.
type SendOtpParams struct {
	// Phone is a Turkmenistan phone number, e.g. "+99361234567". SMS only.
	Phone string
	// Email is any valid email address; works worldwide.
	Email string
	// Project is a free-text label shown in your dashboard/webhooks.
	Project string
	// Lang selects which language to send the OTP in: "tk" | "ru" | "en".
	// Since this SDK calls the API server-to-server, the Accept-Language
	// header is unreliable — always set Lang explicitly if you support
	// multiple languages for your end users. Defaults to "tk" if empty.
	Lang string
}

// SendOtpResult is the response from SendOtp.
type SendOtpResult struct {
	Success      bool   `json:"success"`
	OtpID        int    `json:"otpId,omitempty"`
	Target       string `json:"target,omitempty"`
	Channel      string `json:"channel,omitempty"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
	ExpiresAtTM  string `json:"expiresAtTM,omitempty"`
	ExpiresIn    int    `json:"expiresIn,omitempty"`
	RemainingOtp int    `json:"remainingOtp,omitempty"`
	Attempt      int    `json:"attempt,omitempty"`
	MaxAttempts  int    `json:"maxAttempts,omitempty"`
	AutoRead     bool   `json:"autoRead,omitempty"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
	// Code is only present when using a sandbox API key — the OTP code
	// is returned directly for testing (no real SMS/email is sent).
	Code    string `json:"code,omitempty"`
	Sandbox bool   `json:"sandbox,omitempty"`
}

// VerifyOtpParams are the parameters for VerifyOtp. Provide Phone OR Email
// (same one used in SendOtp).
type VerifyOtpParams struct {
	Phone string
	Email string
	// Code is the code the user entered. Required.
	Code string
	// OtpID is optional — the OtpID returned by SendOtp, for extra precision
	// if you have multiple pending OTPs for the same target.
	OtpID int
}

// VerifyOtpResult is the response from VerifyOtp.
type VerifyOtpResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message,omitempty"`
	OtpID        int    `json:"otpId,omitempty"`
	Target       string `json:"target,omitempty"`
	VerifiedAt   string `json:"verifiedAt,omitempty"`
	VerifiedAtTM string `json:"verifiedAtTM,omitempty"`
	Channel      string `json:"channel,omitempty"`
	Project      string `json:"project,omitempty"`
	Error        string `json:"error,omitempty"`
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 15 * time.Second}
}

func (c *Client) baseURL() string {
	if c.BaseURL == "" {
		return DefaultBaseURL
	}
	return strings.TrimRight(c.BaseURL, "/")
}

func (c *Client) request(ctx context.Context, path string, body map[string]any, out any) error {
	if c.APIKey == "" {
		return fmt.Errorf("otpsanly: APIKey is required. Get one at https://otp.sanly.dev/dashboard/api-keys")
	}
	body["apiKey"] = c.APIKey

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("otpsanly: encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL()+path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("otpsanly: building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("otpsanly: request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("otpsanly: reading response: %w", err)
	}

	var generic map[string]any
	_ = json.Unmarshal(raw, &generic)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 || generic["success"] == false {
		msg := fmt.Sprintf("request failed with status %d", resp.StatusCode)
		if errMsg, ok := generic["error"].(string); ok && errMsg != "" {
			msg = errMsg
		}
		return &APIError{Status: resp.StatusCode, Message: msg, Body: generic}
	}

	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("otpsanly: decoding response: %w", err)
		}
	}
	return nil
}

// SendOtp sends an OTP via SMS (Turkmenistan numbers only) or email (worldwide).
func (c *Client) SendOtp(ctx context.Context, p SendOtpParams) (*SendOtpResult, error) {
	if p.Phone == "" && p.Email == "" {
		return nil, fmt.Errorf("otpsanly: SendOtp requires either Phone or Email")
	}
	body := map[string]any{}
	if p.Phone != "" {
		body["phone"] = p.Phone
	}
	if p.Email != "" {
		body["email"] = p.Email
	}
	if p.Project != "" {
		body["project"] = p.Project
	}
	if p.Lang != "" {
		body["lang"] = p.Lang
	}
	var out SendOtpResult
	if err := c.request(ctx, "/api/send-otp", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// VerifyOtp verifies a code the user entered.
func (c *Client) VerifyOtp(ctx context.Context, p VerifyOtpParams) (*VerifyOtpResult, error) {
	if p.Phone == "" && p.Email == "" {
		return nil, fmt.Errorf("otpsanly: VerifyOtp requires either Phone or Email")
	}
	if p.Code == "" {
		return nil, fmt.Errorf("otpsanly: VerifyOtp requires Code")
	}
	body := map[string]any{"code": p.Code}
	if p.Phone != "" {
		body["phone"] = p.Phone
	}
	if p.Email != "" {
		body["email"] = p.Email
	}
	if p.OtpID != 0 {
		body["otpId"] = p.OtpID
	}
	var out VerifyOtpResult
	if err := c.request(ctx, "/api/verify-otp", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
