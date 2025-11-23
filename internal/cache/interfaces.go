package cache

// QueryOptions controls how reads see speculative branches.
type QueryOptions struct {
    // BranchID, if non-empty, enables overlay from that branch.
    BranchID string
}

// Cache is the high-level cache interface used by service layer.
type Cache interface {
    // PutSpec creates a speculative write for a key and returns branch ID.
    PutSpec(key string, value any) (string, error)

    // Get returns value either from confirmed state or with a branch overlay.
    Get(key string, opts QueryOptions) (any, error)

    // Commit merges branch into confirmed state if possible.
    Commit(branchID string) error

    // Rollback discards branch and all its speculative changes.
    Rollback(branchID string) error
}
