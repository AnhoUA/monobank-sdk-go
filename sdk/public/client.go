// Package public provides a Go SDK for the Monobank Public API.
package public

import (
	"context"
	"net/http"

	"github.com/AnhoUA/monobank-sdk-go/sdk/internal/client"
)

// Client is the main client for interacting with the Monobank Public API.
type Client struct {
	base *client.Client
}

// NewClient creates and returns a new Client instance for public API.
func NewClient(opts ...client.Option) *Client {
	return &Client{
		base: client.NewClient("", opts...),
	}
}

// NewClientFromBase creates and returns a new Client instance from a base client.
func NewClientFromBase(base *client.Client) *Client {
	return &Client{base: base}
}

// CurrencyInfo represents information about currency rates.
type CurrencyInfo struct {
	// CurrencyCodeA is the ISO 4217 code of the first currency.
	CurrencyCodeA int32 `json:"currencyCodeA"`
	// CurrencyCodeB is the ISO 4217 code of the second currency.
	CurrencyCodeB int32 `json:"currencyCodeB"`
	// Date of the rate in Unix time.
	Date int64 `json:"date"`
	// RateSell is the selling rate.
	RateSell float64 `json:"rateSell,omitzero"`
	// RateBuy is the buying rate.
	RateBuy float64 `json:"rateBuy,omitzero"`
	// RateCross is the cross rate.
	RateCross float64 `json:"rateCross,omitzero"`
}

// GetCurrency returns a list of current currency rates from Monobank.
func (c *Client) GetCurrency(ctx context.Context) ([]CurrencyInfo, error) {
	var res []CurrencyInfo
	err := c.base.Do(ctx, http.MethodGet, "/bank/currency", nil, &res, false)
	return res, err
}
