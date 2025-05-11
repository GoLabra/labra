# Contributing to LabraGo

Thanks for your interest in contributing! We're building a clean, developer-first headless CMS in Go and React — and we'd love your help.

## Getting Started

To contribute:

1. **Fork this repository**  
   Click the "Fork" button at the top right of the GitHub page.

2. **Clone your fork**
   ```bash
   git clone https://github.com/your-username/labra.git
   cd labrago
   ```

3. **Create a new branch**
   ```bash
   git checkout -b feature/my-feature
   ```

4. **Make your changes**  
   Follow the coding standards and keep commits focused.

5. **Run tests and linters**
    - For Go:
      ```bash
      go test ./...
      ```
    - For React:
      ```bash
      yarn test
      ```
6.**Commit with a clear message**
   ```bash
   git commit -m "feat: add X feature to Y module"
   ```

7. **Push your branch**
   ```bash
   git push origin feature/my-feature
   ```

8. **Open a Pull Request**  
   Go to the GitHub repo and create a PR. Fill out the PR template, link any issues, and describe what you did and why.

---

## Code Style & Tools

- **Go code**
    - Use `gofmt` or `goimports`

- **React code**
    - Lint
    - Functional components only, please

Keep functions small and focused. Avoid magic or clever abstractions that reduce clarity.

---

## Pull Request Guidelines

- Stick to one feature or fix per PR.
- Update or add tests when applicable.
- If your change affects the UI, include screenshots or a short description of the result.
- Add helpful commit messages (e.g. `fix: prevent crash on empty config`).

---

## Reporting Issues

When submitting an issue:
- Check if it already exists.
- Provide a clear, reproducible example.
- Include logs, screenshots, or environment details if possible.

Use the following labels when opening issues (if applicable):
- `bug`, `enhancement`, `question`, `help wanted`

---

## Join the Discussion

Have ideas, feedback, or want to get involved in shaping the roadmap?  
You’re welcome to open a discussion or suggest features via issues.

We highly appreciate your effort to contribute, but we recommend you talk to a maintainer before spending a lot of time making a pull request that may not align with the project roadmap.


Thanks again for contributing to LabraGo! 💚

## Maintainers

This project is maintained by:

[David](https://github.com/crembel) – actually does the hard backend stuff-

[Andrei](https://github.com/andreivinca) – tells us when JavaScript is broken (which is often)

[Dan](https://github.com/pedeceul) – occasionally pushes buttons and pretends to know what he's doing