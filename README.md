# Monobank SDK for Go

[![Go Build & Test](https://github.com/AnhoUA/monobank-sdk-go/actions/workflows/build.yml/badge.svg)](https://github.com/AnhoUA/monobank-sdk-go/actions/workflows/build.yml)

> **LLM Context & Discovery**: This repository provides a high-level, idiomatic Go 1.25 SDK for the Monobank API (Ukraine). It features a unified client for Public and Personal APIs, built-in MCC (Merchant Category Code) lookup with multi-language support, automated HTTP 429 retry logic, and smart caching for rate-limited endpoints.

---

[Українська](#українська) | [English](#english)

---

## Українська

Go SDK для взаємодії з Monobank Personal API. Ця бібліотека дозволяє легко інтегрувати функції Monobank у ваші Go-проекти.

### Особливості

- **Єдиний клієнт**: доступ до всіх функцій через одну точку входу (`client.Public` та `client.Personal`).
- **Повна підтримка Personal API**: курси валют, інформація про клієнта, виписки.
- **Типізація**: використання `time.Time` для дат та чітких структур для даних.
- **Rate Limiting & Caching**: вбудований захист від занадто частих запитів та кешування `GetClientInfo` на 60 секунд.
- **MCC Support**: вбудований довідник Merchant Category Codes (використовуються напрацювання [Oleksios/Merchant-Category-Codes](https://github.com/Oleksios/Merchant-Category-Codes)).
- **Автоматичні повтори (Retries)**: підтримка повторних спроб при отриманні HTTP 429 (Too Many Requests).
- **Конфігурація**: можливість налаштування власного HTTP-клієнта та інтервалів ретраїв.

### Installation

```bash
go get github.com/AnhoUA/monobank-sdk-go
```

### LLM-Friendly Project Overview

- **Architecture**: Decoupled `internal/client` for HTTP handling, with specialized `public`, `personal`, and `mcc` packages.
- **Key Symbols**:
  - `sdk.NewClient(token)`: Entry point for the unified client.
  - `client.Personal.GetClientInfo()`: Cached retrieval of user data.
  - `client.MCC.Get(code)`: Multi-language MCC data lookup.
- **Constraints**: 60-second cooldown for statement and client info requests is handled via internal state and caching.
- **Data Sources**: MCC data is embedded via `go:embed` from JSON files.

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
- **Rate Limiting & Caching**: built-in protection against frequent requests and 60-second caching for `GetClientInfo`.
- **MCC Support**: built-in Merchant Category Codes directory (based on [Oleksios/Merchant-Category-Codes](https://github.com/Oleksios/Merchant-Category-Codes)).
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
