package cache
import "testing"

func TestConfirmed_Table(t *testing.T){
    tests := []struct{
        key string
        val any
    }{
        {"a", 1},
        {"b", 2},
        {"c", 3},
    }

    for _, tt := range tests {
        t.Run(tt.key, func(t *testing.T){
            s := NewInMemoryConfirmedStore()
            s.setUnsafe(tt.key, tt.val)
            vv, ok := s.get(tt.key)
            if !ok || vv.Value != tt.val {
                t.Fatalf("expected %v got %v", tt.val, vv.Value)
            }
        })
    }
}
