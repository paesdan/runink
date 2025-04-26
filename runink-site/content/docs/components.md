# Runink Component Reference

This document describes each component of the Runink architecture, providing clear definitions, roles, and interactions within the overall system.

<img src="./images/components.png" width="680"/>

---
### ⚙️ Runink CLI

> Role: Developer interface for scaffolding, testing, and running pipelines.  

**The Runink CLI (`runi`)** is the developer-first control surface for the entire data lifecycle — from contracts and pipelines to tests, observability, and deployment. It acts as the **orchestration layer for humans**, offering a simple, scriptable, and powerful interface for everything Runink can do.

Whether you're scaffolding a new project, compiling `.dsl` files into Go pipelines, validating schema contracts, testing with golden files, or pushing artifacts to a metadata registry — `runi` gives you the tools to do it fast and reliably, right from your terminal.

***Functions:***

- **Project Scaffolding** for quickly bootstrapping new repositories with features, contracts, tests, and CI configs  
- **Contract Management** including generation from Go structs, diffing, rollback, and drift tracking  
- **Pipeline Compilation** from `.dsl` scenarios into DAG-aware Go code with batch and streaming support  
- **Golden Testing & Synthetic Data** using built-in Faker-based generation and contract-backed verification  
- **Metadata Publishing** to push lineage, changelogs, and SBOMs to registries or dashboards  
- **Interactive REPL & Exploration** for debugging, jq-style filtering, DataFrame operations, and live SQL queries  
- **Security & DevOps** integrations for SBOM, artifact signing, and secure publishing pipelines  

For a full list of commands and examples, see [`cli`](/docs/cli).  

The CLI embeds core functionality like the **Pipeline Generator**, **Execution Engine**, and **Governance Layer**, turning complex distributed operations into single, declarative, CLI-driven workflows.

With `runi`, developers get a **portable, scriptable, and audit-friendly control plane** — no YAML, no drag-and-drop UIs, just pipelines that flow from contract to cluster with confidence.


---

### 🤖 Runi Agent

**Runi** (the Runink Agent) runs on every compute node, bridging the Control Plane and the isolated worker slices. It registers node resources and health, fetches secrets, and manages the lifecycle of pipeline steps within lightweight, cgroup‑isolated “slices.” Runi ensures that each step executes in the proper Herd context, collects observability data, and enforces resource quotas.

***Functions:***

- **Node Registration** with Runink Control Plane (health, capacity, labels)  
- **Slice Manager** that creates cgroups and namespaces per Herd, launches worker processes, and applies resource limits  
- **Secrets Fetching** via secure Raft‑backed Secrets Manager  
- **Observability Forwarding** metrics to Prometheus, logs to Fluentd, traces to OpenTelemetry  
- **Lifecycle Management** handles retries, restarts, graceful shutdowns, and cleanup  

Runi turns each compute node into a **secure**, **observed**, and **policy‑enforced** worker — ready to run any pipeline slice on demand.

---

### 🐑 Herd Namespace

**Herd Namespaces** is Runink’s lightweight domain isolation layer, leveraging Linux namespaces to enforce **shared‑nothing** principles. Each “herd” represents a logical group of related pipeline steps or data slices, isolated in its own namespace to guarantee resource control, security boundaries, and reproducibility.

***Functions:***

- **Domain Isolation:** Creates per‑herd namespaces so pipelines for different domains (e.g., finance, ecommerce) cannot interfere  
- **Resource Control:** Assigns cgroups for CPU, memory, and I/O limits on each herd’s worker slices  
- **Secure Multi‑Tenancy:** Ensures data and processes are confined to their designated namespaces  
- **Dynamic Scaling** Spins up and tears down herds on demand, coordinating with bigmachine for serverless‑style compute  
- **Observability Hooks** Exposes namespace‑scoped metrics and logs, tagged by herd ID and domain  
- **Multi‑Tenant Governance** separate catalogs, DLQs, and audit trails per Herd  

With Herd, Runink achieves fine‑grained isolation for domain-driven workloads — combining security, performance, and multi‑tenant governance in a single, native Linux construct. Empowering teams to collaborate on a shared platform without stepping on each other’s toes.

---

### 🖥️ Herd Control Plane

