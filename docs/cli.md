# 🧰 Runink CLI Reference (`runi`)

The `runi` CLI is the command-line interface to everything in the **Runink** data platform. It’s your developer-first companion for defining, testing, running, securing, and publishing data pipelines — all from declarative `.dsl` files and Go-native contracts.

This reference describes all available commands, grouped by capability.

---

## 🧱 Project & Pipeline Lifecycle

| Command | Description |
|--------|-------------|
| `runi herd init [project-name]` | Scaffold a new workspace with starter contracts, features, CI config |
| `runi compile --scenario <file>` | Generate Go pipeline code from `.dsl` files |
| `runi run --scenario <file> --contract <contract.json>` | Run pipelines locally or remotely |
| `runi watch --scenario <file>` | Auto-compile & re-run scenario on save |

---

## 📁 Schema & Contract Management

| Command | Description |
|--------|-------------|
| `runi contract gen --struct <pkg.Type>` | Generate a contract from Go struct |
| `runi contract diff --old v1.json --new v2.json` | Show schema drift between versions |
| `runi contract rollback` | Revert to previous contract version |
| `runi contract history --name <contract>` | Show all versions and changelog entries |
| `runi contract validate --file <file>` | Validate a file against a contract |
| `runi contract catalog` | Create an index of all contracts in the repo |
| `runi contract hash` | Generate contract hash for versioning |

---

## 🧪 Testing & Data Validation

| Command | Description |
|--------|-------------|
| `runi test --scenario <file>` | Run tests using golden files |
| `runi synth --contract <contract.json>` | Generate synthetic golden test data |
| `runi verify-contract --contract <file>` | Validate pipeline against contract |

---

## 🔐 Security, Publishing & Compliance

| Command | Description |
|--------|-------------|
| `runi secure [--sbom|--sign|--scan]` | Run security audits and generate SBOM |
| `runi publish` | Push metadata, lineage, and contracts to registry |
| `runi sbom export [--format spdx]` | Export SPDX-compliant software bill of materials |
| `runi changelog gen` | Generate changelogs from contract/feature diffs |

---

## 🔍 Observability & Lineage

| Command | Description |
|--------|-------------|
| `runi lineage --run-id <uuid>` | Show DAG lineage for a run |
| `runi lineage track --source A --sink B` | Manually link lineage metadata |
| `runi lineage graph --out file.dot` | Export lineage graph in DOT format |
| `runi metadata get --key <name>` | Retrieve stored metadata for a step |
| `runi metadata annotate --key <name>` | Attach annotation to pipeline metadata |
| `runi logs --run-id <uuid>` | View logs for a specific run |
| `runi status --run-id <uuid>` | Check status of a pipeline execution |

---

## 🤖 Distributed Execution (Remote)

| Command | Description |
|--------|-------------|
| `runi deploy --target <k8s|bigmachine>` | Deploy Runi workers to a remote cluster |
| `runi start --slice <file> --herd <namespace>` | Start execution of a scenario remotely |
| `runi kill --run-id <uuid>` | Terminate running scenario |

---

## 💪 Control Plane & Agents

| Command | Description |
|--------|-------------|
| `runi herdctl create` | Create a new Herd (namespace + quotas + policies) |
| `runi herdctl delete` | Delete a Herd |
| `runi herdctl update` | Update Herd quotas, RBAC, metadata |
| `runi herdctl list` | List all Herds and resource states |
| `runi herdctl quota set <herd>` | Update CPU/mem quotas live |
| `runi herdctl lineage <herd>` | View lineage graphs scoped to a Herd |
| `runi agentctl list` | List active Runi agents, resource usage, labels |
| `runi agentctl status <agent>` | Detailed agent status (health, registered slices, metrics) |
| `runi agentctl drain <agent>` | Mark agent as unschedulable (cordon) |
| `runi agentctl register` | Manually register agent (optional bootstrap) |
| `runi agentctl cordon <agent>` | Prevent slice scheduling |

---

## 🌐 Worker Slice Management

| Command | Description |
|--------|-------------|
| `runi slicectl list --herd <id>` | List all active slices for a Herd |
| `runi slicectl logs <slice-id>` | Fetch logs for a given slice |
| `runi slicectl cancel <slice-id>` | Cancel a running slice gracefully |
| `runi slicectl metrics <slice-id>` | Show real-time metrics for a slice |
| `runi slicectl promote <slice-id>` | Checkpoint a slice mid-run |

