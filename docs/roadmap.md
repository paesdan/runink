# 🗺️ Runink Roadmap

Welcome to the official **Runink Roadmap** — our evolving guide to what we're building, where we’re headed, and how you can get involved.

Runink is built on the belief that modern data platforms should be **safe by default**, **composable by design**, and **collaborative at scale**. This roadmap reflects our commitment to transparency, community-driven development, and rapid iteration.

---

## 🧩 Roadmap Themes

| Theme | Description |
|-------|-------------|
| **Composable Pipelines** | Make it easy to build, reuse, and test pipeline steps across teams and domains. |
| **Secure & Compliant by Default** | Tighten RBAC, data contracts, and observability for enterprise-grade governance. |
| **DevX & Developer Productivity** | Empower devs with a powerful CLI, REPL, codegen, and rapid iteration loops. |
| **Streaming-First DataOps** | Advance real-time use cases with backpressure-safe, contract-aware streaming. |
| **Interoperability & Ecosystem** | Play well with FDC3, CDM, OpenLineage, Kafka, Snowflake, and more. |

---

## 🧭 Current Focus (Q2 2025)

These items are in active development or early testing:

- [ ] **Herd Namespace Isolation** (multi-tenant namespace support)
- [ ] **Golden Test Rewrites** for easier review and diffing
- [ ] **CLI REPL SQL Mode** with DataFrame-to-Feature export
- [ ] **RBAC & Token Scoping Enhancements**
- [ ] **Raft-backed Barn & Secrets Manager Integration**
- [ ] **gRPC Streaming Orchestration**
- [ ] **Pipeline Preview Mode (dry-run with metadata only)**
- [ ] **Lineage UI + CLI support**
- [ ] **Remote Artifact Signing & SBOM generation (SLSA-style)**

---

## 🔜 Near-Term (Q3 2025)

Planned next based on user feedback and enterprise needs:

- [ ] **Live Feature File Linter & Formatter**
- [ ] **REPL Session Recorder** (record → replay feature building)
- [ ] **Multi-Herd Scheduling & Cost Reporting**
- [ ] **Secrets Rotation + External Vault Integration**
- [ ] **Contract Diff Web Viewer**
- [ ] **Push-to-Registry UX from CLI**
- [ ] **DLQ Visualization + Retry Tools**
- [ ] **Plugin Marketplace (source/sink/step handlers)**

---

## 🌅 Long-Term Vision (Late 2025+)

Our long-range goals to shape Runink into the **standard platform for responsible data pipelines**:

- ⚙️ **Full No-YAML Orchestration (Declarative-Only Pipelines)**
- 🧠 **AI Copilot for Contract & Scenario Generation**
- 🌐 **Cross-Org Data Mesh Support via Herd Federation**
- 📡 **Runink Cloud** (fully managed, secure, multi-tenant SaaS)
- 🔒 **Zero-Trust Data Contracts (ZK + Provenance)**

---

## 🧠 Ideas We're Exploring

These are in research/design phases — feedback welcome!

- ✨ Feature DSL Step Suggestions in CLI
- 🔀 Schema Merge Conflict Resolution UX
- 📥 Native ingestion support for S3/Parquet/Arrow
- 🔎 Full integration with OpenLineage + dbt Core
- 🧾 GitHub Copilot integration for contract authoring

---

## 🙋‍♀️ Contribute to the Roadmap

We prioritize what the community and users need most. If there’s a feature you’d love to see:

- Open an issue using the `Feature Request` template
- Upvote existing roadmap items via 👍 reactions
- Join upcoming roadmap discussions (Discord coming soon!)
- PRs welcome for anything marked as `help-wanted`

---

## 🔄 Release Cadence

We aim for:

- **Minor releases** every 4–6 weeks (feature drops, improvements)
- **Patch releases** as needed (hotfixes, regressions)
- **Major milestones** every ~6 months with community showcases

Track progress in [`CHANGELOG.md`](./CHANGELOG.md)

---

Thanks for being part of the journey — we’re building Runink **with** you, not just **for** you. Let’s define the future of safe, modular, and explainable data platforms together.

— Team Runink 🐑