package httpapi

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/swampus/forcache/internal/service"
)

// NewRouter wires HTTP handlers.
func NewRouter(cacheSvc *service.CacheService) http.Handler {
    h := &Handler{svc: cacheSvc}
    r := mux.NewRouter()

    r.HandleFunc("/spec/{key}", h.handlePutSpec).Methods(http.MethodPut)
    r.HandleFunc("/commit/{branchID}", h.handleCommit).Methods(http.MethodPost)
    r.HandleFunc("/rollback/{branchID}", h.handleRollback).Methods(http.MethodPost)
    r.HandleFunc("/value/{key}", h.handleGet).Methods(http.MethodGet)

    return r
}
