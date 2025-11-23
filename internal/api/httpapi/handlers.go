package httpapi

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/swampus/forcache/internal/cache"
    "github.com/swampus/forcache/internal/service"
)

// Handler holds dependencies for HTTP layer.
type Handler struct {
    svc *service.CacheService
}

type putSpecRequest struct {
    Value any `json:"value"`
}

type putSpecResponse struct {
    BranchID string `json:"branch_id"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) handlePutSpec(w http.ResponseWriter, r *http.Request) {
    key := mux.Vars(r)["key"]
    var req putSpecRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }

    id, err := h.svc.PutSpec(key, req.Value)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    writeJSON(w, http.StatusOK, putSpecResponse{BranchID: id})
}

func (h *Handler) handleCommit(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["branchID"]
    if err := h.svc.Commit(id); err != nil {
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleRollback(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["branchID"]
    if err := h.svc.Rollback(id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
    key := mux.Vars(r)["key"]
    branch := r.URL.Query().Get("branch")

    val, err := h.svc.Get(key, cache.QueryOptions{BranchID: branch})
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{
        "key":   key,
        "value": val,
    })
}
