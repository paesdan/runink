# Getting Started with Runink

Welcome to **Runink**! This quick-start guide will help you get up and running with Runink to effortlessly build, test, and run data pipelines.

---

## 🚀 **1. Installation**

Make sure you have Go installed (v1.20 or later). Then install Runink:

```bash
go install github.com/runink/runink@latest
```

Ensure `$GOPATH/bin` is in your `$PATH`.

---

## 🛠 **2. Initialize Your Project**

Create a new Runink project in seconds:

```bash
runink init myproject
cd myproject
```

This command generates:
- Initial Go module
- Sample contracts
- Example `.dsl` files
- Golden file tests
- Dockerfile and CI/CD configs

---

## 📋 **3. Explore the Project Structure**

Your project includes:

```
myproject/
├── bin/                  -> CLI
├── api/                  -> gRPC API layer (authn/z, herd context)
├── contracts/            -> Schema contracts and transformation logic on go struts. 
├── features/             -> Scenarios definitions for each feature from the `.dsl` files.
├── golden/     -> Golden files used on regression testing with examples and synthetic data.
├── dags/                 -> Generated DAG code from the contracts and features to be executed by runi.
├── herd/                 -> Domain Service Control Policies Herd context isolation
├── docs/                 -> Markdown docs, examples, use cases, and playbooks.
├── config/runi/          -> Agent: runner, cgroup manager, metrics, logs
├── config/raftstore/     -> State sync using Raft
├── config/scheduler/     -> DAG-aware scheduler logic
├── config/parser/        -> DSL parser, scenario runner
├── config/observability/ -> Tracing, logging, Prometheus exporters
└── .github/workflows/
```

---

## ⚙️ **4. Compile and Run Pipelines**

Compile your first pipeline:

```bash
runink compile --scenario features/example.dsl --out pipeline/example.go --herd my-data-herd
```

Execute a scenario:

```bash
runink run --scenario features/example.dsl --herd my-data-herd
```

---

## ✅ **5. Test Your Pipelines**

Use built-in testing and golden files to ensure correctness:

```bash
runink test --scenario features/example.dsl --herd my-data-herd
```

If the pipeline logic changes and the test is intentionally updated, regenerate golden files:

```bash
runink test --scenario features/example.dsl --update --herd my-data-herd
```

---

## 🔍 **6. Interactive REPL**

Interactively explore data and debug transformations:

```bash
runink repl --scenario features/example.dsl --herd my-data-herd
```

Example REPL commands:

```bash
load csv://data/input.csv
apply MyTransform
show
```

---

## 📚 **7. Next Steps**

- Explore advanced [feature DSL syntax](./feature-dsl.md)
- Learn about [data lineage and metadata tracking](./data-lineage.md)
- Understand [schema and contract management](./schema-contracts.md)

---

## 🚧 **Support & Community**

Need help or have suggestions?

- Open an [issue on GitHub](https://github.com/runink/runink/issues)
- Join our community discussions and get involved!

Let's start building amazing data pipelines with Runink!