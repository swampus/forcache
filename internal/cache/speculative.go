package cache

import (
    "sync"
    "time"

    "github.com/google/uuid"
)

// InMemorySpeculativeStore keeps all active branches in memory.
type InMemorySpeculativeStore struct {
    mu       sync.RWMutex
    branches map[string]*Branch
}

func NewInMemorySpeculativeStore() *InMemorySpeculativeStore {
    return &InMemorySpeculativeStore{
        branches: make(map[string]*Branch),
    }
}

// createBranch creates a new branch with given base version.
func (s *InMemorySpeculativeStore) createBranch(base Version) *Branch {
    id := uuid.NewString()
    b := &Branch{
        ID:          id,
        BaseVersion: base,
        CreatedAt:   time.Now(),
        State:       BranchPending,
        Overlay:     make(map[string]Value),
    }
    s.mu.Lock()
    s.branches[id] = b
    s.mu.Unlock()
    return b
}

func (s *InMemorySpeculativeStore) getBranch(id string) (*Branch, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    b, ok := s.branches[id]
    return b, ok
}

func (s *InMemorySpeculativeStore) deleteBranch(id string) {
    s.mu.Lock()
    delete(s.branches, id)
    s.mu.Unlock()
}
