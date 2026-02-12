# Contributing to budge

Thank you for considering contributing to budge! This document provides guidelines for contributing to this project.

## Code of Conduct

Please be respectful and constructive in all interactions. We're building an open, welcoming community.

## How to Contribute

### Reporting Bugs

Before creating a bug report:
- Check the [issue tracker](https://github.com/chhlga/budge/issues) for existing reports
- Verify you're using the latest version

When reporting a bug, include:
- Your operating system and Go version
- Steps to reproduce the issue
- Expected vs. actual behavior
- Email provider (Gmail, Outlook, etc.) if relevant
- Any error messages or logs

### Suggesting Features

Feature requests are welcome! Please:
- Search existing issues to avoid duplicates
- Clearly describe the use case and benefit
- Consider if it fits the project's scope (lightweight, keyboard-driven TUI email client)

### Pull Requests

1. **Fork the repository** and create a branch from `main`
2. **Make your changes**:
   - Follow existing code style and conventions
   - Add tests for new functionality
   - Update documentation if needed
3. **Test your changes**:
   ```bash
   go test ./...
   go build -o budge .
   ./budge  # Manual testing
   ```
4. **Commit your changes** with clear, descriptive messages
5. **Push to your fork** and submit a pull request

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Keep functions focused and reasonably sized
- Add comments for complex logic
- Handle errors explicitly; don't ignore them
- Use meaningful variable names

### Testing

- Add tests for new features in `*_test.go` files
- Ensure existing tests pass: `go test ./...`
- Test manually with different IMAP providers if possible

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/budge.git
cd budge

# Install dependencies
go mod download

# Build
go build -o budge .

# Run tests
go test ./...
```

## Project Structure

```
budge/
├── main.go              # Entry point
├── internal/
│   ├── cache/          # Email caching
│   ├── config/         # Configuration management
│   ├── email/          # Email parsing & rendering
│   ├── imap/           # IMAP client
│   └── tui/            # Terminal UI (Bubble Tea)
```

## Questions?

Open an issue or discussion on GitHub. We're happy to help!
