package httpapi

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/swampus/forcache/internal/cache"
    "github.com/swampus/forcache/internal/service"
)

func newTestServer(t *testing.T) *httptest.Server {
    t.Helper()

    confirmed := cache.NewInMemoryConfirmedStore()
    spec := cache.NewInMemorySpeculativeStore()
    merge := cache.NewSimpleMergeEngine()
    svc := service.NewCacheService(confirmed, spec, merge)

    router := NewRouter(svc)
    return httptest.NewServer(router)
}

func TestHTTP_FullSpeculativeFlow(t *testing.T) {
    t.Parallel()

    srv := newTestServer(t)
    defer srv.Close()

    // 1) create speculative branch
    body := map[string]any{"value": 123}
    buf, _ := json.Marshal(body)

    resp, err := http.Put(srv.URL+"/spec/foo", "application/json", bytes.NewReader(buf))
    if err != nil {
        t.Fatalf("PUT /spec error: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("status = %d, want 200", resp.StatusCode)
    }

    var putResp struct {
        BranchID string `json:"branch_id"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&putResp); err != nil {
        t.Fatalf("decode putResp error: %v", err)
    }
    if putResp.BranchID == "" {
        t.Fatalf("empty branch_id in response")
    }

    // 2) read via branch
    resp2, err := http.Get(srv.URL + "/value/foo?branch=" + putResp.BranchID)
    if err != nil {
        t.Fatalf("GET /value with branch error: %v", err)
    }
    defer resp2.Body.Close()

    if resp2.StatusCode != http.StatusOK {
        t.Fatalf("status = %d, want 200", resp2.StatusCode)
    }

    var valResp map[string]any
    if err := json.NewDecoder(resp2.Body).Decode(&valResp); err != nil {
        t.Fatalf("decode valueResp error: %v", err)
    }
    if valResp["value"] != float64(123) { // JSON numbers decode as float64
        t.Fatalf("value = %v, want 123", valResp["value"])
    }

    // 3) commit branch
    req, _ := http.NewRequest(http.MethodPost, srv.URL+"/commit/"+putResp.BranchID, nil)
    resp3, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("POST /commit error: %v", err)
    }
    defer resp3.Body.Close()

    if resp3.StatusCode != http.StatusNoContent {
        t.Fatalf("status = %d, want 204", resp3.StatusCode)
    }

    // 4) read from confirmed
    resp4, err := http.Get(srv.URL + "/value/foo")
    if err != nil {
        t.Fatalf("GET /value confirmed error: %v", err)
    }
    defer resp4.Body.Close()

    if resp4.StatusCode != http.StatusOK {
        t.Fatalf("status = %d, want 200", resp4.StatusCode)
    }

    var valResp2 map[string]any
    if err := json.NewDecoder(resp4.Body).Decode(&valResp2); err != nil {
        t.Fatalf("decode valueResp2 error: %v", err)
    }
    if valResp2["value"] != float64(123) {
        t.Fatalf("confirmed value = %v, want 123", valResp2["value"])
    }
}
