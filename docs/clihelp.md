
# 🆘 Runink CLI: Help Template

This is a developer-friendly help template for implementing consistent `runi <command> --help` outputs.

---

## 🧱 Format: Basic Help Command

```
runi <command> [subcommand] [flags]

Usage:
  runi <command> [options]

Options:
  -h, --help       Show this help message and exit
  -v, --verbose    Show detailed logs and diagnostics
```

---

## Example: `runi init --help`

```
Initialize a new Runink project.

Usage:
  runi init [project-name]

Flags:
  -h, --help     Show help for init
```

---

## 🔄 Example: `runi compile --help`

```
runi compile --scenario <file.dsl>

Description:
  Compile a `.dsl` scenario and its contract into an executable Go DAG.
  Generates a Go file under `rendered/` based on contract-linked step tags.

Usage:
  runi compile --scenario features/trade_cdm.dsl

Flags:
  --scenario    Path to a DSL scenario file
  --out         Optional: custom output path for DAG (default: rendered/<name>.go)
  --dry-run     Only validate scenario and contract, do not write DAG
  --verbose     Show full DAG step resolution logs
```

---

## 🧪 Example: `runi test --help`

```
runi test --scenario <file.dsl>

Description:
  Execute a feature scenario with golden test inputs and compare output.
  Supports diff mode and golden update flows.

Usage:
  runi test --scenario features/onboard.dsl

Flags:
  --scenario       DSL file to test
  --golden         Optional: override path to golden test folder
  --update         Automatically update golden output on success
  --only <step>    Run test up to a specific pipeline step
```

---

## 🔐 Example: `runi contract gen --help`

```
runi contract gen --struct <package.Type> --out <file>

Description:
  Generate a JSON contract definition from a Go struct. Includes schema, access tags, and validation metadata.

Usage:
  runi contract gen --struct contracts.Customer --out contracts/customer.json

Flags:
  --struct     Fully qualified Go type (e.g. contracts.Customer)
  --out        Output contract file path
  --flatten    Inline nested types into flat fields
  --herd       Optional: attach to specific herd (e.g. finance)
```

---

---

## `runi contract diff --help`

```
Diff two versions of a contract and show schema drift.

Usage:
  runi contract diff --old v1.json --new v2.json
```

---

## `runi run --help`

```
Run a compiled pipeline with data inputs.

Usage:
  runi run --scenario <file.dsl> [--contract file] [--herd name]

Flags:
  --scenario     Scenario to execute
  --contract     Optional explicit contract
  --herd         Herd to run pipeline in
  --dry-run      Preview DAG resolution only
```

---

## `runi lineage --help`

```
Show lineage metadata for a run.

Usage:
  runi lineage --run-id <id>

Flags:
  --run-id       Unique run identifier
  --output       Format (json|csv|graph)
```

---

## `runi publish --help`

```
Publish contracts, lineage, and tags to metadata registry.

Usage:
  runi publish --herd <name> [--scenario file]
```

---

## `runi repl --help`

```
Start interactive REPL for querying test inputs or contract data.

Usage:
  runi repl --scenario <path>
```

---

## 🤖 Example: `runi deploy --help`

```
runi deploy --target <target>

Description:
  Deploy Runi workers and slices to a remote orchestration cluster.

Usage:
  runi deploy --target k8s

Flags:
  --target      Target platform (k8s, bigmachine)
  --herd        Herd (namespace) to deploy into
  --dry-run     Simulate deployment without applying
  --confirm     Require manual confirmation for remote changes
```

---

## `runi schedule --help`

```
Schedule a pipeline scenario for recurring execution.

Usage:
  runi schedule --scenario <file> --cron "0 6 * * *"

Flags:
  --scenario     DSL file
  --cron         Cron-style expression
```

---

## `runi audit --help`

```
Show schema contract change history and approvals.

Usage:
  runi audit --contract <file>
```
---

