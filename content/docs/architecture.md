# Runink Architecture: Go/Linux Native Distributed Data Environment

Goal: Define a self-sufficient, distributed environment for orchestrating and executing data pipelines using Go and Linux primitives. This system acts as the cluster resource manager and scheduler (replacing Slurm), provides Kubernetes-like logical isolation and RBAC, integrates data governance features, and ensures robust security and observability. It aims for high efficiency by avoiding traditional virtualization or container runtimes like Docker. Define a self-sufficient, distributed environment for orchestrating and executing data pipelines using Go and Linux primitives, with enhanced metadata capabilities designed to support standard data governance (lineage, catalog) AND future integration of LLM-generated annotations.

## High-Level Architecture

Runink operates with a Control Plane managing multiple Worker Nodes, each running a Runi Agent.

```plaintext
+---------------------------------------------------------------+
|                      Developer / Authoring Layer              |
|---------------------------------------------------------------|
|  - CLI (runi)                                                 |
|  - REPL (DataFrame-style, step-by-step)                       |
|  - DSL Authoring (`features/`)                                |
|  - Contract Definitions (`contracts/`)                        |
|  - Test Engine (Golden tests, synthetic fixtures)             |
+---------------------------------------------------------------+
                              ↓
+---------------------------------------------------------------+
|                   Control Plane (Domain + Orchestration)      |
|---------------------------------------------------------------|
|  - API Server (REST/gRPC, OIDC auth, herd scoping)            |
|  - Scheduler (declarative constraint solver, DAG generator)   |
|  - Metadata Catalog (contracts, lineage, tags, scenarios)     |
|  - Secrets Manager (Herd-scoped secrets over Raft)            |
|  - Compliance & RBAC Enforcer (contract-level access control) |
|  - Checkpoint Coordinator (tracks partial run status + resume)|
+---------------------------------------------------------------+
                              ↓
+---------------------------------------------------------------+
|                 Distributed State Layer (Consensus + Storage) |
|---------------------------------------------------------------|
|  - Raft Consensus Group (Herd + Scheduler + State sync)       |
|  - BadgerDB-backed volumes for:                               |
|      • Pipeline run metadata                                  |
|      • Slice-local state volumes                              |
|      • Checkpointed outputs (resume-able DAG stages)          |
|  - Herd Definitions, Contract Versions, RBAC, Secrets         |
+---------------------------------------------------------------+
                              ↓
+---------------------------------------------------------------+
|                   Execution Plane (Slices + Agents)           |
|---------------------------------------------------------------|
|  - Runi Agent (on every node, manages slices + volumes)       |
|  - Slice Group Supervisor (windowed runs, state volumes)      |
|  - Slice Process:                                             |
|      • Mounts ephemeral namespace, PID net, user, cgroup      |
|      • Loads step function (from contract)                    |
|      • Streams data with io.Pipe (zero-copy)                  |
|      • Accesses read/write state volume via volume proxy      |
|      • Writes to sink, emits lineage                          |
|  - Windowed Join Runner:                                      |
|      • Centralized lineage rehydration across slice groups    |
|      • Co-group on keys across parallel slice outputs         |
+---------------------------------------------------------------+
                              ↓
+---------------------------------------------------------------+
|              Governance, Observability, and Quality Plane     |
|---------------------------------------------------------------|
|  - Data Quality Rules (contract tags + field checks)          |
|  - Lineage Engine (runID, contract hash, input/output edges)  |
|  - Metrics Exporter (Prometheus, Fluentd, OpenTelemetry)      |
|  - DLQ + Routing Controller (based on step & validation tags) |
|  - Audit Log Engine (per slice, signed logs)                  |
+---------------------------------------------------------------+
                              ↓
+---------------------------------------------------------------+
|             External Sources / Sinks / State APIs             |
|---------------------------------------------------------------|
|  - Kafka, S3, Snowflake, GCS, Postgres, Redshift              |
|  - APIs (FDC3, CDM, internal REST sources)                    |
|  - Optional backing state for key joins (Redis, RocksDB)      |
+---------------------------------------------------------------+
```

### Detailed architecture

