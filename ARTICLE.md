# Speculative Branching Cache: Managing Temporary State Without Transactions or Raft

## 1. Introduction

This article describes a lightweight approach to handling temporary application state using **speculative branches** over a confirmed key–value store. The goal is to provide a simple and explicit **commit/rollback** model without using transactions, locking, MVCC, or distributed consensus algorithms such as Raft.

The design is intentionally small, local, and easy to embed into services. It can be useful for simulations, previews, rule engines, and local-first application logic where temporary state is needed, but full transactional systems would be too heavy.

---

## 2. Motivation

Modern systems usually fall into two categories:

### • Caches (Redis, Memcached, in-memory KV)
Fast and simple, but with only one visible state. There is no structured way to:
- simulate alternative scenarios,
- test business rules in isolation,
- show previews before applying changes.

Developers typically resort to:
- cloning objects,
- keeping temporary maps,
- embedding "draft" fields in business models.

### • Transactional systems (SQL DBs, MVCC engines)
Powerful and reliable, but:
- relatively heavy,
- not ideal for embedded use,
- unnecessary if all you need is “try local changes and decide later”.

### The missing middle
Forcache fills the conceptual gap:
> **Local speculative state + explicit commit/rollback, without transactional overhead.**

---

## 3. Core Idea

Forcache introduces two layers:

### ✔ Confirmed State
A stable key–value store with a global monotonic version.

### ✔ Speculative Branches
Short-lived overlays attached to a specific `BaseVersion`.

A branch contains:
```
BranchID
BaseVersion
Overlay (key → value)
State (Pending / Committed / Rejected)
```

Reads follow a simple rule:
1. If branch overlay contains the key → return overlay value.
2. Otherwise → return confirmed value.

---

## 4. Why Optimistic Versioning?

Instead of using locks or full MVCC, Forcache uses a global version number.

Commit rule:
```
if Confirmed.Version == Branch.BaseVersion:
    apply overlay and commit
else:
    reject due to conflict
```

### Advantages:
- no locking,
- no deadlocks,
- simple reasoning,
- predictable failures,
- minimal overhead.

This is essentially a small, explicit optimistic concurrency mechanism that fits well into services with moderate contention.

---

## 5. Architecture

```
                 ┌──────────────────────────┐
                 │      Cache Service       │
                 │ (public API: get/put,    │
                 │   spec, commit, rollback)│
                 └─────────────┬────────────┘
                               │
        ┌──────────────────────┴─────────────────────────┐
        │                                                │
┌───────▼────────┐                            ┌──────────▼─────────┐
│  Speculative    │                            │   Confirmed Store  │
│     Store       │                            │ (KV + global ver.) │
│ - branches      │                            │ - stable state     │
│ - overlays      │                            └──────────┬─────────┘
└───────┬────────┘                                       │
        │                                                │
        └──────────────────────┬──────────────────────────┘
                               │
                   ┌───────────▼────────────┐
                   │      Merge Engine       │
                   │   version check,        │
                   │   commit/reject         │
                   └─────────────────────────┘
```

---

## 6. Branch Lifecycle

```
             ┌────────────────┐
             │    Created      │
             │ BaseVersion = X │
             └───────┬────────┘
                     │
                     ▼
             ┌────────────────┐
             │    Active      │
             │  (overlay)     │
             └───────┬────────┘
             ┌────────┴────────┐
             ▼                 ▼
    ┌────────────────┐   ┌──────────────────┐
    │  Commit OK      │   │  Commit Failed   │
    │ versions match  │   │ version changed  │
    └────────┬────────┘   └──────────┬──────┘
             │                        │
             ▼                        ▼
     ┌──────────────┐         ┌────────────────┐
     │  Committed    │         │   Rejected     │
     └──────────────┘         └────────────────┘
```

---

## 7. Read Path

```
                Read(key, branch?)

         ┌──────────────────────────┐
         │ branch specified?        │
         └─────────────┬────────────┘
                       │ yes
                       ▼
          ┌─────────────────────────┐
          │ branch.overlay has key? │
          └───────┬─────────────────┘
                  │ yes
                  ▼
        return overlay[key]

                  │ no
                  ▼
        return confirmed[key]

                       │ no (no branch)
                       ▼
             return confirmed[key]
```

