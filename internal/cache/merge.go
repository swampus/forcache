package cache

import "fmt"

// MergeEngine defines how branches are merged into confirmed state.
type MergeEngine interface {
    CommitBranch(branchID string, confirmed *InMemoryConfirmedStore, spec *InMemorySpeculativeStore) error
    RollbackBranch(branchID string, spec *InMemorySpeculativeStore) error
}

// SimpleMergeEngine implements a basic merge policy:
// branch can commit only if confirmed version has not changed
// since branch.BaseVersion.
type SimpleMergeEngine struct{}

func NewSimpleMergeEngine() *SimpleMergeEngine {
    return &SimpleMergeEngine{}
}

func (m *SimpleMergeEngine) CommitBranch(branchID string, confirmed *InMemoryConfirmedStore, spec *InMemorySpeculativeStore) error {
    branch, ok := spec.getBranch(branchID)
    if !ok {
        return fmt.Errorf("branch not found: %s", branchID)
    }
    currentVersion := confirmed.getCurrentVersion()
    if currentVersion != branch.BaseVersion {
        branch.State = BranchRejected
        return fmt.Errorf("branch %s rejected: version conflict", branchID)
    }
    confirmed.commitOverlay(branch)
    branch.State = BranchCommitted
    spec.deleteBranch(branchID)
    return nil
}

func (m *SimpleMergeEngine) RollbackBranch(branchID string, spec *InMemorySpeculativeStore) error {
    branch, ok := spec.getBranch(branchID)
    if !ok {
        return fmt.Errorf("branch not found: %s", branchID)
    }
    branch.State = BranchRejected
    spec.deleteBranch(branchID)
    return nil
}
