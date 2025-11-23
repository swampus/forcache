package cache

import "testing"

func BenchmarkConfirmedSet(b *testing.B){
    s := NewInMemoryConfirmedStore()
    for i := 0; i < b.N; i++ {
        s.SetForTest("k", i)
    }
}
