# 🧰 Runink CLI Reference (`runi`)

The `runi` CLI is the command-line interface to everything in the **Runink** data platform. It’s your developer-first companion for defining, testing, running, securing, and publishing data pipelines — all from declarative `.dsl` files and Go-native contracts.

This reference describes all available commands, grouped by capability.

---

## 🧱 Project & Pipeline Lifecycle

| Command | Description |
|--------|-------------|
| `runi init [project-name]` | Scaffold a new workspace with starter contracts, features, CI config |
| `runi compile --scenario <file>` | Generate Go pipeline code from `.dsl` files |
| `runi run --scenario <file> --contract <contract.json>` | Run pipelines locally or remotely |
| `runi watch --scenario <file>` | Auto-compile & re-run scenario on save |

---

## 📑 Schema & Contract Management

| Command | Description |
|--------|-------------|
| `runi contract gen --struct <pkg.Type>` | Generate a contract from Go struct |
| `runi contract diff --old v1.json --new v2.json` | Show schema drift between versions |
| `runi contract rollback` | Revert to previous contract version |
| `runi contract history --name <contract>` | Show all versions and changelog entries |

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

## 🔄 Introspection & Visualization

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

## 🧰 Dev Tools & Generators

| Command | Description |
|--------|-------------|
| `runi gen --dsl input.json` | Generate feature from sample input |
| `runi contract from-feature <file>` | Extract contract from `.dsl` spec |
| `runi schema hash` | Generate contract fingerprint |
| `runi bump` | Auto-increment contract version with changelog |

---

## 🧩 Plugins & Extensions

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

## 💡 Example Flow

```bash
runi init finance-platform
cd finance-platform

runi contract gen --struct contracts.Customer --out contracts/customer.json
runi synth --contract contracts/customer.json --scenario features/onboard.dsl

runi compile --scenario features/onboard.dsl
runi test --scenario features/onboard.dsl

runi run --scenario features/onboard.dsl --contract contracts/customer.json
runi lineage --run-id 5ab2...
runi publish
```

---

> 💬 Use `runi <command> --help` for flags, options, and examples.
>  
> Runink's CLI gives you a full stack data engine in your terminal — from contracts to clusters, from `.dsl` to full observability.