---

## 🔀 Introspection & Visualization

| Command | Description |
|--------|-------------|
| `runi explain --scenario <file>` | Describe DAG and step resolution logic |
| `runi graphviz --scenario <file>` | Render DAG as a .png, .svg, or .dot |
| `runi diff --feature old.dsl --feature new.dsl` | Compare feature files and show logic drift |

---

## 🧪 REPL & Exploratory Commands

| Command | Description |
|--------|-------------|
| `runi repl` | Launch interactive DataFrame, SQL, JSON REPL |
| `runi json explore -f file.json -q '.email'` | Run jq-style query on JSON |
| `runi query -e "SELECT * FROM dataset"` | Run SQL-like query on scenario input |

---

## 🛠️ Dev Tools & Generators

| Command | Description |
|--------|-------------|
| `runi gen --dsl input.json` | Generate feature from sample input |
| `runi contract from-feature <file>` | Extract contract from `.dsl` spec |
| `runi schema hash` | Generate contract fingerprint |
| `runi bump` | Auto-increment contract version with changelog |

---

## 🧹 Plugins & Extensions

| Command | Description |
|--------|-------------|
| `runi plugin install <url>` | Install external plugin |
| `runi plugin list` | List installed extensions |
| `runi plugin run <name>` | Execute a plugin subcommand |

---

## 📦 Packaging & CI/CD

| Command | Description |
|--------|-------------|
| `runi build` | Compile pipeline bundle for remote use |
| `runi pack` | Zip workspace for deployment/distribution |
| `runi upgrade` | Self-update the CLI and plugins |
| `runi doctor` | Diagnose CLI and project setup |

---

## 📅 Runtime Lifecycle

| Command | Description |
|--------|-------------|
| `runi restart --run-id <uuid>` | Restart a pipeline from last successful checkpoint |
| `runi resume --run-id <uuid>` | Resume paused pipeline without reprocessing |
| `runi checkpoint --scenario <file>` | Create a persistent step-based checkpoint marker |

---

## 💬 Collaboration & Governance

| Command | Description |
|--------|-------------|
| `runi comment --contract <file>` | Leave inline comments for review (contract-level QA) |
| `runi request-approval --contract <file>` | Submit contract for governance approval |
| `runi feedback --scenario <file>` | Attach review notes to a scenario |

---

## 🛡️ Privacy, Redaction & Data Escrow

| Command | Description |
|--------|-------------|
| `runi redact --contract <file>` | Automatically redact PII based on tags |
| `runi escrow --run-id <uuid>` | Encrypt pipeline outputs for future unsealing |
| `runi anonymize --input <file>` | Generate synthetic version of a sensitive input file |

---

## 🗓 Event-Based Execution

| Command | Description |
|--------|-------------|
| `runi trigger --on <webhook|s3|pubsub>` | Set up trigger-based pipeline starts |
| `runi listen --event <type>` | Listen for external event to start scenario |
| `runi subscribe --stream <source>` | Subscribe to stream source with offset recovery |

---

## 🔄 Pipeline & Contract Lifecycle

| Command | Description |
|--------|-------------|
| `runi freeze --scenario <file>` | Lock DAG hash and contract state as immutable |
| `runi archive --herd <name> --keep <N>` | Archive old scenarios/runs beyond retention policy |
| `runi retire --contract <file>` | Deprecate contract from active use |

---

## 🧬 Metadata Graph & Semantic Search

| Command | Description |
|--------|-------------|
| `runi knowledge export --format turtle` | Export contract and DAG metadata as RDF |
| `runi query lineage` | Run SQL-style queries across lineage metadata |

---

## 🧪 Experimental / LLM-integrated

| Command | Description |
|--------|-------------|
| `runi openai audit --contract <file>` | Use LLM to summarize or describe schema changes |
| `runi sandbox --scenario <file>` | Run pipeline in ephemeral namespace sandbox |
| `runi simulate --input <file> --window 5m` | Replay stream input over a time window |
| `runi mint-token --role <name>` | Mint temporary JWT token scoped to herd/contract |

---

> 💬 Use `runi <command> --help` for flags, options, and examples.
>  
> Runink's CLI gives you a full stack data engine in your terminal — from contracts to clusters, from `.dsl` to full observability.