```plaintext
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
|                             External User / Client                                                        |
|               (Runink CLI, UI, API Client via OIDC/API Key Auth)                                          |
|                (Requests scoped to specific 'Herds'/Namespaces)                                           |
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
                                     | (gRPC/REST API Calls - TLS Encrypted)
                                     | (Pipeline Definitions, Queries within Herd context)
                                     v
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
|                HERD CONTROL PLANE (HA Replicas using Raft)                                                |
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
|                                                                                                           |
|  [1. API Server (Replica 1..N)] <───────╮ (Internal gRPC - TLS)                                           |
|  |- Entry Point                         │                                                                 |
|  |- AuthN/AuthZ (OIDC, RBAC) <──────────+──────[4. Identity & RBAC Mgr]                                   |
|  |  *Enforces RBAC per Herd*            │      (Runs Co-located or Standalone)                            |
|  |- Validation (inc. Herd scope)        │                                                                 |
|  |- Gateway                             │                                                                 |
|  |                                      │                                                                 |
|  |──────────────────────────────────────╯                                                                 |
|  │ ▲ │                                                                                                    |
|  │ │ └────────────────────────────────────────────────────────────────────────────────┐                   |
|  │ │ (State Read/Write Requests to Raft Leader)                                       |                   |
|  │ │ (*Read* Request potentially served by Leader/Follower)                           |                   |
|  │ ▼                                                                                  |                   |
|  [2. Cluster State Store / "Barn" (Nodes 1..N)] <──────────> [Raft Consensus Protocol]|                   |
|  - *Achieves State Machine Replication via Raft                                       |                   |
|  |       Consensus for HA & Consistency*                                              |                   |
|  |       (Leader Election, Log Replication)                                           |                   |
|  - Stores: Nodes, Pipelines, Slices,                                                  |                   |               
|  |         Logical Herds (Definition, Quotas),                                        |                   |
|  |         RBAC Policies, Secrets, Core Metadata Refs                                 |                   |
|  │                                                                                    |                   |
|  │  [ Raft Leader ] <------(Write Op)-- Clients (API Srv, Scheduler, etc.)            |                   |
|  │   |                                                                                |                   |
|  │   | 1. Append Op to *Own* Log                                                      |                   |
|  │   v                                                                                |                   |
|  │  [ Replicated Log ] ----(Replicate Entry)--> [ Raft Followers (N-1) ]              |                   |
|  │   |                     <--(Acknowledge)------                                     |                   |
|  │   | 2. Wait for Quorum                                                             |                   |
|  │   v                                                                                |                   |
|  │  [ Commit Point ] -----(Notify Apply)-----> [ Leader & Followers ]                 |                   |
|  │   |                                                                                |                   |
|  │   v 3. Apply Committed Op to *Local* State Machine                                 |                   |
|  │  [ State Machine (KV Store) @ Leader ]                                             |                   |
|  │                                                                                    |                   |
|  └───────── // Same process applies committed Ops on Followers // ────────────────────┘                   |
|  │ ▲          [ State Machine (KV Store) @ Follower 1..N-1 ]                                              |
|  │ │                                                                                                      |
|  │ │ (Consistent Reads, Watches on Committed State)                                                       |
|  │ │                                                                                                      |
|  │ ▼                                                                                                      |
|  [3. Scheduler (Active Leader or Replicas)] ────────→┤ (Placement Cmds via API Server)                    |
|  | - Reads Consistent State from Raft Store          │                                                    |
|  | - Considers Task Reqs, Herd Quotas, Affinity      │                                                    |
|  |                                                   │                                                    |
|   ───────────────────────────────────────────────────┘                                                    |
|   │ ▲ │                                                                                                   |
|   │ │ └───────────────────────────────────────────────────┐                                               |
|   │ │ (Secrets Read/Write Requests to Raft Leader)        │                                               |
|   │ │                                                     │                                               |
|   │ ▼                                                     │                                               |
|  [5. Secrets Manager]────────────────────────────────────>┤ (Secrets Delivery via Runi Agent)             |
|  |- Reads/Writes Encrypted Secrets via Raft Store         │                                               |
|  |- RBAC for Secrets Access (Enforced via API Server)     │                                               |
|  |- Interact with State Store via API Server/Directly     |                                               |
|  |- Read consistent state, submit changes via Leader      |                                               |                                                      
|  +────────────────────────────────────────────────────────┘                                               |
|   │ ▲                                                                                                     |
|   │ │ (Metadata Read/Write Requests to Raft Leader/Service)                                               |
|   │ │                                                                                                     |
|   │ ▼                                                                                                     |
|  [6. Data Governance Service]───────────────────────+─> (Queries via API Server)                          |
|  |- Stores Core Metadata/Refs via Raft Store        │                                                     |
|  |- May use separate backend for Annotations        │                                                     |
|  |- Provides Metadata gRPC API                      │                                                     |
|  |                                                  │                                                     |
+────────────────────────────────────^────────────────^─────────────────────────────────────────────────────+
|                                    │                │ (Internal gRPC/TLS)                                 |
| (Agent Registration, Health,       │                │ (Cmds to Agent from Control Plane Leader/API)       |
|  Metrics Reporting, Log Forwarding)│                │ (Cmds: Launch Slice *in Herd*, Secrets Req, Status) | 
|                                    │                │ (Metadata Reporting via Worker to Gov Service)      |
+────────────────────────────────────v────────────────v─────────────────────────────────────────────────────+
|                      RUNI AGENT (Daemon on each Worker Node)                                              |
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
|  - Registers with Control Plane (via API Server to Raft Store)                                            |
|  - Reports Node Resources & Health (via API Server to Raft Store)                                         |
|  - Receives Commands (from Scheduler via API Server)                                                      |
|  - Manages Local `Slice Manager` (cgroups, namespaces, exec)                                              |
|  - Manages Local Worker Lifecycle                                                                         |
|  - Securely Fetches Secrets (from Secrets Mgr via API Server)                                             |
|  - Sets up Network & Storage for Workers                                                                  |
|  - Collects/Forwards Worker Observability (Logs to Fluentd, Metrics Scrape)                               |
|  - Exposes Own & Aggregated Metrics                                                                       |
|                                                                                                           |
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
|           (Launch Cmd *within Herd context*,    │ ▲    (Log/Metrics Collection)                           |
|            Config, Secrets)                     │ |                                                       |
|                                                 ▼ │                                                       |
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
|                  RUNI SLICE (Isolated Go Process)                                                         |
|                    *Runs within the context of a specific Herd*                                           |
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
|  - Executes User Pipeline Step Code                                                                       |
|  - Isolated via Namespaces (User, PID, Net, Mount, etc.)                                                  |
|  - Constrained by Cgroups (*reflecting Herd quotas*)                                                      |
|  - Runs as Ephemeral User (mapped from Service Account, authorized for Herd)                              |
|  - Receives Config/Secrets                                                                                |
|                                                                                                           |
|  **Interactions:**                                                                                        |
|  - Communicates with Other Slices (gRPC, TLS, AuthN)                                                      |
|  - Accesses External Services (DBs, MinIO, OpenAI)                                                        |
|  - Reports Lineage/Annotations to `Data Governance Service` (*tagged by Herd*)                            |
|  - Exposes */metrics* endpoint (scraped by Runi Agent)                                                    |
|  - Writes Structured Logs to stdout/stderr (collected by Runi Agent)                                      |
|                                                                                                           |
+───────────────────────────────────────────────────────────────────────────────────────────────────────────+
```

