package main

import (
    "log"
    "net/http"

    "github.com/swampus/forcache/internal/api/httpapi"
    "github.com/swampus/forcache/internal/cache"
    "github.com/swampus/forcache/internal/service"
)

func main() {
    confirmed := cache.NewInMemoryConfirmedStore()
    spec := cache.NewInMemorySpeculativeStore()
    mergeEngine := cache.NewSimpleMergeEngine()

    cacheSvc := service.NewCacheService(confirmed, spec, mergeEngine)
    router := httpapi.NewRouter(cacheSvc)

    log.Println("forcache server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
