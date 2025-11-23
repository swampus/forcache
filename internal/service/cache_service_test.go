package service

import (
    "testing"

    "github.com/swampus/forcache/internal/cache"
)

func TestCacheService_FullSpeculativeFlow(t *testing.T) {
    t.Parallel()

    confirmed := cache.NewInMemoryConfirmedStore()
    spec := cache.NewInMemorySpeculativeStore()
    merge := cache.NewSimpleMergeEngine()

    s := NewCacheService(confirmed, spec, merge)

    // STEP 1: speculative write
    branchID, err := s.PutSpec("a", 10)
    if err != nil {
        t.Fatalf("PutSpec() error = %v", err)
    }

    // read inside branch
    got, err := s.Get("a", cache.QueryOptions{BranchID: branchID})
    if err != nil {
        t.Fatalf("Get() in branch error = %v", err)
    }
    if got != 10 {
        t.Fatalf("Get() in branch = %v, want 10", got)
    }

    // read without branch - should be not found
    if _, err := s.Get("a", cache.QueryOptions{}); err == nil {
        t.Fatalf("expected error when reading key outside branch before commit")
    }

    // STEP 2: commit
    if err := s.Commit(branchID); err != nil {
        t.Fatalf("Commit() error = %v", err)
    }

    // STEP 3: read from confirmed
    got2, err := s.Get("a", cache.QueryOptions{})
    if err != nil {
        t.Fatalf("Get() after commit error = %v", err)
    }
    if got2 != 10 {
        t.Fatalf("Get() after commit = %v, want 10", got2)
    }
}

func TestCacheService_CommitConflict(t *testing.T) {
    t.Parallel()

    confirmed := cache.NewInMemoryConfirmedStore()
    spec := cache.NewInMemorySpeculativeStore()
    merge := cache.NewSimpleMergeEngine()

    s := NewCacheService(confirmed, spec, merge)

    // speculative branch on base version 0
    branchID, err := s.PutSpec("a", 10)
    if err != nil {
        t.Fatalf("PutSpec() error = %v", err)
    }

    // external change -> confirmed version bump
    confirmed.SetForTest("a", 5)

    // commit should now conflict
    if err := s.Commit(branchID); err == nil {
        t.Fatalf("expected conflict error on Commit(), got nil")
    }
}