## Core Principles

* **User Interaction:** Client requests are often scoped to a specific Herd.
* **API Server / RBAC:** Enforces RBAC policies based on user/service account permissions within a target Herd.
* **Cluster State Store:** Explicitly stores Herd definitions and their associated resource quotas.
* **Scheduler:** Considers Herd-level quotas when making placement decisions.
* **Secrets Manager:** Access to secrets might be scoped by Herd.
* **Data Governance:** Metadata (lineage, annotations) can be tagged by or associated with the Herd it belongs to.
* **Runi Agent:** Receives the target Herd context when launching a slice and uses this information to potentially configure namespaces and apply appropriate cgroup limits based on Herd quotas.
* **Runi Slice:** A single instance of a pipeline step running as an isolated Worker Slice Process. Executes entirely within the logical boundary and resource constraints defined by its assigned Herd.
* **Runi Pipes:** Primarily used now for internal communication within the `Runi Agent` to capture logs/stdio from the `Runi Slice` it `exec`s, rather than for primary data transfer between steps.
* **Herd:** A logical grouping construct, similar to a Kubernetes Namespace, enforced via RBAC policies and potentially mapped to specific sets of Linux namespaces managed by Agents. Provides multi-tenancy and team isolation. Quotas can be applied per Herd.
* **Go Native & Linux Primitives:** Core components written in Go, directly leveraging cgroups, namespaces (user, pid, net, mount, uts, ipc), pipes, sockets, and `exec` for execution and isolation.
* **Self-Contained Cluster Management:** Manages a pool of physical or virtual machines, schedules workloads onto them, and handles node lifecycle.
* **Serverless Execution Model:** Users define pipelines and resource requests; Runink manages node allocation, scheduling, isolation, scaling (by launching more slices), and lifecycle. Users are subject to quotas managed via cgroups.
* **Security First:** Integrated identity (OIDC), RBAC, secrets management, network policies, encryption in transit/rest.
* **Data Governance Aware:** Built-in metadata tracking, lineage capture, and support for quality checks. With extension for storage/management of rich data annotations (e.g., from LLMs).
* **Rich Observability:** Native support for metrics (Prometheus) and logs (Fluentd).

---

# 🧠 Programming Approaches: Why They Power Runink

Runink’s architecture isn’t just Go-native — it’s intentionally designed around a few **low-level but high-impact programming paradigms**. These concepts are what let Runink outperform containerized stacks, enforce security without overhead, and keep pipelines testable, composable, and fast.

Below, we walk through the four core techniques and where they show up in Runink’s components.

---

### 1. **Data-Oriented Design (DOD)**
_"Design for the CPU, not the developer."_

Instead of designing around objects and interfaces, optimize around **memory layout and data access patterns**. This is especially important for Runink’s slice execution and contract validation stages, where predictable access to batches of structs (records) matters.

- Use slices of structs over slices of pointers to enable **CPU cache locality**.
- Align field access with columnar memory usage if streaming transforms run across many rows.
- Preallocate buffers in `Runi Agent`’s slice execution path to avoid GC churn.

---

### Core Idea:  
Layout memory for how it will be accessed, not how it's logically grouped. Runink’s slices often scan, validate, or enrich large batches of records — so struct layout, batching, and memory predictability **directly impact performance**.

### Apply DOD in Runink:
- Prefer **flat structs over nested ones** in contracts:
  ```go
  // Better
  type User struct {
    ID    string
    Name  string
    Email string
  }

  // Avoid: nested fields unless necessary
  type User struct {
    Meta struct {
      ID string
    }
    Profile struct {
      Name string
      Email string
    }
  }
  ```


### In `runi slice`:
- Use `sync.Pool` for reusable buffers (especially JSON decoding).
- Pre-size buffers based on contract hints (e.g., `maxRecords=10000`).
- Avoid interface{} — use generated structs via `go/types` or `go:generate`.

### Benefits: 
- Better **memory throughput**, fewer allocations, and **Go GC-friendliness** under load.
- Use **structs of arrays (SoA)** or `[]User` with preallocated slices in transformations.
- Minimize pointer indirection. Use value receivers and avoid `*string`, `*int` unless you need nil.
- Design transforms that operate in **tight loops**, e.g., `for _, rec := range batch`.

---

### 2. **Zero-Copy and Streaming Pipelines**
_"Avoid full in-memory materialization — process as the data flows."_

Avoid unnecessary data marshaling or full deserialization. Rely on:

- **`io.Reader` / `io.Writer` interfaces** for pipeline steps.
- `os.Pipe()`, `bufio`, or `net.Pipe()` for intra-slice streaming.
- Only materialize records **when needed for validation or transformation**.

In slices, structure execution as `stream -> transform -> stream`, not `[]record -> map -> []record`.

---

### Core Idea:
Instead of `[]Record -> Transform -> []Record`, operate on **streams of bytes or structs** using `io.Reader`, `chan Record`, or even UNIX pipes between stages.

