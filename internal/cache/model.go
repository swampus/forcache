package cache

import "time"

// Value represents stored value. For V0 it's kept as `any`.
type Value = any

// Version is a monotonic version of confirmed state.
type Version uint64

// VersionedValue stores value with version metadata.
type VersionedValue struct {
    Value   Value
    Version Version
}

// BranchState describes lifecycle of a branch.
type BranchState string

const (
    BranchPending   BranchState = "PENDING"
    BranchCommitted BranchState = "COMMITTED"
    BranchRejected  BranchState = "REJECTED"
)

// Branch represents speculative overlay for multiple keys.
type Branch struct {
    ID          string
    BaseVersion Version
    CreatedAt   time.Time
    State       BranchState
    Overlay     map[string]Value
}
