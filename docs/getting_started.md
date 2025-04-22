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
- Example `.feature` files
- Golden file tests
- Dockerfile and CI/CD configs

---

## 📋 **3. Explore the Project Structure**

Your project includes:

```
myproject/
├── wranglers/
├── contracts/
├── contracts/golden/
├── features
├── deploy
├── config
├── docs/
└── .github/workflows/
```

📁 Explore the monorepo to find:

- `wranglers/` — ? (ETL Code to transform the data.)
- `features/` — Gherkin scenarios definitions for each feature use case. (`.feature` files).
- `config/` — ? (Maybe extra configs such connection and rbac rules).
- `deploy/` — ? (Maybe the generated code to be deployed. With all commands and its DAG orchestestration using go or github?).
- `contracts/` — ? (Contracts definitions as go struts)
- `contracts/golden` — Golden test files for regression testing. Holds golden test examples, snapshots, synthetic data.
- `docs/` — markdown docs, examples, use cases, and playbooks.

---

## ⚙️ **4. Compile and Run Pipelines**

Compile your first pipeline:

```bash
runink compile --scenario features/example.feature --out pipeline/example.go
```

Execute a scenario:

```bash
runink run --scenario features/example.feature
```

---

## ✅ **5. Test Your Pipelines**

Use built-in testing and golden files to ensure correctness:

```bash
runink test
```

If the pipeline logic changes and the test is intentionally updated, regenerate golden files:

```bash
runink test --update
```

---

## 🔍 **6. Interactive REPL**

Interactively explore data and debug transformations:

```bash
runink repl
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