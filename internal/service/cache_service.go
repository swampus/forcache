package service

import (
    "fmt"

    "github.com/swampus/forcache/internal/cache"
)

// CacheService is an application-level facade over confirmed & speculative stores.
type CacheService struct {
    confirmed   *cache.InMemoryConfirmedStore
    speculative *cache.InMemorySpeculativeStore
    merge       cache.MergeEngine
}

func NewCacheService(
    confirmed *cache.InMemoryConfirmedStore,
    speculative *cache.InMemorySpeculativeStore,
    merge cache.MergeEngine,
) *CacheService {
    return &CacheService{
        confirmed:   confirmed,
        speculative: speculative,
        merge:       merge,
    }
}

// PutSpec creates a new branch and stores speculative value for a key.
func (s *CacheService) PutSpec(key string, value any) (string, error) {
    base := s.confirmed.getCurrentVersion()
    branch := s.speculative.createBranch(base)
    branch.Overlay[key] = value
    return branch.ID, nil
}

// Commit tries to merge branch into confirmed state via merge engine.
func (s *CacheService) Commit(branchID string) error {
    return s.merge.CommitBranch(branchID, s.confirmed, s.speculative)
}

// Rollback discards a branch via merge engine.
func (s *CacheService) Rollback(branchID string) error {
    return s.merge.RollbackBranch(branchID, s.speculative)
}

// Get returns value either from branch overlay (if provided) or confirmed state.
func (s *CacheService) Get(key string, opts cache.QueryOptions) (any, error) {
    if opts.BranchID != "" {
        if branch, ok := s.speculative.getBranch(opts.BranchID); ok {
            if v, ok := branch.Overlay[key]; ok {
                return v, nil
            }
        }
    }
    vv, ok := s.confirmed.get(key)
    if !ok {
        return nil, fmt.Errorf("key not found: %s", key)
    }
    return vv.Value, nil
}
