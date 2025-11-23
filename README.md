
# Forcache – Speculative Branching Cache

Forcache is an experimental high-performance in-memory cache that supports **speculative branching**, allowing clients to perform isolated “what-if” writes on temporary branches and later **commit** or **rollback** them based on business logic.

This design combines ideas from:
- optimistic concurrency control,
- lightweight versioned storage,
- MVCC inspiration without overhead,
- speculative execution models.

The goal is to provide a **simple, explicit and controllable branching model** suitable for edge computation, temporary workflows, simulations and local decision engines.

---
## Documentation

- **Design Specification:** [DESIGN.md](./DESIGN.md)  
  Comprehensive description of the system architecture, versioning model, branching semantics, merge logic and future extensions.

- **Testing Strategy:** [TESTING.md](./TESTING.md)  
  Detailed overview of the test approach (unit, integration, API, benchmarks)

## Articles

- [Speculative Branching Cache: Managing Temporary State Without Transactions or Raft](./ARTICLE.md)


## Key Features

### ✔ Speculative Branches
Create independent temporary overlays on top of confirmed state.

### ✔ Version-Based Commit
Commit succeeds only if the global version hasn't changed.

### ✔ Clean Merge Engine
Explicit conflict model: match version → commit; mismatch → reject.

### ✔ Minimalistic Consistency Model
Simple enough to reason about, strong enough for real workflows.

### ✔ HTTP API + Service Layer
Full end-to-end functional example with REST endpoints.

### ✔ Tests & Benchmarks
Full test suite with:
- unit tests,
- integration tests,
- table-driven tests,
- benchmarks.

---

## Architecture Overview

```
                  +---------------------+
                  |   CacheService      |
                  |  (API Facade Layer) |
                  +----------+----------+
                             |
        -------------------------------------------------
        |                                               |
+-------v-------+                               +-------v-------+
| Speculative   |                               | Confirmed     |
| Store         |                               | Store         |
| - branches    |                               | - KV state    |
| - overlays    |                               | - version     |
+-------+-------+                               +-------+-------+
        |                                               |
        +----------------------+------------------------+
                               |
                        +------v-------+
                        | Merge Engine |
                        | conflict/ok  |
                        +--------------+
```

---

## HTTP API

### Create speculative write
```
PUT /spec/{key}
{
  "value": 123
}
```

### Read
Confirmed:
```
GET /value/foo
```
With branch:
```
GET /value/foo?branch={id}
```

### Commit
```
POST /commit/{branch_id}
```

### Rollback
```
POST /rollback/{branch_id}
```

---

## Installing & Running

```
go build ./cmd/server
./server
```

---

## Testing

```
go test ./...
go test -bench . -run ^$
```

---

## Roadmap

1. Pluggable merge strategies  
2. Snapshotting & persistence  
3. Distributed “speculative replication”  
4. Multi-branch merges  
5. Conflict visualizers  

---

## Status
This project is experimental and aims to explore **branch-based caching models** as a primitive for decision systems and local-first architectures.
