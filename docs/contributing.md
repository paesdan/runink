# 🤝 Contributing to Runink

Welcome, and thank you for your interest in contributing to **Runink** — the data platform that makes pipelines safe, expressive, and composable. Whether you're fixing bugs, improving documentation, writing tests, or building new features, your help makes Runink better for everyone.

We’re excited to have you here! 🚀

---

## 🧭 Quick Start

1. **Fork the repository**
2. **Clone your fork**

```bash
git clone https://github.com/your-username/runink.git
cd runink
```

3. **Install Go (1.21+)** and required tools
4. **Set up your environment**

```bash
make install-tools
make setup
```

5. **Start hacking** 🚧  
Use the CLI or REPL to run tests and iterate locally:

```bash
runink init test-pipeline
runink compile
runink test
```

---

## 📂 Project Structure

Here's where to plug in:

| Folder | Purpose |
|--------|---------|
| `internal/` | Core Go logic — CLI, pipeline engine, agents |
| `features/` | DSL scenarios (`.feature` files) |
| `contracts/` | Data contracts & golden files |
| `docs/` | Documentation and guides |
| `deploy/` | Generated orchestration artifacts |
| `.github/` | CI configs, issue templates, actions |

---

## 🛠️ How to Contribute

### 🐞 Report Bugs

- Use [GitHub Issues](https://github.com/runink/runink/issues)
- Please include:
  - Steps to reproduce
  - Expected vs. actual behavior
  - Runink version (`runink version`)
  - OS / system info

### 🌟 Suggest Features

We love ideas! Start a discussion or open a feature request issue. Try to include:
- A clear problem statement
- Example DSL or contract (if relevant)
- Why this helps the community

### 🧑‍💻 Submit Code

#### 1. Create a new branch

```bash
git checkout -b feat/my-awesome-change
```

#### 2. Make your changes  
Follow existing conventions. Run tests (`make test`) before committing.

#### 3. Format and lint

```bash
make fmt
make lint
```

#### 4. Commit and push

```bash
git commit -m "feat(dsl): add support for new step type"
git push origin feat/my-awesome-change
```

#### 5. Open a Pull Request

- Keep PRs focused and well-scoped
- Include tests and docs if relevant
- Use conventional commit messages (`feat:`, `fix:`, `chore:`, etc.)

---

## 🧪 Testing & Validation

- Use golden files (`contracts/testdata/`) for regression testing
- Run `runink test` to validate scenarios
- Use the REPL for quick exploration and debugging
- Ensure your PR passes CI

---

## 🔐 Security

If you discover a security vulnerability, **please do not open a public issue.**  
Instead, email us at [paes@dashie.ink](mailto:paes@dashie.ink).

---

## ❤️ Code of Conduct

We’re a community of data builders. We expect contributors to be respectful, inclusive, and constructive.

Please read our [Code of Conduct](./CODE_OF_CONDUCT.md) before contributing.

---

## 🧵 Join the Community

- GitHub Discussions (coming soon)
- Discord server (invite coming soon)
- Follow our roadmap in [`docs/roadmap.md`](./docs/roadmap.md)

---

Thanks for helping us build the future of safe, expressive, and reliable data pipelines. We can’t wait to see what you contribute! 🙌

— The Runink Team