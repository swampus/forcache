
# Forcache Design Document

## 1. Overview

Forcache implements a two-layer state model:

- **Confirmed State** — durable, globally visible KV store with monotonically increasing version.
- **Speculative State** — short-lived branches storing overlays.

Branches behave like:
- local transactions without isolation,
- lightweight MVCC snapshots with explicit commit checks.

This approach provides deterministic semantics with extremely low overhead.

---

## 2. Data Model

### ConfirmedStore
```
map[string]Value
Version (uint64)
```

Writes:
- update value,
- increment version by 1 per key.

### Branch
```
BranchID
BaseVersion
State (Pending/Committed/Rejected)
Overlay: map[key]Value
```

### SpeculativeStore
A simple map: `branchID → Branch`.

---

## 3. Consistency Model

### Rules

1. On branch creation:
   ```
   BaseVersion = Confirmed.Version
   ```

2. Reads with branch:
   ```
   overlay[key] if exists
   else confirmed[key]
   ```

3. Commit succeeds iff:
   ```
   Confirmed.Version == Branch.BaseVersion
   ```

4. Merge:
   - apply overlay in order,
   - bump version per applied key,
   - remove branch.

5. Reject:
   - mark branch as REJECTED,
   - remove but keep metadata optional.

---

## 4. Algorithms

### Branch Creation
```
base = confirmed.version
branch = new Branch(base)
speculativeStore.add(branch)
return branchID
```

### Commit
```
branch = get(branchID)
if confirmed.version != branch.base:
    reject
for each (k,v) in overlay:
    confirmed.set(k,v)
return success
```

### Read
```
if opts.branchID != "":
    if overlay has key:
        return overlay[key]
return confirmed[key]
```

---

## 5. Failure Scenarios

### Commit Conflict
If a parallel write occurs before commit, branch is rejected.

### Stale Branches
Branches may become invalid; this is expected.

### No Partial Merge
A commit applies all operations atomically relative to the version check.

---

## 6. Future Extensions

- Per-key merge policies  
- Multi-layered branching (branch-on-branch)  
- Distributed “speculative synchronization”  
- Persistent backend  

---

## 7. Rationale

This model is powerful enough for simulations, local-first workflows, interactive UIs, and pre-validation steps, while remaining much simpler than full transactional subsystems.