---

## 8. Use Cases

### Menu
- [8.1 Rule Engine Simulations](#81-rule-engine-simulations)
- [8.2 UI Preview / Local-First Behaviour](#82-ui-preview--local-first-behaviour)
- [8.3 Pricing, Scoring, Allocation](#83-pricing-scoring-allocation)
- [8.4 Workflow Orchestration](#84-workflow-orchestration)
- [8.5 Edge / Offline Devices](#85-edge--offline-devices)

---

## 8.1 Rule Engine Simulations

Rule engines often modify a set of attributes, flags, or derived metrics while evaluating business logic.  
Normally, developers copy objects, create temporary maps, or use complex “draft” structures.

With speculative branches:

1. A branch is created from the confirmed state.
2. All rule updates are applied inside the branch overlay.
3. The engine evaluates rules using the branched view.
4. The system either:
    - **commits** the branch (if rule evaluation produces an acceptable state), or
    - **rejects** the branch (if a rule fails or a constraint is violated).

This keeps the rule engine pure: it reads/writes into the same interface, but with an isolated state.

---

## 8.2 UI Preview / Local-First Behaviour

A very common scenario is *preview before applying changes*:

- UI users modify settings, pricing, metadata, or configuration.
- Before saving, they want to see the effect of their changes.

With Forcache:

1. UI sends proposed changes → backend creates a branch.
2. UI fetches data referencing the branch.
3. Backend resolves all reads using the overlay + confirmed state.
4. UI shows the preview that exactly matches post-commit behavior.
5. User presses:
    - **Save** → branch commit
    - **Cancel** → rollback

This eliminates the need for:
- draft tables,
- temporary DB records,
- domain objects with “dirty/draft” flags.

The temporary state lives in the cache, not in the domain model.

---

## 8.3 Pricing, Scoring, Allocation

Many services compute results based on tunable parameters:

- discount rules,
- scoring factors,
- risk weights,
- allocation strategies,
- resource balancing,
- routing or matching logic.

These are often expensive to run and risky to apply directly.

Speculative branches allow you to:

1. Apply experimental parameters in a branch.
2. Run the entire pricing/scoring engine on the branched state.
3. Inspect the resulting values.
4. Choose whether to apply the change globally (commit) or discard it.

This is extremely helpful when analysts or automated systems test multiple scenarios.

---

## 8.4 Workflow Orchestration

Workflows consist of sequential steps:

Step 1 → Step 2 → Step 3 → Step 4 → …

Sometimes each step modifies temporary state, but the final result should be applied only if *all* steps succeed.

Buying a ticket, provisioning a resource, managing a financial transaction — all involve multi-step logic.

With Forcache:

1. Workflow starts → create branch.
2. Step 1 writes partial state into branch.
3. Step 2 reads branch view, updates overlay.
4. …
5. Final step:
    - success → commit
    - failure → rollback

The branch becomes an in-memory transactional boundary **without real transactions or locks**.

---

## 8.5 Edge / Offline Devices

Consider:

- IoT nodes,
- offline-first mobile applications,
- retail POS systems,
- drones,
- embedded controllers,
- systems with intermittent connectivity.

These devices often need to:

- temporarily update state,
- simulate changes,
- make local decisions,
- and only later synchronize with a central system.

Speculative branches provide:

- a lightweight, embedded “temporary state” layer,
- cheap local execution,
- rollback on errors or invalid states,
- no need for DB engines.

The device can:
- work locally while offline,
- accumulate confirmed commits,
- sync them when online again.

This aligns with “local-first software” principles.


---

## 9. Limitations

- not a database,
- not durable storage,
- not a distributed store,
- no multi-version snapshots,
- high contention can cause many conflicts.

This is expected: the model is intentionally minimal.
This design is intentionally local and single-process. Forcache is not a distributed system, but it can be extended with a distributed confirmed state if needed.
---

## 10. Conclusion

Forcache is a small experiment exploring how far we can go with:
- one confirmed state,
- short-lived speculative branches,
- and a simple optimistic commit rule.

It provides a structured way to manage temporary state without the weight of transactions, locks, or consensus systems. For many scenarios, this “middle layer” can be surprisingly practical.

Repository:  
https://github.com/swampus/forcache
