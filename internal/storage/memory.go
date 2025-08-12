package storage

import (
	"errors"
	"sync"
	"time"
)

// URLRecord stores link metadata
type URLRecord struct {
	LongURL   string
	CreatedAt time.Time
	ExpiresAt *time.Time // nil means no expiry
	Clicks    int64
}

// MemoryStore is a simple in-memory thread-safe store
type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]URLRecord
}

// NewMemoryStore returns an initialized MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]URLRecord),
	}
}

// Save stores a mapping code -> longURL. If overwrite is false and code exists, returns error.
func (m *MemoryStore) Save(code, longURL string, expiresAt *time.Time, overwrite bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !overwrite {
		if _, ok := m.store[code]; ok {
			return errors.New("code already exists")
		}
	}
	m.store[code] = URLRecord{
		LongURL:   longURL,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: expiresAt,
		Clicks:    0,
	}
	return nil
}

// Get returns the long URL and boolean found
func (m *MemoryStore) Get(code string) (URLRecord, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.store[code]
	return r, ok
}

// IncrementClick increases click count for a code
func (m *MemoryStore) IncrementClick(code string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.store[code]
	if !ok {
		return
	}
	r.Clicks++
	m.store[code] = r
}