**Herd Control Plane** is the high‑availability control plane that glues together pipelines, policies, and metadata. It serves as the central orchestrator and gateway for all user and agent interactions, enforcing authentication, authorization, and multi‑tenant (Herd) scoping. Runink stores and manages the cluster state, schedules work, and provides metadata and governance APIs — all over secure, TLS‑encrypted gRPC/REST endpoints.

***Functions:***

- **API Server** for CLI/UI/SDK requests (authN/authZ via OIDC, RBAC checks per Herd)  
- **Barn** (Cluster State Store under Raft consensus) for pipelines, Herd definitions, RBAC policies, and secrets  
- **Scheduler** that matches pipeline tasks to available nodes, respecting resource quotas and Herd affinity  
- **Secrets Manager** that securely stores and delivers credentials scoped to services and Herds  
- **Data Governance Service** for lineage, quality, and annotation APIs  
- **High‑Availability** via stateless replicas and consensus‑backed state  
- **CLI Hub:** Exposes commands (`runi init`, `runi contract …`, `runi compile`, `runi test`, `runi run`, `runi publish`) for end‑to‑end pipeline management  
- **Project Scaffolding:** Bootstraps new repositories with templates, configs, contracts, and CI/CD pipelines  
- **Feature & Contract Integration:** Parses Feature DSL and schema contracts to drive code generation, testing, and execution  
- **Metadata Registry:** Manages lineage, versioned contracts, run histories, and observability endpoints  
- **Plugin & Extension Manager:** Loads and invokes user‑provided plugins, adapters, and custom steps  

With Herd at the helm, your platform is **secure**, **consistent**, and **multi‑tenant ready** — enforcing policies and providing a single source of truth for all pipelines and metadata. Every aspect of your data platform—from initial spec to production rollout—is **standardized, traceable, and automatable**, so teams can focus on delivering insight rather than wiring infrastructure.

---

### 💾 Barn (Cluster State Store)

**The Barn** is the authoritative, highly available backbone of Runink’s control plane — a **distributed key-value store backed by Raft consensus**. It provides a consistent, fault-tolerant source of truth for all runtime metadata, configurations, and orchestration logic across the platform.

This includes pipeline definitions, Herd boundaries, secrets metadata, RBAC policies, and scheduler state. Every component of the Runink platform — from the API Server to the Scheduler and Governance Layer — relies on this store to read, write, and replicate configuration and execution state. With strong consistency guarantees and built-in fault tolerance, the state store ensures that decisions made by the platform are **safe, synchronized, and auditable**.

***Functions:***

- Persist Herd definitions, RBAC policies, pipeline graphs, and secrets metadata  
- Provide read/write access to the API Server, Scheduler, and Runi Agents  
- Synchronize cluster-wide state using Raft-based consensus  
- Support high-availability deployments with leader election and quorum  
- Enable fault-tolerant execution of distributed pipelines and governance policies  

With the Barn, Runink gains a **single source of truth** that’s distributed by design — making orchestration consistent, secure, and always in sync.

---

### 🔐 Secrets Manager

**The Secrets Manager** is Runink’s secure vault for managing sensitive credentials, tokens, keys, and configuration values across environments and Herds. It ensures that secrets are **stored safely, accessed securely, and delivered precisely when and where they're needed** — with no manual handling required by pipeline authors.

Backed by encrypted storage within the **Barn**, the Secrets Manager enforces strict **per-Herd access controls**, ensuring that only authorized agents and services can retrieve secrets for a given pipeline or slice. Secrets are injected just-in-time into ephemeral slice environments, avoiding persistence or leakage risks.

Whether you’re integrating with databases, APIs, or cloud services, the Secrets Manager ensures that sensitive data stays **encrypted, auditable, and access-controlled** by default.

***Functions:***

- Securely store secrets in an encrypted, Raft-backed data store  
- Scope secret access using RBAC policies and Herd-level permissions  
- Deliver secrets to Runi Agents on demand for authorized slices only  
- Support dynamic secret injection into pipeline execution environments  
- Ensure secrets are never written to disk or logs during runtime  

With the Secrets Manager, Runink provides **seamless and secure secret management** — enabling safe access to credentials without compromising auditability or isolation.

---
### 🔑 Identity & RBAC Manager

**The Identity & RBAC Manager** is the guardian of **who can do what, where, and when** across the Runink platform. It provides secure authentication and fine-grained, role-based authorization by linking users, service accounts, and slices to their permitted actions within specific **Herds**.

