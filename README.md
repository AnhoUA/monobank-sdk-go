# Monobank SDK for Go

[Українська](#українська) | [English](#english)

---

## Українська

Go SDK для взаємодії з Monobank Personal API. Ця бібліотека дозволяє легко інтегрувати функції Monobank у ваші Go-проекти.

### Особливості

- **Єдиний клієнт**: доступ до всіх функцій через одну точку входу (`client.Public` та `client.Personal`).
- **Повна підтримка Personal API**: курси валют, інформація про клієнта, виписки.
- **Типізація**: використання `time.Time` для дат та чітких структур для даних.
- **Rate Limiting**: вбудований захист від занадто частих запитів (60с ліміт для виписок та інфо).
- **Автоматичні повтори (Retries)**: підтримка повторних спроб при отриманні HTTP 429 (Too Many Requests).
- **Конфігурація**: можливість налаштування власного HTTP-клієнта та інтервалів ретраїв.

### Встановлення

```bash
go get github.com/AnhoUA/monobank-sdk-go
```

### Приклад використання

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AnhoUA/monobank-sdk-go/sdk"
	"github.com/AnhoUA/monobank-sdk-go/sdk/internal/client"
)

func main() {
	// Створення єдиного клієнта з налаштуваннями
	c := sdk.NewClient("YOUR_TOKEN", client.WithRetryInterval(2 * time.Second))

	// Публічні дані
	currencies, err := c.Public.GetCurrency(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range currencies {
		fmt.Printf("Валюта: %d, Купівля: %.2f\n", c.CurrencyCodeA, c.RateBuy)
	}

	// Персональні дані
	info, err := c.Personal.GetClientInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Привіт, %s!\n", info.Name)
}
```

---

## English

Go SDK for interacting with Monobank Personal API. This library allows you to easily integrate Monobank features into your Go projects.

### Features

- **Unified Client**: access all features via a single entry point (`client.Public` and `client.Personal`).
- **Full Personal API Support**: currency rates, client info, statements.
- **Strong Typing**: uses `time.Time` for dates and clear data structures.
- **Rate Limiting**: built-in protection against frequent requests (60s limit for statements and info).
- **Automatic Retries**: supports retrying requests when receiving HTTP 429 (Too Many Requests).
- **Configuration**: ability to set a custom HTTP client and retry intervals.

### Installation

```bash
go get github.com/AnhoUA/monobank-sdk-go
```

### Usage Example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AnhoUA/monobank-sdk-go/sdk"
	"github.com/AnhoUA/monobank-sdk-go/sdk/internal/client"
)

func main() {
	// Create unified client with settings
	c := sdk.NewClient("YOUR_TOKEN", client.WithRetryInterval(2 * time.Second))

	// Public data
	currencies, err := c.Public.GetCurrency(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range currencies {
		fmt.Printf("Currency: %d, Buy: %.2f\n", c.CurrencyCodeA, c.RateBuy)
	}

	// Personal data
	info, err := c.Personal.GetClientInfo(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Hello, %s!\n", info.Name)
}
```