### Runink Optimizations:
- Use **`io.Reader → Decoder → Transform → Encoder → io.Writer`** chain.
- Design step transforms like this:

  ```go
  func ValidateUser(r io.Reader, w io.Writer) error {
    decoder := json.NewDecoder(r)
    encoder := json.NewEncoder(w)

    for decoder.More() {
      var user contracts.User
      if err := decoder.Decode(&user); err != nil {
        return err
      }
      if isValid(user) {
        encoder.Encode(user)
      }
    }
    return nil
  }
  ```

- For multi-stage slices, use **`os.Pipe()`**:
  ```go
  r1, w1 := os.Pipe()
  r2, w2 := os.Pipe()

  go Normalize(r0, w1) // input -> step 1
  go Enrich(r1, w2)    // step 1 -> step 2
  go Sink(r2, out)     // step 2 -> sink
  ```

### Benefits:
- **Constant memory** even for massive datasets.
- **Backpressure**: If downstream slows down, upstream blocks — great for streaming (Kafka, etc.).
- Enables **DLQ teeing**: `tee := io.MultiWriter(validOut, dlqSink)`.

---

### 3. **Declarative Scheduling with Constraint Propagation**
_"Schedule via logic, not instructions."_
The `Herd` and `Runi Agent` coordination already benefits from Raft-backed state, but push it further with **affinity-aware, declarative scheduling**:

- Define **resource constraints** (e.g., `@requires(cpu=2, memory=512Mi, label=“gpu”)`) in `.dsl`.
- Propagate **slice placement decisions** through constraint-solving logic instead of imperative scheduling.
- Record constraints in the Raft-backed state store to enforce deterministic task placement.

You can build this as a small DSL-on-DSL layer (e.g. `@affinity(domain="finance", colocate_with="step:JoinUsers")`).

*Benefit:* Stronger **determinism**, **replayability**, and **multi-tenant safety**.

### Core Idea:
Model placement as a **set of constraints**: affinity, herd quota, GPU needs, tenant isolation, etc. Let the scheduler **solve** the placement rather than being told where to run.

### Runink DSL Extension:
In `.dsl`:

```gherkin
@step("RunLLMValidation")
@affinity(label="gpu", colocate_with="step:ParsePDFs")
@requires(cpu="4", memory="2Gi")
```

This can be compiled into metadata stored in the Raft-backed scheduler store.

### Scheduler Logic (Pseudo-Go):
```go
type Constraints struct {
  CPU       int
  Memory    int
  Affinity  string
  Colocate  string
  HerdScope string
}

func ScheduleStep(stepID string, constraints Constraints) (NodeID, error) {
  candidates := filterByHerd(constraints.HerdScope)
  candidates = filterByResources(candidates, constraints.CPU, constraints.Memory)
  if constraints.Colocate != "" {
    candidates = colocateWith(candidates, constraints.Colocate)
  }
  if constraints.Affinity != "" {
    candidates = matchLabel(candidates, constraints.Affinity)
  }
  return pickBest(candidates)
}
```
### Bonus: All decisions can be **Raft-logged**, making slice scheduling **auditable and replayable**.

---

## ✨ Designed for Infra-Native Developers

Whether you're used to **Spark DAGs**, **Airflow tasks**, or **Python dataframes**, Runink gives you:

- **Lower startup latency** than containers
- **Testable pipeline logic** with clear steps
- **No need for external orchestrators**
- **Secure-by-default execution** with zero-trust isolation

By combining Go, Raft, Linux namespaces, and functional data transforms — Runink builds pipelines that feel like **code, not infrastructure**.

---

# 🧩 Go Slices, Linux Primitives, and Raft: Runink’s Execution Philosophy
**Host-Efficient Execution for Secure Data Ingestion**

Runink takes a radically different approach to pipeline execution than traditional data platforms — instead of running heavy containers, JVMs, or external orchestrators, Runink uses **Go-native workers**, Linux primitives like **cgroups and namespaces**, and concepts like **data-oriented design and zero-copy streaming** to deliver **blazing-fast, memory-stable, and secure slices**.


## 🔍 Overview

Runink executes data pipelines using **Go "slices"** — lightweight, isolated execution units designed to model both **bounded** (batch) and **unbounded** (streaming) data workloads. These are:

- Spawned by the **Runi agent** on worker nodes
- Executed as **isolated processes**
- Scoped to **Herd namespaces**
- Constrained by **cgroups**
- Communicate via **pipes**, **sockets**, or **gRPC streams**

This orchestration is **Raft-coordinated**, making every launch deterministic, fault-tolerant, and observable.

---

## 🧬 What Are Bounded and Unbounded Slices?

| Type         | Use Case                            | Description |
|--------------|-------------------------------------|-------------|
| **Bounded**  | Batch ETL, contract validation       | Processes a finite dataset and terminates |
| **Unbounded**| Streaming ingestion, log/event flows | Long-running, backpressured pipelines with checkpointing |

Both types run as **Go processes** within a controlled **Herd namespace**, and can be composed together in DAGs.

---

## 🧰 Slice Internals: Go + Linux Synergy

Each slice is a **native Go process** managed via:

### ✅ Cgroups
- Applied per Herd, per slice
- Limits on CPU, memory, I/O
- Enforced using Linux `cgroupv2` hierarchy
- Supports slice preemption and fair resource sharing

### ✅ Namespaces
- User, mount, network, and PID namespaces
- Enforce isolation between tenants (Herds)
- Prevent noisy-neighbor problems and info leaks

