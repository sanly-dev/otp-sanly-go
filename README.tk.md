# OTP Sanly — Go SDK

[Türkmençe](README.tk.md) | [Русский](README.ru.md) | [English](README.md)

**OTP Sanly** (https://otp.sanly.dev) üçin resmi SDK — Türkmenistan üçin SMS & Email OTP tassyklama hyzmaty.

> **Bellik:** registr-esasly ekosistemalardan (npm, PyPI) tapawutlylykda, Go modullary
> göni şu repo-dan alynýar, şoň üçin `go get` işlemegi üçin repo hemişe açyk (public) bolmaly.

## Gurnamak

```bash
go get github.com/sanly-dev/otp-sanly-go
```

## Çalt başlamak

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

	// OTP iber
	sent, err := client.SendOtp(context.Background(), otpsanly.SendOtpParams{
		Phone:   "+99361234567", // SMS üçin diňe Türkmenistan belgileri. Dünýäniň islendik ýerine ibermek üçin Email ulanyň.
		Project: "Meniň Programmam",
		Lang:    "ru", // "tm" | "ru" | "en" — OTP-iň haýsy dilde ugradylmalydygy
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.OtpID, sent.Message)

	// Ulanyjynyň girizen kodyny barla
	verified, err := client.VerifyOtp(context.Background(), otpsanly.VerifyOtpParams{
		Phone: "+99361234567",
		Code:  "123456",
	})
	if err != nil {
		panic(err)
	}
	if verified.Success {
		// dowam et — kod dogry boldy
	}
}
```

## Synag (sandbox) rejimi

Integrasiýaňyzy hakyky SMS/email ibermezden, mugt synamak üçin [dolandyryş panelinde](https://otp.sanly.dev/dashboard/api-keys)
**sandbox** API açary dörediň. OTP kody göni jogapda gaýtarylýar:

```go
sent, _ := client.SendOtp(context.Background(), otpsanly.SendOtpParams{Phone: "+99361234567"})
fmt.Println(sent.Code) // diňe sandbox açarlar üçin bar
```

## Ýalňyşlyklary dolandyrmak

Şowsuz soraglar `*otpsanly.APIError` gaýtarýar, ol HTTP statusyny we çig jogap gövresini öz içine alýar:

```go
sent, err := client.SendOtp(context.Background(), otpsanly.SendOtpParams{Phone: "+99361234567"})
if err != nil {
	var apiErr *otpsanly.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.Status, apiErr.Message, apiErr.Body)
	}
}
```

## Doly API salgylanmasy

Doly API salgylanmasy, webhook resminamalary we beýleki dillerdäki mysallar üçin
[https://otp.sanly.dev/developers](https://otp.sanly.dev/developers) sahypasyna serediň.

## Ygtyýarnama

MIT
