# OTP Sanly — Go SDK

[Türkmençe](README.tk.md) | [Русский](README.ru.md) | English

Official SDK for [OTP Sanly](https://otp.sanly.dev) — SMS & Email OTP authentication for Turkmenistan.

> **Note:** unlike registry-based ecosystems (npm, PyPI), Go modules are fetched directly from
> this repository, so it must be public for `go get` to work.

## Install

```bash
go get github.com/sanly-dev/otp-sanly-go
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"os"

	otpsanly "github.com/sanly-dev/otp-sanly-go"
)

func main() {
	client := otpsanly.NewClient(os.Getenv("OTP_API_KEY"))

	// Send an OTP
	sent, err := client.SendOtp(context.Background(), otpsanly.SendOtpParams{
		Phone:   "+99361234567", // Turkmenistan numbers only for SMS. Use Email instead for worldwide delivery.
		Project: "My App",
		Lang:    "ru", // "tm" | "ru" | "en" — which language to send the OTP in
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.OtpID, sent.Message)

	// Verify the code the user entered
	verified, err := client.VerifyOtp(context.Background(), otpsanly.VerifyOtpParams{
		Phone: "+99361234567",
		Code:  "123456",
	})
	if err != nil {
		panic(err)
	}
	if verified.Success {
		// proceed — the code was correct
	}
}
```

## Sandbox mode

Create a **sandbox** API key in your [dashboard](https://otp.sanly.dev/dashboard/api-keys) to test your integration
for free, without sending real SMS/email. The OTP code is returned directly in the response:

```go
sent, _ := client.SendOtp(context.Background(), otpsanly.SendOtpParams{Phone: "+99361234567"})
fmt.Println(sent.Code) // only present for sandbox keys
```

## Error handling

Failed requests return an `*otpsanly.APIError`, which includes the HTTP status and the raw response body:

```go
sent, err := client.SendOtp(context.Background(), otpsanly.SendOtpParams{Phone: "+99361234567"})
if err != nil {
	var apiErr *otpsanly.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.Status, apiErr.Message, apiErr.Body)
	}
}
```

## Full API reference

See [https://otp.sanly.dev/developers](https://otp.sanly.dev/developers) for the complete API reference, webhook docs,
and framework-specific examples in other languages.

## License

MIT
