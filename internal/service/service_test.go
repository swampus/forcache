package service

import (
    "testing"
    "github.com/swampus/forcache/internal/cache"
)

func TestService_Table(t *testing.T){
    items := []any{10,20,30}

    for _, v := range items {
        t.Run("case", func(t *testing.T){
            c := cache.NewInMemoryConfirmedStore()
            s := cache.NewInMemorySpeculativeStore()
            m := cache.NewSimpleMergeEngine()
            svc := NewCacheService(c, s, m)

            br, _ := svc.PutSpec("x", v)
            got, _ := svc.Get("x", cache.QueryOptions{BranchID: br})
            if got != v {
                t.Fatalf("expected %v got %v", v, got)
            }
        })
    }
}
