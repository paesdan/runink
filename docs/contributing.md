# 🤝 Contributing to Runink

Welcome! First off, thank you for considering contributing to **Runink**. We deeply appreciate your support and effort to improve our project.

This document will guide you through the process of contributing code, filing issues, suggesting features, and participating in the Runink community.

---

## 📜 Code of Conduct

We expect everyone participating to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md) (to be created). Respect and kindness are the foundation.

---

## 🛠️ How to Contribute

### 1. Fork the Repo

Use GitHub's "Fork" button to create a personal copy of the repository.

### 2. Clone Your Fork

```bash
git clone https://github.com/your-username/runink.git
cd runink
```

### 3. Create a New Branch

Use a clear branch naming convention:

```bash
git checkout -b feature/short-description
# or
git checkout -b fix/bug-description
```

### 4. Make Your Changes

Follow our coding guidelines:
- Write idiomatic Go (gofmt, golint).
- Keep PRs small and focused.
- Update or add tests for your changes.
- Update documentation (`docs/`) if applicable.

### 5. Test Before You Push

Run all tests:

```bash
make lint
make test
```

### 6. Push and Open a Pull Request

Push to your fork and open a Pull Request against the `main` branch.

```bash
git push origin feature/short-description
```

On GitHub, create a new Pull Request and fill in the template (title, description, related issues).

---

## 📋 Development Guidelines

- **CLI Commands:** Place new commands inside their respective domain folder (`barnctl`, `buildctl`, `herdctl`, `runictl`).
- **Testing:** Add unit tests for CLI commands, helpers, validators.
- **Logging:** Use structured logging where needed.
- **Security:** Always consider security (no plaintext secrets, minimal privilege).
- **Performance:** Avoid premature optimization, but don't introduce obvious inefficiencies.

---

## 🔍 Reporting Bugs

- Search existing issues first.
- File a [new issue](https://github.com/your-username/runink/issues/new) with clear reproduction steps.
- Provide logs, stack traces, and your environment (OS, Go version).

> If you discover a security vulnerability, **please do not open a public issue.**  
Instead, email us at [paes@dashie.ink](mailto:paes@dashie.ink).


---

## 🚀 Suggesting Features

- Open an Issue labeled `enhancement`.
- Explain your use case and how it aligns with Runink's vision.

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
## 📅 Regular Updates

We sync main with active development regularly. Expect fast iteration and reviews.

If you have any questions, feel free to open an issue or discussion!

Thanks for being part of the **Runink** Herd and for helping us build the future of safe, expressive, and reliable data pipelines.🐑  

We can’t wait to see what you contribute! 🙌

— The Runink Team