### ✅ Pipes & IPC
- Use of `os.Pipe()` or `io.Pipe()` in Go to model stage-to-stage communication
- `net.Pipe()` and UNIX domain sockets for local transport
- Optionally enhanced via `io.Reader`, `bufio`, or gRPC streams (for cross-node slices)

### ✅ Execution
- `os/exec` with setns(2) and clone(2) to launch each slice
- Environment-injected config and secrets fetched securely via Raft-backed Secrets Manager

---

## 🔄 Raft as Execution Backbone

The **Barn (Cluster State Store)** is Raft-backed. It ensures:

| Raft Role     | Benefit to Runink                                      |
|---------------|---------------------------------------------------------|
| **Leader Election** | Prevents race conditions in pipeline launches |
| **Log Replication** | Guarantees all agents/schedulers share same DAG, lineage, and config |
| **Strong Consistency** | Execution decisions are deterministic and audit-traceable |
| **Fault Tolerance** | Node crashes do not corrupt state or duplicate work |

**Examples of Raft-integrated flows:**
- DAG submission is a Raft log entry
- Herd quota changes update slice scheduling state
- DLQ routing is replicated for contract validation violations
- Slice termination is consensus-driven (no orphaned processes)

---

### 1. **Data-Oriented Design (DOD)**

Rather than designing around objects or classes, Runink favors **struct-of-arrays** and flat memory access patterns. Why?

- Go structs are faster to iterate over than Python dictionaries or Java POJOs.
- Access patterns align with how CPUs fetch and cache data.
- Contract validation and transforms run over **preallocated `[]struct` batches**, not heap-bound objects.

💡 *For Python/Java devs:* Think of this like switching from `dicts of dicts` to **NumPy-like** flat arrays — but in Go, with static types and no GC spikes.

---

### 2. **Zero-Copy Streaming Pipelines**

Every slice in Runink processes data as a stream:

- Uses `io.Reader` / `io.Writer` rather than buffering everything in memory.
- Transforms run as **pipes between goroutines** — like UNIX but typed.
- Memory stays flat, predictable, and bounded — even for 10M+ record streams.

💡 *For pandas/Spark devs:* This is closer to **generator pipelines** or **structured stream micro-batches**, but with Go’s backpressure-aware channels and streaming codecs.

---

### 3. **Slice-Level Isolation via Linux Primitives**

Each pipeline stage runs in a **Runi Slice**, a native Go process launched via:

- **`os/exec`**, inside a unique PID/user/mount/network **namespace**
- With **cgroupv2** constraints for CPU, memory, and I/O
- Injected with secrets via **Raft-backed vault**, scoped to Herd + slice

This gives container-like security and isolation **without Docker**.

💡 *If you're used to Kubernetes or Docker:* Think of slices as **ephemeral containers, but 10x faster** — no image pulls, no pod scheduling latency.

---

### 4. **Declarative Scheduling with Constraint Solving**

Runink doesn't “assign pods” — it **solves constraints**:

- Define resource needs in `.dsl` (`@requires(cpu=4, memory=2Gi)`)
- Add placement rules (`@affinity(colocate_with=step:Join)`)
- Scheduler uses a **Raft-backed store** to enforce and replay decisions

💡 *In Airflow, scheduling is time-based. In Spark, it’s coarse-grained. In Runink, it’s **placement-aware, deterministic, and multi-tenant safe.***

---

### 5. **Golden Testing and Contract-Aware Execution**

All pipelines are backed by **Go structs as schema contracts**:

- Versioned, diffable, validated at runtime
- Pipeline steps are tested using **golden output** and **synthetic data generators**
- Runink can **detect schema drift** and route invalid records to DLQ

💡 *This gives you the safety of DBT with the flexibility of code-first transforms.*

---

## ✨ For Java and Python Developers

**You won’t see:**

- JVM warmup times  
- YAML config sprawl  
- Container image builds  
- Cross-language serialization bugs  

**You’ll get:**

- Sub-millisecond slice startups  
- Schema-first data governance baked in  
- Functional transforms with typed inputs/outputs  
- Pipelines you can test, diff, and re-run like unit tests  

**Runink slices run like fast, secure micro-VMs — written in Go, isolated with Linux, coordinated by Raft.**  
This gives you the **speed of compiled code**, the **safety of contracts**, and the **audibility of Git.**

>***No containers, no clusters — just data pipelines that behave like code.***

---

## 🔍 Why These Approaches Matter (Component-Wise Breakdown)

---

### 🎯 Functional Pipelines  
> Used in: `Pipeline Generator`, `REPL`, `Testing Engine`

- The `.dsl` files describe *what* should happen, not *how*. This declarative structure allows Runink to generate **Go DAGs** with **pure functions** per step (`func Transform(in io.Reader) io.Reader`).
- Each step is **stateless**, making it safe to test in isolation or stub during REPL prototyping.
- Enables **golden file testing** with fixed input/output expectations — like unit tests for data.

🧠 *Why this works:* You get **Airflow DAG logic** with **Go-like testability** and **Unix pipeline semantics** — without Python globals or scheduler state leaks.

---

### 🧱 Data-Oriented Design (DOD)  
> Used in: `Schema Contracts`, `Runi Slice`

- Contracts compile to **flat Go structs**, not JSON blobs or reflective maps.
- In slices, contracts are used in **tight, cache-friendly loops** for validation, enrichment, and routing.
- Avoids pointer chasing (`*string`) and nested structs — aligns memory layout to CPU caching.