Integrated with **OIDC providers** (like Google, Okta, or GitHub), it supports Single Sign-On (SSO) and token-based identity. At the API layer, every request is checked against **RBAC policies**, ensuring only authorized users or agents can create pipelines, access secrets, or run code in a given domain.

Each issued token (JWT) is short-lived, scoped to specific actions and Herds, and cryptographically verifiable — enabling secure, auditable access across distributed components.

***Functions:***

- Authenticate users and services via OIDC and JWT tokens  
- Enforce fine-grained, Herd-scoped RBAC policies on all API and CLI actions  
- Bind identities to service accounts, groups, or workload slices  
- Issue, refresh, and revoke secure tokens with scoped permissions  
- Audit access and permission changes across time and space  

With the Identity & RBAC Manager, Runink ensures **secure, auditable, and granular access control** — empowering collaboration without compromising control.

---

### 🖥️ API Server

**The API Server** is the main gateway to the Runink platform — acting as the **front door for developers, services, and tools**. It exposes both **REST** and **gRPC** interfaces to support CLI commands (`runink`), SDKs, and any web-based UI clients, ensuring seamless access from local scripts to enterprise control planes.

All incoming requests pass through a robust authentication and authorization layer, powered by **OIDC** and **RBAC policies**, and are validated against the current cluster state. The API Server routes tasks, triggers executions, resolves Herd permissions, and coordinates interactions with internal services like the Scheduler, Governance layer, and Secrets Manager.

Whether you're registering a pipeline, validating a contract, or executing a scenario — it all begins here.

***Functions:***

- Accept and route requests from CLI (`runink`), UI, SDKs, and APIs  
- Authenticate and authorize users and agents based on OIDC and RBAC  
- Validate inputs, resolve Herd scoping, and check quota policies  
- Proxy commands to internal services (Scheduler, Governance, etc.)  
- Serve both REST and gRPC endpoints for extensibility and automation  

With the API Server, Runink delivers a **secure, scalable, and consistent interface** for interacting with the entire data platform — keeping pipelines and teams in sync.

---

### 📑 Schema Contracts

**Schema Contracts** are the cornerstone of Runink's commitment to data integrity, trust, and evolution. They define the **authoritative structure** of your data — including fields, types, formats, constraints, and relationships — serving as a shared language between data producers, consumers, and pipeline logic.

Generated directly from Go structs, contracts are versioned, validated, and tracked through every stage of a pipeline. When upstream systems change, Runink detects **schema drift**, highlights diffs, and optionally blocks execution or routes data for quarantine — ensuring nothing breaks silently. Contracts also tie directly into your tests, golden files, and metadata lineage, making them the connective tissue between your code, your data, and your governance. Runink treats contracts not just as schemas, but as **data responsibility guarantees** — machine-verifiable, human-readable, and Git-versioned.

***Functions:***

- Auto-generate contracts from Go structs (`runi contract gen`)  
- Detect and diff schema changes with drift alerts  
- Enforce structural and semantic validation at ingest or transform  
- Integrate with golden testing, feature execution, and lineage tracking  
- Embed contract versions in audit logs, metadata, and observability signals  

With Schema Contracts, Runink gives your team **confidence, control, and compliance** — by design.

---


### 🧱 Feature DSL (Domain-Specific Language)

The **Feature DSL** is Runink’s powerful, readable interface for defining pipelines — purpose-built to align data logic with domain language. Inspired by Gherkin, this DSL uses `.dsl` files to declare **what should happen to the data**, using a structured, declarative flow based on real business scenarios.

Every scenario is composed of modular `@step` blocks, organized using `@when`, `@then`, and `@do` to clearly express conditions, expected outcomes, and the logic to apply. These steps can be reused across domains — for example, `@step(name=NormalizeEmails)` may appear in both ecommerce and finance pipelines, while more specific validations (like salary ranges or KYC checks) are scoped to their respective domain modules.

This structure makes it easy to translate business rules into t.mdestable pipelines, enabling close collaboration between technical and non-technical users — with each scenario acting as both documentation and executable logic.

***Functions:***

- Scenarios are composed of reusable, composable steps across domain contexts  
- Each step uses `@when → @then → @do` semantics for clarity and modularity  
- Supports tag-based step handlers (`@source`, `@step(name=…)`, `@sink`)  
- Domain-specific constraints and logic are cleanly isolated in modules  
- DSL scenarios are fully testable, traceable, and convertible to pipeline code  

