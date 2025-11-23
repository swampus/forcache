
# Forcache Testing Strategy

## 1. Goals
- Validate correctness of branching semantics.
- Ensure merge engine works under optimistic concurrency.
- Guarantee read isolation logic.
- Provide minimal performance benchmarks.
- Maintain clarity and determinism.

---

## 2. Types of Tests

### 2.1 Unit Tests
- `confirmed_test.go`: versioning, set/get semantics, commit logic.
- `speculative_test.go`: branch creation, isolation, deletion.
- `merge_test.go`: commit success, conflict detection, rollback.

All tests are table-driven and use subtests.

### 2.2 Integration Tests
- `service_test.go`: full spec → get → commit → get flow.
- Covers read isolation and merging effects.

### 2.3 API Tests
- `handlers_test.go`: full HTTP flow using `httptest.Server`.

### 2.4 Benchmarks
- `bench_test.go`: benchmark for confirmed store writes.

Example command:
```
go test -bench . -run ^$
```

---

## 3. Test Architecture

```
 Unit Tests → Service Tests → HTTP Tests → Benchmarks
   (internal)      (facade)        (API)         (perf)
```

---

## 4. Race Validation

Recommended:
```
go test -race ./...
```

---

## 5. Coverage

```
go test -cover ./...
```

Recommended threshold: ~70%+ for core logic.

---

## 6. Conclusion

The test suite ensures deterministic correctness, and covers core invariants crucial for caching systems:
- version monotonicity,
- branch isolation,
- conflict detection,
- merge correctness,
- API contract stability.
