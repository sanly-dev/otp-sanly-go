# OTP Sanly — Go SDK

[Türkmençe](README.tk.md) | Русский | [English](README.md)

Официальный SDK для [OTP Sanly](https://otp.sanly.dev) — сервис SMS и Email OTP-аутентификации для Туркменистана.

> **Примечание:** в отличие от экосистем на базе реестров (npm, PyPI), модули Go загружаются
> напрямую из этого репозитория, поэтому он должен быть публичным, чтобы `go get` работал.

## Установка

```bash
go get github.com/sanly-dev/otp-sanly-go
```

## Быстрый старт

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

	// Отправить OTP
	sent, err := client.SendOtp(context.Background(), otpsanly.SendOtpParams{
		Phone:   "+99361234567", // Для SMS только номера Туркменистана. Для доставки по всему миру используйте Email.
		Project: "Моё приложение",
		Lang:    "ru", // "tm" | "ru" | "en" — на каком языке отправить OTP
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.OtpID, sent.Message)

	// Проверить код, введённый пользователем
	verified, err := client.VerifyOtp(context.Background(), otpsanly.VerifyOtpParams{
		Phone: "+99361234567",
		Code:  "123456",
	})
	if err != nil {
		panic(err)
	}
	if verified.Success {
		// продолжайте — код верный
	}
}
```

## Тестовый (sandbox) режим

Создайте **тестовый (sandbox)** API-ключ в [панели управления](https://otp.sanly.dev/dashboard/api-keys), чтобы
бесплатно протестировать интеграцию без отправки реальных SMS/email. Код OTP возвращается прямо в ответе:

```go
sent, _ := client.SendOtp(context.Background(), otpsanly.SendOtpParams{Phone: "+99361234567"})
fmt.Println(sent.Code) // доступно только для тестовых ключей
```

## Обработка ошибок

Неудачные запросы возвращают `*otpsanly.APIError`, содержащий HTTP-статус и тело ответа:

```go
sent, err := client.SendOtp(context.Background(), otpsanly.SendOtpParams{Phone: "+99361234567"})
if err != nil {
	var apiErr *otpsanly.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.Status, apiErr.Message, apiErr.Body)
	}
}
```

## Полная документация по API

Полную документацию по API, описание webhook и примеры на других языках см. на странице
[https://otp.sanly.dev/developers](https://otp.sanly.dev/developers).

## Лицензия

MIT