---

## `runi restart --help`

```
Restart a failed or incomplete pipeline run from its last checkpoint.

Usage:
  runi restart --run-id <uuid>

Flags:
  --run-id     Run ID to restart from
  --force      Ignore checkpoint and rerun from start
```

---

## `runi resume --help`

```
Resume an interrupted or paused pipeline.

Usage:
  runi resume --run-id <uuid>
```

---

## `runi checkpoint --help`

```
Create a DAG state checkpoint for partial run recovery.

Usage:
  runi checkpoint --scenario <file>
```

---

## `runi comment --help`

```
Leave inline comments for contracts or fields (used in review tools).

Usage:
  runi comment --contract <file> --field <path> --note <text>
```

---

## `runi request-approval --help`

```
Submit a contract for governance approval and audit.

Usage:
  runi request-approval --contract <file>
```

---

## `runi feedback --help`

```
Attach feedback note to a scenario feature for peer review.

Usage:
  runi feedback --scenario <file> --note <text>
```

---

## `runi redact --help`

```
Automatically redact fields marked pii:"true" in a contract schema.

Usage:
  runi redact --contract <file> --out <file>
```

---

## `runi escrow --help`

```
Encrypt and store output data for delayed release or approval.

Usage:
  runi escrow --run-id <uuid> --out <vault.json>
```

---

## `runi anonymize --help`

```
Create a non-sensitive version of input using faker + tags.

Usage:
  runi anonymize --input <file> --contract <file> --out <file>
```

---

## `runi trigger --help`

```
Define an event trigger for this scenario.

Usage:
  runi trigger --scenario <file> --on webhook|s3|pubsub
```

---

## `runi listen --help`

```
Start a listener to monitor incoming event and dispatch pipeline.

Usage:
  runi listen --event <type>
```

---

## `runi subscribe --help`

```
Subscribe to a streaming topic or channel with offset tracking.

Usage:
  runi subscribe --stream <topic> --window 5m
```

---

## `runi freeze --help`

```
Freeze contract + scenario versions with hashes for snapshot validation.

Usage:
  runi freeze --scenario <file>
```

---

## `runi archive --help`

```
Archive old versions of scenarios and their runs by herd.

Usage:
  runi archive --herd <name> --keep 3
```

---

## `runi retire --help`

```
Retire a contract so it cannot be used in future scenarios.

Usage:
  runi retire --contract <file>
```

---

## `runi lineage graph --help`

```
Export full DAG and contract lineage as GraphViz dot file.

Usage:
  runi lineage graph --out lineage.dot
```

---

## `runi knowledge export --help`

```
Export pipeline metadata using RDF serialization (Turtle/N-Triples).

Usage:
  runi knowledge export --format turtle
```

---

## `runi query lineage --help`

```
Query lineage metadata using SQL-like syntax.

Usage:
  runi query lineage --sql "SELECT * WHERE pii = true"
```

---

## `runi openai audit --help`

```
Use an LLM to summarize contract diffs or suggest field comments.

Usage:
  runi openai audit --contract <file>
```

---

## `runi sandbox --help`

```
Execute scenario in a secure ephemeral environment.

Usage:
  runi sandbox --scenario <file>
```

---

## `runi simulate --help`

```
Replay input data as a stream window to test stateful logic.

Usage:
  runi simulate --input <file> --window 5m
```

---

## `runi mint-token --help`

```
Generate a short-lived JWT for scoped access by herd or scenario.

Usage:
  runi mint-token --herd finance --role analyst --ttl 5m
```

---

## 🧠 Best Practices

- ✅ Describe what the command **does**, not how it's implemented
- ✅ Include at least 1 usage example
- ✅ Use consistent flags: `--scenario`, `--contract`, `--out`, `--herd`
- ✅ Provide guidance for `--dry-run`, `--verbose`, `--help`
- ✅ Include multi-step examples if command touches multiple files

---

