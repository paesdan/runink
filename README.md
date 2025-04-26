# Runink: Cloud Native Distributed Data Environment

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://example.com/build/status)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/your_org/runink)](https://goreportcard.com/report/github.com/your_org/runink)

---

## Overview

<img src="./runink-site/content/docs/images/logo.png" width="150"/>

***Runink*** is an ambitious project aiming to define a **self-sufficient, distributed environment** specifically designed for orchestrating and executing data pipelines. Built natively in **Go** and leveraging core **Linux primitives** (cgroups, namespaces, `exec`), Runink acts as its own cluster resource manager and scheduler, providing a vertically integrated platform that replaces the need for separate systems like Slurm or complex Kubernetes setups for data workloads.

Our goal is to provide a highly efficient, secure, and governance-aware platform with a **serverless execution model** for data engineers and scientists. Define your pipelines declaratively, and let Runink handle the distributed execution, isolation, resource management, security, lineage, and observability.

## Core Principles

* **Go Native & Linux Primitives:** Maximize performance and minimize overhead by using Go and direct Linux kernel features (cgroups, namespaces, pipes, sockets, `exec`).
* **Self-Contained Cluster Management:** Manages a pool of nodes, schedules workloads (`Runi Slices`), and handles node lifecycles.
* **Serverless Execution Model:** Users define pipelines; Runink manages resource allocation (within quotas), execution isolation (`Herds`/`Runi Slices`), and scaling.
* **Security First:** Integrated identity (OIDC), Role-Based Access Control (RBAC) scoped by `Herd`, secrets management, network policies via namespaces, and end-to-end encryption (TLS for gRPC).
* **Data Governance Aware:** Built-in metadata catalog, automatic data lineage tracking, data quality framework, and support for rich annotations (including LLM-generated metadata).
* **Rich Observability:** Native support for metrics (Prometheus format) and structured logs (for Fluentd aggregation).

## Key Features

* **Feature DSL:** Define complex pipelines using a human-readable, Gherkin-inspired `.dsl` file format.
* **Integrated Orchestration:** A built-in Control Plane parses DSLs into DAGs and schedules distributed execution across managed nodes.
* **Runi/Herd Execution Model:** Lightweight, isolated pipeline steps (`Runi Slices`) execute within resource-constrained cgroups and logically separated `Herd` namespaces for multi-tenancy and domain isolation.
* **Data Governance & Lineage:** Automatic lineage tracking, integrated data catalog, quality rule management, and rich annotation support (incl. LLM outputs) via a central `Data Governance Service`.
* **Built-in Security:** Native RBAC, OIDC integration, secrets management, mTLS for internal communication, and namespace-based isolation.
* **Native Observability:** Components are instrumented for Prometheus metrics and structured logging out-of-the-box.
* **Schema Contracts:** Define and enforce data structure guarantees throughout the pipeline lifecycle.

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

*(For more details, see [`docs/architecture.md`](./docs/architecture.md))*

## Key Concepts

---

<table>
  <tr>
    <th><img src="./docs/images/runink.png" width="450"/></th>
    <th><h4>The golang code base to deploy features from configurations files deployed by command actions over the CLI/API.</h4></th>
  </tr>
  <tr>
    <th>Runink</th>
  </tr>
</table>
<br>
<table>
  <tr>
    <th><img src="./docs/images/runi.png" width="450"/></th>
    <th><h4>A single instance of a pipeline step running as an isolated <i>Runi Slice Process</i> managed by a <i>Runi Agent</i> within the constraints of a specific <i>Herd</i></h4></th>
  </tr>
  <tr>
    <th>Runi</th>
  </tr>
</table>
<br>
<table>
  <tr>
    <th><img src="./docs/images/herd.png" width="450"/></th>
    <th><h4>A logical grouping construct, similar to a Kubernetes Namespace, enforced via RBAC policies and resource quotas. Provides multi-tenancy and domain isolation.</h4></th>
  </tr>
  <tr>
    <th>Herd</th>
  </tr>  
</table>

<br>

---

# Runink Component Reference

This section provides a brief overview of the main Runink components. For detailed descriptions, please refer to [`components`](/docs/components/) page.

---

### ⚙️ Runink CLI (`runi`)

> Role: Developer interface for scaffolding, testing, and running pipelines.

The primary command-line tool for interacting with the Runink platform, managing contracts, features, tests, and deployments.

---

### 🤖 Runi Agent

> Role: Node daemon managing local slice execution and reporting.

Runs on each compute node, manages `Runi Slice Process` lifecycle within cgroups/namespaces, fetches secrets, and forwards logs/metrics to the observability stack.

---

### 🐑 Herd Namespace

> Role: Logical isolation boundary for multi-tenancy and domain separation.

