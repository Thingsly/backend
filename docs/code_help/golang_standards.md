# Golang Conventions

## golang-standards

**Reference Repository**: [https://github.com/golang-standards](https://github.com/golang-standards)

**Documentation**: [https://github.com/golang-standards/project-layout/blob/master/README.md](https://github.com/golang-standards/project-layout/blob/master/README.md)

1. **Positive Feedback**: Many developers praise the clear structure and well-organized project layout provided. Especially for beginners and developers transitioning from other programming languages to Go, these guidelines serve as a great starting point, helping them quickly get started and follow generally accepted best practices.

2. **Criticism**: Some advanced Go developers point out that these structures might be too complex or unnecessary, especially for smaller or simpler projects. The philosophy of Go encourages simplicity and directness, and the template provided by `project-layout` may lead to excessive structuring and unnecessary complexity.

This project only partially references `project-layout`.

## Core Principles and Philosophy of Go

- **Simplicity ("Less is exponentially more")**: This principle, highlighted by Rob Pike in his famous blog post, emphasizes that by reducing the number of elements in the language, greater expressiveness and flexibility can be achieved.

- **Efficient Concurrency ("Concurrency is not parallelism")**: Pike stresses the distinction between concurrency and parallelism. Go's concurrency model (goroutines and channels) is designed as an efficient way to manage concurrent operations, not just to achieve multi-core parallel computing.

- **Clear Error Handling ("Clear is better than clever")**: This core idea in Go emphasizes code clarity and readability. For instance, Go avoids traditional exception handling mechanisms and instead uses explicit error checks, making error handling clearer and more direct.

- **Unified Code Formatting ("Gofmt's style is no one's favorite, yet gofmt is everyone's favorite")**: Go introduces the `gofmt` tool, which enforces a uniform code format. The goal is to eliminate debates over code formatting, ensuring all Go code looks consistent and improving readability.

- **Avoid Hidden Complexity ("Do less. Enable more.")**: Another core principle of Go is to **avoid excessive abstraction and hidden complexity**. Go encourages direct solutions rather than methods that might seem clever but could hide complexity and potential issues.

- **Tools and Community ("The bigger the interface, the weaker the abstraction")**: This principle is about interface design. It encourages **creating small, focused interfaces rather than large, all-encompassing ones**. This aligns with the Go community and toolchain's design philosophy of providing simple, effective tools and encouraging community contributions and collaboration.
