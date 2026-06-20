// Package sdk provides a unified Go SDK for the Monobank API.
package sdk

import (
	"github.com/AnhoUA/monobank-sdk-go/sdk/internal/client"
	"github.com/AnhoUA/monobank-sdk-go/sdk/mcc"
	"github.com/AnhoUA/monobank-sdk-go/sdk/personal"
	"github.com/AnhoUA/monobank-sdk-go/sdk/public"
)

// Client is the unified client for interacting with the Monobank API.
type Client struct {
	Public   *public.Client
	Personal *personal.Client
	MCC      *mcc.MCCService
}

// NewClient creates and returns a new unified Client instance.
func NewClient(token string, opts ...client.Option) *Client {
	base := client.NewClient(token, opts...)

	return &Client{
		Public:   public.NewClientFromBase(base),
		Personal: personal.NewClientFromBase(base),
		MCC:      &mcc.MCCService{},
	}
}
