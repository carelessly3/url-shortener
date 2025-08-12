package storage

import (
	"testing"
	"time"
)

func TestMemorySaveGet(t *testing.T) {
	s := NewMemoryStore()
	now := time.Now().UTC()
	err := s.Save("abc123", "https://example.com", &now, false)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	rec, ok := s.Get("abc123")
	if !ok {
		t.Fatal("expected to find code")
	}
	if rec.LongURL != "https://example.com" {
		t.Fatalf("unexpected url: %s", rec.LongURL)
	}
	s.IncrementClick("abc123")
	rec2, _ := s.Get("abc123")
	if rec2.Clicks != 1 {
		t.Fatalf("expected clicks 1 got %d", rec2.Clicks)
	}
}