With the Feature DSL, Runink turns business scenarios into living, executable specifications — bringing pipelines and domain intent closer together than ever before.

---

### 🧪 Testing Engine

> Role: Ensures pipeline correctness and prevents regressions.

The **Testing Engine** is Runink’s built-in safety net — ensuring that every pipeline transformation is **correct, traceable, and regression-proof**. At its heart is the **golden file system**, where each scenario’s expected output is saved and used for future comparisons. If anything changes during a test run, Runink shows a clear and structured diff, making it easy to spot regressions or confirm intentional updates.

Tests can be auto-generated using **synthetic data**, and pipelines can be validated in full using **integration test suites**. Whether you’re writing fine-grained unit tests or verifying complete feature flows, Runink leverages Go’s standard testing ecosystem alongside its own CLI to make pipeline testing reliable and repeatable.

Combined with contracts, scenarios, and metadata tracking, this engine ensures that changes are safe, intentional, and observable — without sacrificing developer speed.

***Functions:***

- Golden file testing with structured diff reporting  
- Synthetic data generation for repeatable, contract-aligned tests  
- CLI-driven integration testing with `.dsl` scenarios  
- Catch schema drift, logic changes, or data regressions early  
- All tests are human-readable and Git-friendly for review  

With Runink’s Testing Engine, quality isn’t bolted on — it’s baked into every scenario, from the first line to the final output.

---

### 🔍 Interactive REPL: DataFrame / SQL Interface / JSON Explorer

The **Interactive REPL** is your live terminal into the data universe — purpose-built for **exploration**, **debugging**, and **interactive prototyping**. Whether you're slicing through a JSON blob, cleaning up a CSV, or building a pipeline from scratch, the REPL gives you fast, expressive control over your data in real time.

Runink’s REPL combines a **DataFrame-inspired API**, **jq-style JSON navigation**, and **SQL-like querying** into a single, cohesive shell. You can apply transformations step by step, inspect intermediate states, and convert your exploratory logic into reusable pipeline components — all without rerunning an entire scenario or writing full Go code. It also supports live previews of how your SQL or REPL commands will map into structured Runink pipeline steps, helping bridge the gap between experimentation and production pipelines.

***Functions:***

- Live scenario prototyping and debugging, one step at a time  
- JSON exploration with jq-style dot access and filtering  
- Interactive DataFrame operations (`Filter`, `Join`, `GroupBy`, etc.)  
- SQL-like query parsing for quick insights and structured filters  
- Preview and export queries into `.dsl` pipeline DSL code  

Whether you’re writing tests, preparing contracts, or validating the latest stream of events — the REPL keeps you **curious**, **confident**, and **connected** to your data.

---

### 🚧 Pipeline Generator

The **Pipeline Generator** is the core of Runink’s automation engine — turning human-readable `.dsl` files into optimized, production-ready Go pipelines. It analyzes each scenario and compiles it into a **Directed Acyclic Graph (DAG)** of transformations, ensuring that all steps are executed in the correct order while honoring their data dependencies.

Each tag (`@source`, `@step`, `@sink`) in the DSL is parsed with context-aware parameters, allowing you to reuse logic across platforms and domains. Whether you're sourcing from Snowflake, applying a standard validation like `@step(name=ValidateSIN)`, or writing to a streaming sink, Runink generates tailored Go code behind the scenes to handle the orchestration.

The generator intelligently creates **temporary views** for intermediate pipeline steps and defines **final outputs** as durable tables — giving you the flexibility to move between exploration and production without rewriting logic. Whether running in batch or real-time streaming mode, the generated pipelines are built with checkpointing, retries, and resilience out of the box.

**Functions:**

- Compiles `.dsl` DSL into fully executable Go code  
- Builds DAGs (Directed Acyclic Graphs) from ordered step dependencies  
- Supports parameterized tags like:  
  - `@source(type=snowflake|kafka, schema|topic=…, table|format=…)`  
  - `@step(name=ValidateSIN, domain=…)`  
  - `@sink(type=snowflake|kafka, schema|topic=…, table|format=…)`  
- Generates **views** for intermediate pipeline objects  
- Writes final outputs as persistent **database tables**  
- Supports both **batch and real-time streaming pipelines**  
- Includes built-in checkpointing, retries, and fault-tolerant flow control  

With the Pipeline Generator, you go from **business intent to running code** in seconds — no YAML files, no glue code, just pure scenario-driven automation.

