package mcc

import (
	"testing"
)

func TestGet(t *testing.T) {
	// Test English (default)
	info, ok := Get("0742")
	if !ok {
		t.Fatal("expected to find MCC 0742")
	}
	if info.MCC != "0742" {
		t.Errorf("expected MCC 0742, got %s", info.MCC)
	}
	if info.ShortDescription != "Veterinary Services" {
		t.Errorf("expected 'Veterinary Services', got '%s'", info.ShortDescription)
	}

	// Test Ukrainian
	info, ok = GetUk("0742")
	if !ok {
		t.Fatal("expected to find MCC 0742 in Ukrainian")
	}
	if info.ShortDescription != "Ветеринарні послуги" {
		t.Errorf("expected 'Ветеринарні послуги', got '%s'", info.ShortDescription)
	}
	if info.GroupName != "Сільськогосподарські послуги" {
		t.Errorf("expected 'Сільськогосподарські послуги', got '%s'", info.GroupName)
	}

	// Test non-existent
	_, ok = Get("none")
	if ok {
		t.Error("expected not to find MCC 'none'")
	}
}

func TestAll(t *testing.T) {
	all, err := All()
	if err != nil {
		t.Fatalf("All() failed: %v", err)
	}
	if len(all) == 0 {
		t.Fatal("expected All() to return records")
	}
}
