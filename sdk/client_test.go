package sdk

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	token := "test-token"
	client := NewClient(token)

	if client.Personal == nil {
		t.Error("Personal client is nil")
	}
	if client.Public == nil {
		t.Error("Public client is nil")
	}
}