---

### 🔄 Feature Orchestration

> Role: Handles pipeline dependency resolution and execution scheduling.

The **Feature Orchestration** engine is the backbone of Runink’s intelligent pipeline execution. It transforms `.dsl` files into **distributed DAGs (Directed Acyclic Graphs)** that honor data dependencies, execution order, and domain boundaries. Each step — tagged in your DSL — is independently executable and can be scheduled as a microservice or serverless job, enabling true modular execution at scale.

Orchestration works seamlessly across **batch** and **streaming** modes, supporting windowed ingestion and real-time event slices with equal efficiency. It leverages Go-native primitives (like `bigmachine`) and Linux cgroups to spin up lightweight, isolated **Runi** compute slices. These slices execute in logical groupings — called **Herds** — using Linux namespaces to separate domain-specific logic, making the entire system cloud-native and data-domain aware by design.

Whether executing in a local cluster or across distributed environments like Kubernetes, Runink’s orchestration ensures resilience through retries, observability, auto-scaling, and minimal resource overhead — giving you a serverless experience with full pipeline traceability and control.

***Functions:***

- Robust streaming and windowed processing, with batch ingestion support  
- Parallel and distributed DAG-based pipeline execution  
- Execution of tagged pipeline steps as microservices or serverless units  
- Smart compute orchestration via **Runi slices** (cgroups)  
- Domain boundary isolation using **Herd namespaces**  
- Health checks, automatic retries, and self-healing execution plans  
- Ingests both bounded and unbounded data with slice-based identity  
- Matches pending pipeline tasks to available capacity.  
- Reads desired tasks from the Barn
- Applies Herd quotas, node labels, and resource requirements  
- Sends placement commands to Runi Agents  

With this architecture, Runink doesn’t just run pipelines — it orchestrates **data responsibility**, **compute efficiency**, and **modular autonomy** at enterprise scale.

---

### 🧭 Data Lineage & Metadata

The **Data Lineage & Metadata** layer in Runink gives you full transparency into how data flows — and why it changes. As each record moves through a pipeline, Runink automatically tracks what transformed it, under what contract, when it happened, and where it ran. This built-in traceability turns every pipeline into an auditable, explorable map of data activity.

Runink attaches rich metadata to every record — including `RunID`, `stage name`, `contract version`, timestamps, and transformation paths — all without requiring extra instrumentation. This metadata powers compliance reports, observability dashboards, and rollback logic. It also supports regulatory requirements like GDPR, HIPAA, and internal DataOps standards.

Visual lineage graphs can be rendered from pipeline metadata to show how datasets were derived, making both debugging and documentation seamless.

***Functions:***

- Automatic lineage tracking across all pipeline stages  
- Visualize transformations via DAGs and lineage graphs  
- Metadata tagging (RunID, stage, contract hash, timestamps)  
- Immutable audit logs for compliance and traceability  
- Hooks into REPL, test output, and contract validation reports  

With lineage and metadata as a first-class feature, Runink ensures your data is **not just accurate**, but **accountable**.

---

### ✅ Data Quality & Validation

> Role: Maintains data integrity and compliance throughout the pipeline.

The **Data Quality & Validation** layer in Runink ensures that your pipelines don’t just run — they run with **integrity, accuracy, and trust**. It allows you to define validations as part of your pipeline logic, using both schema-driven rules and domain-specific assertions. From pre-ingestion *contract validation* checks to mid-pipeline *schema enforcement* validations and final output guarantees, Runink makes quality a first-class citizen.

Non-technical users can define business rules using `.dsl` files — asserting required fields, value ranges, regex patterns, or even cross-field relationships. Invalid records can be logged, routed to a **Dead Letter Queue (DLQ)**, or flagged for review, without halting the entire pipeline.

Built-in quality metrics and thresholds provide insight into pipeline health, and integrate seamlessly with observability tools for alerting and monitoring.

***Functions:***

- Validate records against schema contracts at multiple stages  
- Enforce field-level, type-level, and domain-specific constraints  
- Route invalid records to DLQs with full metadata for debugging  
- Monitor and alert on data quality thresholds  
- Allow business and data teams to codify validation rules in `.dsl` files  

With Runink, data quality becomes proactive, testable, and deeply integrated — not an afterthought.

---

### 🔒 Data Governance & Compliance

> Role: Ensures adherence to regulatory and internal data governance standards.

