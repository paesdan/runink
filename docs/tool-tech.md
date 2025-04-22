# Tools & Technologies Used in Runink

This page lists the major tools, libraries, and technologies that power the Runink platform. Understanding these will help you contribute, debug, or extend Runink more effectively.

---

### ⚙️ Core Language

*Go (Golang)*

The foundation of Runink — fast, statically typed, and great for concurrency.
Used for CLI, compilation, pipeline execution, and all orchestration.
Emphasizes simplicity, performance, and safe parallelism with goroutines and channels.

---

### 🛠 CLI and Code Generation

*Cobra*

Library for building powerful CLI applications in Go.
Handles command definitions, flags, help messages, and auto-completion.
Templates are used for human-readable scaffolding.

---

### 🔌 Protocols and APIs

*gRPC*

Efficient binary protocol for remote execution of pipeline stages.
Enables distributed pipelines and microservice communication.

*Protobuf (Protocol Buffers)*

Used to define strongly typed data structures across services.
Great for defining contracts in a language-agnostic way.

---

### 📈 Observability Stack

*Prometheus*

Collects and stores pipeline metrics (record counts, error rates, etc).
Easy integration with Grafana dashboards.

*Grafana*

Visualizes Prometheus metrics with dashboards and alerts.
Provides real-time insights into pipeline performance.

*OpenTelemetry*

Traces events and spans across pipeline stages.
Integrates with Jaeger or Grafana Tempo for distributed tracing.

*Zerolog / Zap*

Structured logging libraries for Go.
Logs include run_id, stage, contract_hash, and more.

---

### 🧱 Infrastructure & Build Tools

*Bigmachine*

Go library for building distributed, auto-scaling workers.
Supports running stateless compute tasks across many nodes.

*GoReleaser*

Automates versioned releases, cross-compilation, and SBOM generation.
Used to publish CLI binaries, plugins, and release metadata.

*Podman, Skopeo and Buildah*

Containerization of the CLI and distributed services.
Enables repeatable builds and deployments.

*Kubernetes*

Optional orchestration layer for scheduling pipelines at scale.

---

### 🔒 Security & Compliance

*SPDX License Scanning*

Ensures licensing metadata is maintained across dependencies.

*SBOM (Software Bill of Materials)*

Ensures visibility of dependencies and compliance with open-source best practices.

*gosec* 

Static analyzer that detects security vulnerabilities in Go code.

*OIDC / JWT* 

Used for identity-based authentication and access control (RBAC).

---

### ✨ Development & Community

*GitHub Actions*

CI/CD automation for linting, testing, building, and releasing Runink.

*GitHub Discussions*

Community questions, ideas, and roadmap visibility.