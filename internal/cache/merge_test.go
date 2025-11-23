package cache
import "testing"

func TestMerge_RollbackTable(t *testing.T){
    confirmed := NewInMemoryConfirmedStore()
    spec := NewInMemorySpeculativeStore()
    m := NewSimpleMergeEngine()

    versions := []int{0,1,2}

    for _, base := range versions {
        t.Run("v", func(t *testing.T){
            br := spec.createBranch(confirmed.getCurrentVersion())
            br.Overlay["x"] = base
            _ = m.RollbackBranch(br.ID, spec)
        })
    }
}
