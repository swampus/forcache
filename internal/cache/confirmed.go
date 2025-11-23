package cache

import "sync"

// InMemoryConfirmedStore is a simple thread-safe confirmed layer.
type InMemoryConfirmedStore struct {
    mu      sync.RWMutex
    data    map[string]VersionedValue
    version Version
}

func NewInMemoryConfirmedStore() *InMemoryConfirmedStore {
    return &InMemoryConfirmedStore{
        data:    make(map[string]VersionedValue),
        version: 0,
    }
}

// getCurrentVersion returns current global version.
func (s *InMemoryConfirmedStore) getCurrentVersion() Version {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.version
}

// get returns confirmed value (without speculative overlay).
func (s *InMemoryConfirmedStore) get(key string) (VersionedValue, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.data[key]
    return v, ok
}

// setUnsafe sets value and bumps global version. Caller must hold write lock.
func (s *InMemoryConfirmedStore) setUnsafe(key string, value Value) {
    s.version++
    s.data[key] = VersionedValue{
        Value:   value,
        Version: s.version,
    }
}

// commitOverlay merges all keys from a branch overlay into confirmed state.
// This is used by merge engine when a branch is accepted.
func (s *InMemoryConfirmedStore) commitOverlay(branch *Branch) {
    s.mu.Lock()
    defer s.mu.Unlock()
    for k, v := range branch.Overlay {
        s.setUnsafe(k, v)
    }
}
