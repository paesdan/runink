## 📦 CHANGELOG

All notable changes to **Runink** will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

> Ongoing work and unreleased features

### 🚀 Added

- _Initial entries go here..._

### 🐛 Fixed

- _Bugfixes not yet released..._

### ♻️ Changed

- _Refactors, renames, behavior tweaks..._

### 🔐 Security

- _Security-related updates pending release..._

---

## [v0.6.0] - 2025-04-22

### 🚀 Added

- gRPC support for real-time feature orchestration
- Barn-backed Secrets Manager with scoped delivery
- SQL mode in REPL with `.dsl` export preview
- Token-scoped RBAC policies with CLI support
- New `runink publish` command for pushing lineage + SBOM

### 🐛 Fixed

- REPL crash on malformed JSON inputs
- Schema drift diff falsely triggering on optional fields

### ♻️ Changed

- Updated logging to structured format across all CLI commands
- Refactored `runi test` to support golden file overrides

---

## [v0.5.0] - 2025-03-10

### 🚀 Added

- Herd Namespace isolation for multi-tenant workloads
- Feature DSL step handlers now support parameter validation
- Dead Letter Queue (DLQ) support with retry metadata

### 🐛 Fixed

- Lineage metadata missing contract hash in REPL output
- `runi compile` now correctly resolves nested includes

### 🔐 Security

- Signed artifacts via GitHub Actions with SBOM & hash outputs

---

## [v0.4.0] - 2025-01-28

### 🚀 Added

- CLI-based project scaffolding (`runink init`)
- Contract generator from Go structs (`runi contract gen`)
- Golden file testing engine
- Feature DSL execution for batch pipelines
- Interactive REPL for pipeline debugging

---

## [v0.3.0] - 2024-12-15

Initial open-source public release 🎉

---

## 🔖 Tags

- **🚀 Added** – New features or modules
- **🐛 Fixed** – Bugfixes
- **♻️ Changed** – Refactors, non-breaking changes
- **🛠 Deprecated** – Features slated for removal
- **💥 Removed** – Breaking changes or removals
- **🔐 Security** – Fixes related to vulnerabilities

---

> For a full list of commits, see the [GitHub Release History](https://github.com/runink/runink/releases)