A logical grouping (similar to K8s namespaces) enforced via RBAC and resource quotas, providing isolation for `Runi Slices` using Linux namespaces.

---

### 🖥️ Herd Control Plane

> Role: Central orchestration and management brain of the Runink cluster.

The distributed set of services (API Server, Scheduler, State Store, etc.) managing cluster state, scheduling workloads, enforcing policies, and providing APIs.

---

### 💾 Barn (Cluster State Store)

> Role: Distributed, consistent storage for all cluster state.

The authoritative, HA key-value store (backed by Raft) holding configurations, pipeline state, RBAC policies, secrets metadata, and core governance info.

---

### 🔐 Secrets Manager

> Role: Secure vault for managing and delivering sensitive credentials.

Stores encrypted secrets (scoped by Herd/RBAC) and securely delivers them to `Runi Agents` for use by `Runi Slices`.

---

### 🔑 Identity & RBAC Manager

> Role: Manages authentication and authorization policies.

Handles user/service identities (OIDC), defines and enforces fine-grained, Herd-scoped RBAC policies via the API Server.

---

### 🖥️ API Server

> Role: Main gateway for all interactions with the Runink platform.

Exposes gRPC/REST APIs, handles AuthN/AuthZ, validates requests, and routes commands to other control plane services.

---

### 📑 Schema Contracts

> Role: Defines authoritative data structures, versions, and validations.

Versioned definitions of data structure used for validation within pipelines and tracked via the governance service.

---

### 🧱 Feature DSL (Domain-Specific Language)

> Role: Human-readable language for defining pipeline logic.

A Gherkin-inspired format (`.dsl` files) used to declare pipeline steps and scenarios, parsed by the Control Plane.

---

### 🧪 Testing Engine

> Role: Ensures pipeline correctness and prevents regressions.

Supports golden file testing, synthetic data generation, and integration tests, driven via the `Runink CLI`.

---

### 🔍 Interactive REPL

> Role: Interface for live data exploration, debugging, and prototyping.

A command-line shell combining DataFrame API, JSON navigation, and SQL-like querying for interactive work.

---

### 🚧 Pipeline Generator

> Role: Translates Feature DSL into an executable DAG plan.

Conceptually, the Control Plane component (Scheduler) that interprets the DSL and plans the execution DAG for `Runi Slices`.

---

### 🔄 Feature Orchestration

> Role: Handles pipeline dependency resolution and execution scheduling.

The core scheduling and execution logic managed by the `Herd Control Plane Scheduler` and `Runi Agents`, executing steps as `Runi Slices` within `Herd` boundaries.

---

### 🧭 Data Lineage & Metadata

> Role: Central service for tracking data flow, context, and annotations.

Managed by the `Data Governance Service`, providing APIs to record and query lineage graphs, catalog info, annotations (incl. LLM), and quality results.

---

### ✅ Data Quality & Validation

> Role: Maintains data integrity and compliance throughout the pipeline.

Enforces rules defined in `Schema Contracts` and `Feature Domain Structured Language` within `Worker Slices`, with results tracked by the `Data Governance Service`.

---

### 🔒 Data Governance & Compliance

> Role: Ensures adherence to regulatory and internal data governance standards.

Encompasses RBAC (scoped by `Herd`), lineage tracking, audit trails, and security controls managed across the Control Plane and Agents.

---

### 🛡️ Security & Enterprise DataOps

> Role: Integrates security practices into the pipeline lifecycle.

Combines platform security features (RBAC, secrets, encryption, isolation via namespaces/cgroups) with recommended CI/CD practices for the platform code itself.

---

### 🌐 Observability & Monitoring

> Role: Enables comprehensive visibility into pipeline health and performance.

Provides metrics (Prometheus) and logs (Fluentd) collected by `Runi Agents` from `Runi Slices` and Agents themselves.

---

## Getting Started

*(TODO: Add instructions for installing the Runink CLI, setting up a minimal cluster (Herd Control Plane + Runi Agents), and running a basic example pipeline using the Feature DSL.)*

```bash
# Example (Conceptual)
runi herd create my-data-herd
runi test \
  --scenario features/trade_cdm.dsl \
  --golden golden/stream_and_enrich \
  --herd finance
runi apply \
  --scenario features/trade_cdm.dsl \
  --contract contracts/trade_cdm_multi.go \
  --herd finance \
  --out rendered/trade_cdm_dag.go
runi run \
  --scenario features/trade_cdm.dsl \
  --herd finance \
  --input testdata/stream_and_enrich/input.json
runi status --runid RUN-20240424-XYZ --herd finance
```

## Development Status

**Alpha / Conceptual:** Runink is currently under active development and should be considered experimental. The architecture and features described represent the target state.


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](./LICENSE) file for details.