💡 *Why this works:* Validating 1M records becomes a **loop over structs**, not a heap traversal. GC pressure drops. Throughput climbs.

---

### 🔄 Zero-Copy Streaming  
> Used in: `Runi Agent`, `Slice Manager`, `REPL`

- Transforms use `io.Reader → Decoder → Transform → Encoder → io.Writer` instead of `[]record → []record`.
- `os.Pipe()` and `net.Pipe()` connect slice stages like Unix pipes, allowing **constant memory flow**.
- If downstream slows, upstream **naturally blocks** — no need for buffering, polling, or timeouts.

⚡ *Why this works:* You process 100GB streams with **<100MB RAM**, and gain DLQ, audit, and observability inline — perfect for streaming Kafka or big batch ETL.

---

### 🧠 Declarative Scheduling  
> Used in: `Scheduler`, `Barn`, `API Server`

- Constraints like `@requires(cpu=2, label="gpu")` are parsed and **stored in Raft** as part of scheduling state.
- The scheduler solves placement **logically**, not imperatively: it picks a node that *satisfies* the constraint — not one hardcoded by a script.
- Quotas, affinity, and co-location rules are enforced **at the control plane**, across all agents.

🔐 *Why this works:* Scheduling is **deterministic**, **auditable**, and **multi-tenant safe** — slices run where they should, not where they might.

---

| **Approach**               | **Runink Component**                | **Why It Matters**                                                                                                                                   | **Key Gains**                                      |
|---------------------------|-------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------|
| **Functional Pipelines**  | `Pipeline Generator`, `REPL`, `Testing Engine` | `.dsl` scenarios compile into pure, testable Go transforms — no hidden state, no globals. Steps are reusable and composable.                         | Predictable DAGs, step reuse, snapshotable output  |
| **Data-Oriented Design**  | `Schema Contracts`, `Slice Internals` | Contracts compile into flat Go structs, used in hot paths like slice execution. Avoids nested heap objects that hurt performance under load.         | GC control, memory locality, faster contract eval  |
| **Zero-Copy Streaming**   | `Runi Agent`, `Runi Slice`, `REPL`   | Stages communicate via `io.Reader`/`Writer`, `os.Pipe`, and JSON streams — not in-memory records. Transforms process as they receive data.          | Constant memory, low-latency, backpressure-safe    |
| **Declarative Scheduling**| `Scheduler`, `API Server`, `Barn`    | Placement decisions are driven by constraints in `.dsl` (`@requires`, `@affinity`) and stored in the Raft-backed store for deterministic execution. | Multi-tenant fairness, auditability, resilience    |

---

## 🧠 Raft-Backed Architecture

At the heart of Runink's distributed control plane is Raft consensus, which ensures strong consistency and availability across multiple replicas in a multi-tenant setup. Here’s how Raft plays a crucial role:

### Cluster State Store (Barn):

*Raft ensures* consistency between multiple replicas of the state store, preventing conflicts and enabling data to be synchronized across different services. With critical data, such as pipeline definitions, herd configurations, secrets, and RBAC policies, is replicated and stored in the Barn, which is under Raft consensus.

*Writes and Reads*: All writes (e.g., task assignments, secrets, RBAC policies) go through the Raft leader to ensure consistency. If the leader fails, Raft elects a new leader to ensure continuous operation without data loss.

### Scheduler:

Uses Raft's consensus for task scheduling to ensure that nodes in the cluster agree on which worker slices should be assigned the next task.
Task assignments are synchronized across the cluster to avoid conflicts and duplications.

### Secrets Manager:

Managed via Raft to ensure that secrets are consistently available across all nodes, with encrypted storage and strict RBAC.
Guarantees that secrets are securely injected into Runi agents and worker slices at runtime, with access scoped per Herd.

### Data Governance & Lineage:

Lineage and metadata are tagged and tracked consistently across the cluster, with Raft-backed guarantees to prevent discrepancies in historical records.

### Resilience & High Availability:

Raft consensus enables Runink to maintain availability and consistency in the face of node failures by synchronizing state across multiple replicas.
The API Server, Scheduler, Secrets Manager, and other components benefit from Raft to ensure high availability, disaster recovery, and distributed state consistency.

With Raft integrated into Runink, the system can operate fault-tolerantly, ensuring that data across the entire platform remains consistent, even when network partitions or node failures occur. This guarantees that your data pipelines, metadata, secrets, and workload scheduling are managed reliably in any cloud-native or distributed environment

---

## Architecture Components

