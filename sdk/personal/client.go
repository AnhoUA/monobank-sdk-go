// Package personal provides a Go SDK for the Monobank Personal API.
package personal

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/AnhoUA/monobank-sdk-go/sdk/internal/client"
)

// Client is the main client for interacting with the Monobank Personal API.
type Client struct {
	base *client.Client

	mu                 sync.Mutex
	lastStatementCall  time.Time
	lastClientInfoCall time.Time
}

// NewClient creates and returns a new Client instance with the given token and options.
func NewClient(token string, opts ...client.Option) *Client {
	return &Client{
		base: client.NewClient(token, opts...),
	}
}

// NewClientFromBase creates and returns a new Client instance from a base client.
func NewClientFromBase(base *client.Client) *Client {
	return &Client{base: base}
}

// UserInfo represents information about the client and their accounts.
type UserInfo struct {
	// Name of the client.
	Name string `json:"name"`
	// WebHookURL for receiving notifications.
	WebHookURL string `json:"webHookUrl"`
	// Accounts associated with the client.
	Accounts []Account `json:"accounts"`
}

// Account represents information about an individual account.
type Account struct {
	// ID is the unique identifier of the account.
	ID string `json:"id"`
	// Balance of the account in the smallest currency unit (e.g., cents).
	Balance int64 `json:"balance"`
	// CreditLimit of the account.
	CreditLimit int64 `json:"creditLimit"`
	// CurrencyCode is the ISO 4217 code of the account currency.
	CurrencyCode int32 `json:"currencyCode"`
	// CashbackType can be "None", "UAH", or "Miles".
	CashbackType string `json:"cashbackType"`
	// IBAN is the International Bank Account Number.
	IBAN string `json:"iban,omitzero"`
	// Type is the account type.
	Type string `json:"type,omitzero"`
}

// StatementItem represents a single transaction in the account statement.
type StatementItem struct {
	// ID is the unique identifier of the transaction.
	ID string `json:"id"`
	// Time of the transaction in Unix time.
	Time int64 `json:"time"`
	// Description of the transaction.
	Description string `json:"description"`
	// MCC is the Merchant Category Code.
	MCC int32 `json:"mcc"`
	// Hold indicates if the transaction is on hold.
	Hold bool `json:"hold"`
	// Amount of the transaction in the smallest currency unit.
	Amount int64 `json:"amount"`
	// OperationAmount is the amount in the original currency.
	OperationAmount int64 `json:"operationAmount"`
	// CurrencyCode is the ISO 4217 code of the transaction currency.
	CurrencyCode int32 `json:"currencyCode"`
	// CommissionRate is the commission charged for the transaction.
	CommissionRate int64 `json:"commissionRate"`
	// CashbackAmount is the cashback earned from the transaction.
	CashbackAmount int64 `json:"cashbackAmount"`
	// Balance after the transaction.
	Balance int64 `json:"balance"`
	// Comment added to the transaction.
	Comment string `json:"comment,omitzero"`
	// ReceiptID for the transaction.
	ReceiptID string `json:"receiptId,omitzero"`
	// CounterEdrpou is the EDRPOU code of the counterparty.
	CounterEdrpou string `json:"counterEdrpou,omitzero"`
	// CounterIban is the IBAN of the counterparty.
	CounterIban string `json:"counterIban,omitzero"`
}

// IsTooManyRequests checks if the error is a result of exceeding request limits.
func IsTooManyRequests(err error) bool {
	return client.IsTooManyRequests(err)
}

// GetClientInfo returns information about the authorized client and their accounts.
func (c *Client) GetClientInfo(ctx context.Context) (*UserInfo, error) {
	c.mu.Lock()
	if time.Since(c.lastClientInfoCall) < 60*time.Second {
		c.mu.Unlock()
		return nil, fmt.Errorf("too many requests: GetClientInfo can be called once per 60 seconds")
	}
	c.lastClientInfoCall = time.Now()
	c.mu.Unlock()

	var res UserInfo
	err := c.base.Do(ctx, http.MethodGet, "/personal/client-info", nil, &res, true)
	return &res, err
}

// SetWebhook sets the webhook URL for receiving transaction notifications.
func (c *Client) SetWebhook(ctx context.Context, webHookURL string) error {
	body := map[string]string{"webHookUrl": webHookURL}
	return c.base.Do(ctx, http.MethodPost, "/personal/webhook", body, nil, true)
}

// GetStatement returns a list of transactions for the specified account and time period.
// The maximum allowed period is 31 days + 1 hour.
func (c *Client) GetStatement(ctx context.Context, account string, from, to time.Time) ([]StatementItem, error) {
	c.mu.Lock()
	if time.Since(c.lastStatementCall) < 60*time.Second {
		c.mu.Unlock()
		return nil, fmt.Errorf("too many requests: GetStatement can be called once per 60 seconds")
	}
	c.lastStatementCall = time.Now()
	c.mu.Unlock()

	fromUnix := from.Unix()
	var toUnix int64
	if !to.IsZero() {
		toUnix = to.Unix()
		if toUnix-fromUnix > 2682000 {
			return nil, fmt.Errorf("max statement period is 31 days + 1 hour (2682000 seconds)")
		}
	}

	path := fmt.Sprintf("/personal/statement/%s/%d", account, fromUnix)
	if toUnix > 0 {
		path = fmt.Sprintf("%s/%d", path, toUnix)
	}
	var res []StatementItem
	err := c.base.Do(ctx, http.MethodGet, path, nil, &res, true)
	return res, err
}
