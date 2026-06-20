package public

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetCurrency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bank/currency" {
			t.Errorf("Expected path /bank/currency, got %s", r.URL.Path)
		}
		res := []CurrencyInfo{
			{CurrencyCodeA: 840, CurrencyCodeB: 980, Date: 1718880000, RateBuy: 40.5, RateSell: 41.0},
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient()
	client.base.BaseURL = server.URL

	currencies, err := client.GetCurrency(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(currencies) != 1 {
		t.Errorf("Expected 1 currency, got %d", len(currencies))
	}

	if currencies[0].CurrencyCodeA != 840 {
		t.Errorf("Expected currencyCodeA 840, got %d", currencies[0].CurrencyCodeA)
	}
}