The **Data Governance & Compliance** layer in Runink ensures that your pipelines don’t just move data — they move it **securely**, **transparently**, and **in compliance** with both regulatory and internal standards. From role-based access control to encryption and masking, Runink provides the tools and policies to safeguard sensitive data while enabling responsible data sharing across domains.

Whether you’re building pipelines for internal analytics or external-facing services, Runink enforces security best practices through clearly defined policies and supports key compliance frameworks like **SOC 2**, **GDPR**, and **ISO 27001**. Each component is designed to be traceable and accountable, with **SPDX-compliant licensing**, **SBOM reporting**, and secure identity handling via **OIDC** and **JWT**.

Data sensitivity is handled directly in the pipeline — with built-in support for **anonymization**, **data masking**, and **field-level encryption** — so governance becomes part of execution, not an afterthought.

***Functions:***

- Exposes gRPC/REST API for UI, CLI, and policy checks  
- Ingests annotations, metadata, and audit logs from slices  
- Enables lineage visualization and quality dashboards  
- SPDX-compliant licensing and SBOM (Software Bill of Materials) exports  
- Integration and compliance with FDC3 and CDM standards
- Support for SOC 2, GDPR, ISO 27001, HIPAA, and more  
- Schema alignment and data interoperability across financial systems
- Schema versioning and drift monitoring specific to industry models
- Applied FinOps on data

Runink brings **data responsibility** and **regulatory readiness** directly into the pipeline, helping your organization scale securely — and confidently.

---


### 🛡️ Security & Enterprise DataOps

> Role: Secure systems during its conception in an automated manner. 

The **Security & Enterprise DataOps** layer in Runink ensures that your pipelines are **secure from day one** — by automating protection across the entire software and data lifecycle. Security isn’t an afterthought; it’s integrated directly into how pipelines are built, tested, and deployed.

Runink automates **vulnerability scanning** and license auditing during CI/CD, using GitHub Actions and GoReleaser to generate secure, signed artifacts with **SBOMs** and **hash verification** to prevent tampering. It supports **environment-specific configurations**, managing secrets and access per environment, whether local, staging, or production.

For runtime security, pipelines are governed by **RBAC** and authenticated via **OIDC/JWT**, while data is protected at multiple layers — with **encryption at rest and in transit**, **field-level masking**, and **anonymization** for sensitive fields.

***Functions:***

- Automated security scanning and vulnerability management during CI/CD  
- Service Mesh integration
- Dynamic configuration management
- GitHub Actions + GoReleaser pipelines with reproducible builds  
- Data encryption at rest (e.g., file systems, S3) and in transit (TLS)  
- Environment-specific pipeline config + secret injection  
- Fine-grained network access control and isolation  
- Role-based access control (RBAC) with secure auth via OIDC/JWT  
- Field-level encryption, masking, and anonymization for PII/PHI  
- Signed hash verification of pipeline artifacts and outputs to prevent tampering  

Runink enables **secure-by-default pipelines** for modern, responsible data platforms — empowering you to scale DataOps without sacrificing trust, traceability, or control.

---

### 🌐 Observability & Monitoring

> Role: Enables comprehensive visibility into pipeline health and performance.

Observability means you can **see what your pipelines are doing, when, and why**. In Runink, every feature scenario and pipeline stage is fully instrumented with **real-time metrics**, **structured logs**, and **traceable telemetry** — giving you deep, actionable visibility into your data workflows.

From **record counts** and **stage durations** to **error rates** and **schema contract versions**, Runink captures detailed metrics tied to each run, stage, and data slice. These are exposed via Prometheus, OpenTelemetry, or your own observability stack, enabling rich dashboards, alerting, and deep debugging capabilities.

Structured logs include tags like `RunID`, `Stage`, and `ContractHash`, while telemetry tracing shows full DAG execution across parallel or distributed nodes. Whether you're tuning performance, chasing bugs, or watching for data anomalies — Runink helps you stay one step ahead.

***Functions:***

- Monitor pipelines in real time with streaming metrics  
- Debug slow or failing stages with stage-level performance data  
- Alert on errors, retries, or threshold breaches  
- Integrate data quality thresholds into observability dashboards  
- Structured logs enriched with run metadata for traceability  
- Visualize end-to-end pipeline execution with OpenTelemetry traces  

With Runink, observability is not an add-on — it's an embedded guarantee that your data platform stays transparent, reliable, and responsive.