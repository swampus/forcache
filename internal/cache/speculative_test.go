package cache

import "testing"

func TestSpeculativeStore_CreateBranch(t *testing.T) {
    t.Parallel()

    s := NewInMemorySpeculativeStore()

    base := Version(5)
    b := s.createBranch(base)

    if b.BaseVersion != base {
        t.Fatalf("BaseVersion = %d, want %d", b.BaseVersion, base)
    }
    if b.State != BranchPending {
        t.Fatalf("State = %s, want %s", b.State, BranchPending)
    }
    if len(b.Overlay) != 0 {
        t.Fatalf("expected empty overlay on new branch")
    }

    // ensure branch is stored
    if _, ok := s.getBranch(b.ID); !ok {
        t.Fatalf("branch %q not found in store", b.ID)
    }
}

func TestSpeculativeStore_OverlayIsolationBetweenBranches(t *testing.T) {
    t.Parallel()

    s := NewInMemorySpeculativeStore()

    b1 := s.createBranch(1)
    b1.Overlay["a"] = 10

    b2 := s.createBranch(1)

    if _, ok := b2.Overlay["a"]; ok {
        t.Fatalf("branch overlay leaked between branches")
    }
}

func TestSpeculativeStore_DeleteBranch(t *testing.T) {
    t.Parallel()

    s := NewInMemorySpeculativeStore()

    b := s.createBranch(0)
    s.deleteBranch(b.ID)

    if _, ok := s.getBranch(b.ID); ok {
        t.Fatalf("branch %q still present after delete", b.ID)
    }
}