```plaintext
I. Runink Control Plane (Distributed & Persistent)
   │
   ├─ Description: The central brain of the Runink cluster. Likely needs to be run in a High Availability (HA) configuration using multiple replicas.
   │
   ├─ 1. API Server (Go Service)
   │   │
   │   ├─ Role: Primary entry point for all external (CLI, UI, API clients) and internal (Node Agents) communication.
   │   ├─ Responsibilities:
   │   │   • Exposes gRPC (primary) and potentially REST APIs.
   │   │   • Authentication: Handles user/client authentication (e.g., OIDC integration, API keys). Issues internal tokens.
   │   │   • Authorization: Enforces RBAC policies on all incoming requests based on user/service account identity.
   │   │   • Request Validation: Validates pipeline definitions, resource requests, etc.
   │   │   • Gateway to other control plane services.
   │   │   • Stores pointers/core info for Data Governance.
   │   └─ Security: Uses TLS for all communication.
   │
   ├─ 2. Cluster State Store (e.g., Raft-based Go Store)
   │   │
   │   ├─ Role: Distributed, consistent, and persistent storage for all cluster state. Replaces the need for external databases or etcd for core state.
   │   ├─ Responsibilities:
   │   │   • Stores information about: Nodes (registration, resources, health), pending/running pipelines & slices, logical Herds (namespaces), Service Accounts, RBAC Policies, Secrets, Data Catalog/Lineage metadata.
   │   │   • Uses a consensus algorithm (like Raft - e.g., `hashicorp/raft`) for HA and consistency across control plane replicas.
   │   │   • Data typically persisted to disk on control plane nodes (potentially encrypted at rest).
   │   └─ Implementation: Could use libraries like `hashicorp/raft` combined with embedded key-value stores (BoltDB, BadgerDB).
   │
   ├─ 3. Scheduler (Go Service)
   │   │
   │   ├─ Role: Matches pending pipeline steps (Worker Slices) with available resources on registered `Runi Agents`. Replaces Slurm's scheduling function.
   │   ├─ Responsibilities:
   │   │   • Watches the `Cluster State Store` for new pending tasks.
   │   │   • Maintains an inventory of available resources on each `Runi Agent` node.
   │   │   • Considers pipeline requirements (CPU, memory, specific hardware), Herd quotas, RBAC permissions, and potentially data locality hints.
   │   │   • Makes placement decisions and instructs the target `Runi Agent` (via API Server) to launch the Worker Slice.
   │   │   • GPU aware for local LLM nodes.
   │   └─ Sophistication: Can range from simple bin-packing to more complex policies.
   │
   ├─ 4. Identity & RBAC Manager (Go Service / Library)
   │   │
   │   ├─ Role: Manages identities and access control policies.
   │   ├─ Responsibilities:
   │   │   • Manages internal Service Accounts used by system components and pipelines.
   │   │   • Handles OIDC federation for user identities.
   │   │   • Generates and validates internal authentication tokens (e.g., JWTs) used for secure communication between components.
   │   │   • Provides interface for defining and querying RBAC policies stored in the `Cluster State Store`.
   │   └─ Integration: Tightly integrated with the `API Server` for enforcement.
   │
   ├─ 5. Secrets Manager (Go Service / Library)
   │   │
   │   ├─ Role: Securely stores and manages sensitive information (e.g., database passwords, API keys).
   │   ├─ Responsibilities:
   │   │   • Stores secrets encrypted within the `Cluster State Store`.
   │   │   • Provides mechanisms for `Runi Agents` to securely retrieve secrets needed by `Worker Slices`.
   │   │   • Integrates with RBAC to control access to secrets.
   │   │   • API keys for external LLMs or credentials for internal ones.
   │   └─ Delivery: Secrets delivered to workers via secure environment variables or mounted tmpfs files within the slice's namespace.
   │
   └─ 6. Data Governance Service (Go Service)
       │
       ├─ Role: Central hub for all metadata: catalog, lineage, quality, and extensible annotations. Acts as the system of record for understanding data context.
       ├─ Responsibilities:
       │   • **Data Catalog:** Stores schemas, descriptions, owners, tags, business glossary links.
       │   • **Data Lineage:** Records the graph of data transformations (inputs -> transform step -> outputs) reported by Worker Slices. Links data artifacts across pipeline steps.
       │   • **Annotation Management:**
       │     * Stores extensible key-value annotations associated with specific data artifacts (datasets, columns, records, pipeline runs).
       │     * **LLM Annotation Schema:** Defines structure to capture LLM-specific details:
       │       - `annotation_type` (e.g., 'classification', 'summary', 'entity_extraction')
       │       - `annotation_value` (the actual text/tag generated)
       │       - `source_model` (e.g., 'openai:gpt-4-turbo', 'minio:llama3-8b-v1.0')
       │       - `confidence_score` (float, if applicable)
       │       - `prompt_reference` (optional, hash or ID of prompt used)
       │       - `timestamp`
       │       - `annotated_data_ref` (pointer to the specific data element)
       │   • **Quality Rules & Results:** Stores data quality rule definitions and results reported by quality check steps.
       │   • **API:** Provides gRPC API for:
       │     * Registering/Querying catalog entries.
       │     * Reporting/Querying lineage graphs.
       │     * Writing/Querying annotations (supporting filtering by type, source, data ref, etc.).
       │     * Managing quality rules/results.
       ├─ **Storage Backend Strategy:**
       │   • **Core Metadata (Catalog, Lineage Graph Structure):** Can reside within the primary `Cluster State Store` (Raft/KV store) for consistency with cluster operations.
       │   • **Annotation / Detailed Metadata:** For potentially high-volume annotations or rich metadata, consider an abstracted storage layer. Initially, could use the main State Store, but allow plugging in a dedicated database (e.g., Document DB for flexible schemas, Graph DB for complex relationships, Search Index like OpenSearch for fast querying) if annotation volume requires it. This avoids bloating the core Raft state.
       └─ Integration: `Worker Slices` report metadata via API; users query via `API Server`.

II. Runi Agent (Daemon on each Worker Node)
    │
    ├─ Role: Runs on every compute node managed by the Runink cluster. Acts as the local representative of the Control Plane. Analogous to Kubernetes Kubelet.
    ├─ Responsibilities:
    │   • Node Registration & Health: Registers the node with the Control Plane, reports resources (CPU, memory, etc.), and sends regular heartbeats.
    │   • Worker Lifecycle Management: Receives commands from the `Scheduler` (via `API Server`) to start, stop, and monitor `Worker Slice Processes`.
    │   • Local Slice Execution: Uses internal `Slice Manager` logic to:
    │     * Create cgroups for strict resource quotas (CPU, memory) per slice.
    │     * Create namespaces for isolation (PID, Net, Mount, User, UTS, IPC). Maps Service Account identity to ephemeral local UID via user namespace.
    │     * Launch worker processes using `os/exec`.
    │   • Secret Delivery: Securely retrieves secrets from the `Secrets Manager` and makes them available to the worker slice (env var, file mount).
    │   • Networking: Configures network interfaces within the slice's network namespace, manages connections for inter-slice communication (gRPC). Potentially enforces basic network policies locally.
    │   • Storage: Manages mounting necessary storage (e.g., shared FS paths, perhaps encrypted local tmpfs) into the slice's mount namespace. Needs integration for encryption at rest if handling local ephemeral data.
    │   • **GPU Awareness (Optional):** Can report presence and status of GPUs to the `Scheduler` for placement of GPU-intensive tasks (like local LLM inference).
    │   • Observability Collection:
    │     * Collects structured logs from worker stdout/stderr. Forwards logs to a configured Fluentd endpoint/aggregator.
    │     * Scrapes Prometheus metrics from the worker's exposed HTTP endpoint. Exposes aggregated node + worker metrics for scraping by a central Prometheus server.
    └─ Security: Authenticates to the API Server using node-specific credentials/tokens. Ensures local operations are authorized.

III. Runi Slice (Isolated Go Process)
     │
     ├─ Role: Executes a single step of a user's data pipeline within a secure, isolated environment managed by the `Runi Agent`.
     ├─ Environment:
     │   • Runs as an ephemeral, non-privileged user (via user namespace) corresponding to a specific Runink `Service Account`.
     │   • Constrained by cgroups (CPU, memory quotas).
     │   • Isolated via Linux namespaces (sees own PID space, network stack, filesystem view).
     │   • Provided with necessary secrets and configuration via env vars or mounted files.
     │   • Issued short-lived tokens for authenticated communication.
     ├─ Responsibilities:
     │   • Executes the specific Go function/code for the pipeline step.
     │   • Communicates with other Runi Slices (potentially on other nodes) via authenticated, encrypted gRPC over the network established by the `Runi Agent`.
     │   • Interacts with external data stores (DBs, streams, files) using credentials provided securely.
     │   • Executes step logic (e.g., data transformation, quality check, LLM interaction).
     │   • **LLM Interaction Steps:**
     │     * Securely fetches data (e.g., from MinIO via S3 API using provided credentials).
     │     * Calls external (OpenAI) or internal (e.g., served from MinIO/GPU node) LLM APIs using secure credentials/network config.
     │     * Parses LLM responses.
     │   • Reports lineage data (inputs, outputs, transform details) to the`Data Governance Service` via gRPC.
     │     * **Reports LLM Annotations:** Formats generated annotations according to the defined schema and sends them to the `Data Governance Service` via its gRPC API.
     │   • Communicates with other slices via authenticated/encrypted gRPC.
     │   • Exposes Prometheus metrics endpoint.
     │   • Writes structured logs to stdout/stderr by the`Runi Agent`.
     │   • Data I/O: Primarily uses gRPC with TLS for inter-slice communication. May interact with shared filesystems mounted into its namespace by the `Runi Agent`.
     └─ Data I/O & Security: Requires appropriate RBAC permissions and Secrets (API keys, internal service URLs) to access data sources (MinIO) and LLM endpoints. Network policies may be needed to allow egress to external APIs or internal LLM services.
```

