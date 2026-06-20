package personal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AnhoUA/monobank-sdk-go/sdk/internal/client"
)

func TestClient_GetClientInfo(t *testing.T) {
	token := "test-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != token {
			t.Errorf("Expected X-Token header %s, got %s", token, r.Header.Get("X-Token"))
		}
		res := UserInfo{
			Name: "Ivan Ivanov",
			Accounts: []Account{
				{ID: "acc1", Balance: 1000},
			},
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient(token)
	client.base.BaseURL = server.URL

	info, err := client.GetClientInfo(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if info.Name != "Ivan Ivanov" {
		t.Errorf("Expected name Ivan Ivanov, got %s", info.Name)
	}
}

func TestClient_GetStatement(t *testing.T) {
	token := "test-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expected path: /personal/statement/0/1718880000/1718883600
		expectedPath := "/personal/statement/0/1718880000/1718883600"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		res := []StatementItem{
			{ID: "stmt1", Amount: -100, Time: 1718881000},
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient(token)
	client.base.BaseURL = server.URL

	from := time.Unix(1718880000, 0)
	to := time.Unix(1718883600, 0)

	stmt, err := client.GetStatement(context.Background(), "0", from, to)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(stmt) != 1 {
		t.Errorf("Expected 1 statement item, got %d", len(stmt))
	}
}

func TestClient_GetStatement_LimitError(t *testing.T) {
	client := NewClient("token")
	from := time.Unix(1000000000, 0)
	to := time.Unix(1000000000+2682001, 0) // Limit + 1 second

	_, err := client.GetStatement(context.Background(), "0", from, to)
	if err == nil {
		t.Fatal("Expected error due to time limit, got nil")
	}

	expectedErr := "max statement period is 31 days + 1 hour (2682000 seconds)"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestClient_GetClientInfo_Caching(t *testing.T) {
	token := "test-token"
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		res := UserInfo{
			Name: fmt.Sprintf("Call %d", calls),
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	c := NewClient(token)
	c.base.BaseURL = server.URL

	// First call
	info1, err := c.GetClientInfo(context.Background())
	if err != nil {
		t.Fatalf("First call failed: %v", err)
	}
	if info1.Name != "Call 1" {
		t.Errorf("Expected 'Call 1', got '%s'", info1.Name)
	}
	if calls != 1 {
		t.Errorf("Expected 1 server call, got %d", calls)
	}

	// Second call immediately after - should be cached
	info2, err := c.GetClientInfo(context.Background())
	if err != nil {
		t.Fatalf("Second call failed: %v", err)
	}
	if info2.Name != "Call 1" {
		t.Errorf("Expected cached 'Call 1', got '%s'", info2.Name)
	}
	if calls != 1 {
		t.Errorf("Expected still 1 server call, got %d", calls)
	}

	// Manipulate last call time to simulate expiration
	c.mu.Lock()
	c.lastClientInfoCall = time.Now().Add(-61 * time.Second)
	c.mu.Unlock()

	// Third call after expiration - should call server again
	info3, err := c.GetClientInfo(context.Background())
	if err != nil {
		t.Fatalf("Third call failed: %v", err)
	}
	if info3.Name != "Call 2" {
		t.Errorf("Expected 'Call 2', got '%s'", info3.Name)
	}
	if calls != 2 {
		t.Errorf("Expected 2 server calls, got %d", calls)
	}
}

func TestClient_RetryOn429(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"errorDescription":"Too many requests"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"Test"}`))
	}))
	defer server.Close()

	// Setting a short interval for the test
	client := NewClient("token", client.WithRetryInterval(10*time.Millisecond))
	client.base.BaseURL = server.URL

	// Since we are testing GetClientInfo, we check retry logic
	_, err := client.GetClientInfo(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}
