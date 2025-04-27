# 📚 Project Structure Reference

This document provides a detailed reference for every folder and file within the **Runink** CLI project. It will be used to guide development, onboarding, and maintenance.

---

## Top-Level Files

| File | Purpose |
|:---|:---|
| `Makefile` | Build, run, install, lint, clean, and test project easily from the command line |
| `go.mod` | Go module definition (package name, dependencies) |
| `go.sum` | Go dependency checksums |
| `main.go` | Root CLI command. Wires all domain commands into `runi` |
| `scripts/` | Shell or helper scripts (setup, migrations, local dev bootstrap) |

---

## `barnctl/` — 🌾 Herd Control Plane

| File | Purpose |
|:---|:---|
| `cli.go` | Cobra command root for `barnctl` subcommands |
| `audit.go` | Audit logs for herd and secret actions |
| `checker.go` | Validation helpers for herd definitions and requests |
| `client.go` | API client for herd control plane operations |
| `crypto.go` | Secret encryption/decryption (AES-GCM) |
| `errors.go` | Domain-specific error types |
| `handlers.go` | API handlers (gRPC/HTTP) for herd operations |
| `herd_create.go` | Create a new herd (namespace, quotas) |
| `herd_delete.go` | Delete an existing herd |
| `herd_update.go` | Update herd quotas, policies |
| `keys.go` | Key and secret management helpers |
| `logs.go` | Herd logs management (event sourcing, rotation) |
| `manager.go` | Main herd manager, orchestration layer |
| `metadata.go` | Herd metadata structs and conventions |
| `middleware.go` | HTTP/gRPC middleware (auth, validation) |
| `models.go` | Data models for herd, secrets, registration |
| `options.go` | CLI flag option structs for `barnctl` |
| `rbac.go` | RBAC policy enforcement for herds |
| `reencrypt.go` | Secret re-encryption and key rotation logic |
| `register.go` | Herd and agent registration logic |
| `registry.go` | Metadata registry management (Raft-backed) |
| `response.go` | Standard API response structs |
| `router.go` | HTTP router wiring all handlers |
| `runi_create.go` | Create Runi slices from herd context |
| `runi_delete.go` | Delete a Runi slice |
| `runi_update.go` | Update a Runi slice |
| `sbom.go` | Generate and manage SBOM (Software Bill of Materials) |
| `server.go` | HTTP and gRPC server lifecycle |
| `snapshot.go` | Create/restore herd snapshots (metadata durability) |
| `storage.go` | Low-level secret storage interactions |
| `store.go` | High-level Store abstraction (Badger, Raft-backed) |
| `tracing.go` | OpenTelemetry tracing for herd ops |
| `utils.go` | Small utilities (hashing, ID generation) |
| `validate.go` | Schema and data validation helpers |
| `validate_api.go` | Request/response validation for HTTP API |

---

## `buildctl/` — 🛠️ Compiler

| File | Purpose |
|:---|:---|
| `cli.go` | Cobra command root for `buildctl` subcommands |
| `builder.go` | Build DAGs from DSL and Contract inputs |
| `compliance.go` | Inject SLA/compliance metadata into DAGs |
| `contract.go` | Contract loader and parser |
| `dag.go` | DAG structures and traversal utilities |
| `feature.go` | Feature DSL parsing and struct models |
| `graphviz.go` | Graphviz `.dot` file export of DAGs |
| `hooks.go` | Compiler event hooks (pre/post build) |
| `metadata.go` | DAG metadata tagging |
| `optimizer.go` | DAG optimization (chaining, streaming) |
| `parser.go` | DSL parsing engine |
| `render.go` | Render compiled DAG to Go code |
| `step.go` | Step representations inside DAGs |
| `steps_dsl.go` | DSL step parsing helpers |
| `stream.go` | Streaming pipeline render engine |
| `utils.go` | Helper functions |
| `validate.go` | Validate DAG integrity and feature compliance |
| `validate_dsl.go` | Validate DSL syntax and semantics |

---

## `herdctl/` — 🐑 Metadata Governance

| File | Purpose |
|:---|:---|
| `cli.go` | Cobra command root for `herdctl` subcommands |
| `api.go` | gRPC/HTTP handlers for herd governance |
| `approvals.go` | Manage contract/feature approval flows |
| `changelog.go` | Generate changelogs between schema/feature versions |
| `exporter.go` | Export metadata, lineage, or artifacts |
| `graph.go` | Lineage graph builders (DAG of runs) |
| `hashing.go` | Contract fingerprinting and hashing helpers |
| `herdloader.go` | Load herd metadata for governance |
| `hooks.go` | Event hooks for governance actions |
| `lineage.go` | Lineage emission and management |
| `metadata.go` | Metadata formats and versioning |
| `metrics.go` | Emit governance Prometheus metrics |
| `policy.go` | Manage herd-wide RBAC policies |
| `rdf.go` | Export metadata as RDF triples |
| `redactions.go` | Handle PII field redactions |
| `scanner.go` | Governance rules engine (compliance checks) |
| `scopes.go` | Manage metadata scopes for herds |
| `semantic_query.go` | Semantic search over lineage/metadata |
| `signer.go` | Artifact signing |
| `token.go` | Short-lived token minting (JWT) |
| `utils.go` | Helper utilities |
| `validate.go` | Validate governance metadata integrity |

---

## `runictl/` — 🚀 Slice Scheduling & Execution

| File | Purpose |
|:---|:---|
| `cli.go` | Cobra command root for `runictl` subcommands |
| `affinity.go` | Slice affinity rules (node selection) |
| `agent.go` | Agent lifecycle management (heartbeat, registration) |
| `cgroup.go` | Cgroup enforcement for slices |
| `constraints.go` | Scheduling constraint validation |
| `envloader.go` | Load environment variables into slices |
| `execution.go` | Run slice DAGs (task orchestration) |
| `heartbeat.go` | Agent heartbeat emission |
| `lineage.go` | Slice-level lineage tracking |
| `metadata.go` | Metadata enrichment for runs |
| `monitor.go` | Slice SLA monitoring |
| `namespace.go` | Namespace (user/mount) creation for slices |
| `parser.go` | Request parser for slices |
| `pipes.go` | Unix pipe chaining between slice stages |
| `placement.go` | Node scoring and placement logic |
| `preemption.go` | Slice preemption under resource pressure |
| `runner.go` | Low-level slice execution engine |
| `sandbox_agent.go` | Agent running in isolated Linux namespaces |
| `scheduler.go` | Core scheduler loop |
| `secrets.go` | Secrets injection into slice env |
| `solver.go` | Constraint solver for scheduler |
| `utils.go` | Retry, error wrapping utilities |
| `validate_dag.go` | Validate DAGs before execution |
| `validate_request.go` | Validate slice execution requests |

---

## `docs/`

| File | Purpose |
|:---|:---|
| `architecture.md` | Explain Runink architecture (slices, herds, control plane) |
| `developer-guide.md` | Developer workflows, contributing guide |
| `governance.md` | Metadata governance flows, approvals, versioning |

---

## `internal/`

| Folder | Purpose |
|:---|:---|
| `api/` | Internal protobuf models, gRPC utils |
| `auth/` | Authentication modules (OIDC/JWT) |
| `config/` | CLI configs parsing (env, yaml, flags) |
| `telemetry/` | OpenTelemetry tracing, metrics setup |

---

> This document is updated every major refactor or package addition.

