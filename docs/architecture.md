# Runink Architecture: Go/Linux Native Distributed Data Environment

Goal: Define a self-sufficient, distributed environment for orchestrating and executing data pipelines using Go and Linux primitives. This system acts as the cluster resource manager and scheduler (replacing Slurm), provides Kubernetes-like logical isolation and RBAC, integrates data governance features, and ensures robust security and observability. It aims for high efficiency by avoiding traditional virtualization or container runtimes like Docker. Define a self-sufficient, distributed environment for orchestrating and executing data pipelines using Go and Linux primitives, with enhanced metadata capabilities designed to support standard data governance (lineage, catalog) AND future integration of LLM-generated annotations.


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

*  **User Interaction:** Client requests are often scoped to a specific Herd.
*  **API Server / RBAC:** Enforces RBAC policies based on user/service account permissions within a target Herd.
*  **Cluster State Store:** Explicitly stores Herd definitions and their associated resource quotas.
*  **Scheduler:** Considers Herd-level quotas when making placement decisions.
*  **Secrets Manager:** Access to secrets might be scoped by Herd.
*  **Data Governance:** Metadata (lineage, annotations) can be tagged by or associated with the Herd it belongs to.
*  **Runi Agent:** Receives the target Herd context when launching a slice and uses this information to potentially configure namespaces and apply appropriate cgroup limits based on Herd quotas.
*  **Runi Slice:** A single instance of a pipeline step running as an isolated Worker Slice Process. Executes entirely within the logical boundary and resource constraints defined by its assigned Herd.
* **Runi Pipes:** Primarily used now for internal communication within the `Runi Agent` to capture logs/stdio from the `Runi Slice` it `exec`s, rather than for primary data transfer between steps.
* **Herd:** A logical grouping construct, similar to a Kubernetes Namespace, enforced via RBAC policies and potentially mapped to specific sets of Linux namespaces managed by Agents. Provides multi-tenancy and team isolation. Quotas can be applied per Herd.
* **Go Native & Linux Primitives:** Core components written in Go, directly leveraging cgroups, namespaces (user, pid, net, mount, uts, ipc), pipes, sockets, and `exec` for execution and isolation.
* **Self-Contained Cluster Management:** Manages a pool of physical or virtual machines, schedules workloads onto them, and handles node lifecycle.
* **Serverless Execution Model:** Users define pipelines and resource requests; Runink manages node allocation, scheduling, isolation, scaling (by launching more slices), and lifecycle. Users are subject to quotas managed via cgroups.
* **Security First:** Integrated identity (OIDC), RBAC, secrets management, network policies, encryption in transit/rest.
* **Data Governance Aware:** Built-in metadata tracking, lineage capture, and support for quality checks. With extension for storage/management of rich data annotations (e.g., from LLMs).
* **Rich Observability:** Native support for metrics (Prometheus) and logs (Fluentd).


---

Below is a **technology-focused benchmark** you can embed into your `architecture.md` (or a dedicated benchmarking section) to highlight the **architectural gains** Runink achieves by choosing Go, functional pipelines, Linux primitives, and Raft—rather than simply comparing tools:

---

## 📈 Technology Benchmark: Go + Raft + Linux Primitives vs. JVM + Containerized Stacks

| **Dimension**                          | **Traditional (JVM / Container / DataFrame)**                   | **Runink (Go / FP Pipelines / Raft)**                                                   | **Gains**                                                                                                  |
|----------------------------------------|-----------------------------------------------------------------|-----------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| **Language Runtime**                   | Java Virtual Machine                                            | Go runtime                                                                                | • **Lower startup latency** (milliseconds vs. seconds) • **Smaller memory footprint** per process (tens of MBs vs. hundreds of MBs)       |
| **Programming Paradigm**               | OOP (classes, inheritance, heavy type hierarchies)             | Functional-ish pipelines in Go (pure functions, composable stages)                      | • **Simpler abstractions** (no deep inheritance)  • **Easier reasoning** about data flow  • **Higher testability** via pure stage functions      |
| **Data Processing Model**              | Row-based DataFrames             | Channel-based streams of structs (Go channels → functional transforms)                   | • **Constant memory** for streaming (no full in-memory tables)  • **Back-pressure** and windowing naturally via channels  • **Composable, step-by-step transforms**                 |
| **Isolation & Sandboxing**             | Containers + VM overhead                     | cgroups + namespaces + Go processes                                                      | • **0.5–1 ms** slice startup vs. **100 ms+** container spin-up  • **Fine-grained resource limits** per slice  • **No dockerd overhead**                              |
| **Inter-Stage Communication**          | Shuffle via network storage or distributed file systems        | In-memory pipes, UNIX sockets, or gRPC streams within Go                                | • **<1 ms** handoff latency  • **Zero-copy** possible with `io.Pipe`  • **Strong type safety** end-to-end                      |
| **Distributed Consensus**               | External etcd or no coordination    | Built-in Raft                                                    | • **Linearizable consistency** for control plane  • **Leader election** and automatic recovery  • **Atomic updates** across scheduler, secrets, and metadata |
| **Multi-Tenancy & Security**           | K8s namespaces + complex RBAC, network policies               | Herd namespaces + cgroup quotas + OIDC/JWT + mTLS                                       | • **Single-layer isolation**  • **Per-slice ephemeral UIDs** for zero-trust  • **Simpler, declarative Herd-scoped policies**               |
| **Observability & Lineage**            | External stacks                  | Native agents + Raft-consistent metadata store                                          | • **Always-consistent** lineage  • **Unified telemetry** (metrics, logs, traces)  • **Real-time auditability** without manual integration       |
| **Deployment & Ops**                   | Helm charts, CRDs, multi-tool CI/CD                           | Single Go binary + Raft cluster bootstrap                                              | • **One artifact** for all services  • **Simpler upgrades** via rolling Raft restarts  • **Built-in health & leader dashboards**                  |
| **Functional Extensibility**           | Plugins in Java/Python with JNI or UDFs (heavy)               | Go plugins or simple function registration                                              | • **No cross-language bridges**  • **First-class Go codegen** from DSL  • **Easier on-the-fly step injection**                         |

---

### Key Takeaways

1. **Speed & Efficiency**  
   Go’s minimal runtime and direct use of Linux primitives deliver sub-millisecond task launches and minimal per-slice overhead—vs. multi-second container or JVM startups.

2. **Deterministic, Fault-Tolerant Control**  
   Raft gives Runink a single, consistent source of truth across all core services—eliminating split-brain and eventual-consistency pitfalls that plague layered stacks.

3. **Composable, Functional Pipelines**  
   Channel-based, stage-oriented transforms allow lean, testable, streaming-first ETL—rather than bulk-loading entire tables in memory.

4. **Unified, Secure Multi-Tenancy**  
   Herds + namespaces + cgroups + OIDC/JWT yield strong isolation and a simpler security model than stitching together K8s, Vault, and side-cars.

5. **Built-In Governance & Observability**  
   Raft-backed metadata ensures lineage and schema contracts are always in sync—no external governance platform required.

By grounding Runink’s design in **Go**, **functional pipelines**, **Linux primitives**, and **Raft consensus**, you get a single, vertically integrated platform that outperforms and out-operates the conventional “JVM + containers + orchestration + governance” stack—delivering **predictable**, **low-overhead**, and **secure** data pipelines at scale.

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