## Benefits & Tradeoffs

* **Potential Performance:** By avoiding full container runtimes and using direct `exec` with namespaces/cgroups, startup times and per-slice overhead *could* be lower than Kubernetes pods. Go's efficiency helps.
* **Full Control & Customization:** Tailored specifically for Go-based data workloads.
* **Huge Complexity:** Building and maintaining a reliable, secure distributed operating system/scheduler is immensely complex. Replaces Slurm, competes conceptually with Kubernetes/Nomad.
* **Ecosystem:** Lacks the vast ecosystem, tooling, and community support of Kubernetes or HPC schedulers.
* **Security Burden:** All security aspects (RBAC, secrets, TLS, identity, namespace config) must be implemented correctly and securely from scratch.

## Future LLM Integration

1. **Pipeline Definition:** A user defines a pipeline step specifically for LLM annotation. This step specifies the input data source (e.g., path on shared FS/MinIO), the target LLM (e.g., OpenAI model name or internal service endpoint), the prompt, and potentially the output format.
2. **Scheduling:** The `Scheduler` assigns this step to a `Runi Agent`. If targeting an internal LLM requiring specific hardware (GPU), the scheduler uses node resource information (reported by Agents) for placement.
3. **Execution:** The `Runi Agent` launches a `Worker Slice Process` for this step.
4. **Credentials:** The Worker Slice receives necessary credentials (e.g., OpenAI API key, MinIO access key) securely via the `Secrets Manager`.
5. **LLM Call:** The worker reads input data, constructs the prompt, calls the relevant LLM API (external or internal), potentially handling batching or retries.
6. **Metadata Persistence:** Upon receiving results, the worker extracts the annotations, formats them according to the `Data Governance Service` schema, and sends them via gRPC call to the service, linking them to the input data reference. It also reports standard lineage (input data -> LLM step -> annotations).
7. **Usage:** Downstream pipeline steps or external users can then query the `Data Governance Service` (via `API Server`) to retrieve these annotations for further processing, reporting, or analysis.