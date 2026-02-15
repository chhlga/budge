# Agent Guide for budge

This guide helps AI agents work effectively in the **budge** codebase - a terminal-based email client built with Go and Bubble Tea.

## Project Overview

**budge** is a lightweight, keyboard-driven TUI email client with:
- IMAP support (Gmail, Outlook, self-hosted servers)
- Rich HTML email rendering via Markdown
- Vim-style navigation (hjkl, gg, G, /)
- Built with Charm's Bubble Tea framework
- Async operations with push notifications

**Tech Stack:**
- Go 1.25.6+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework (Elm architecture)
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [go-imap/v2](https://github.com/emersion/go-imap) - IMAP client library

## Essential Commands

### Build & Run
```bash
# Build the binary
go build -o budge .

# Run directly
go run .

# Install to Go bin
go install github.com/chhlga/budge@latest
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/tui/
go test ./internal/imap/
```

### Development
```bash
# Download dependencies
go mod download

# Tidy dependencies (fixes the "should be direct" warning)
go mod tidy

# Format code
gofmt -w .

# Check for common issues
go vet ./...
```

## Project Structure

```
budge/
├── main.go                    # Entry point - loads config, creates client, starts TUI
├── config.example.yaml        # Example configuration template
├── internal/
│   ├── cache/                # LRU cache for email bodies (thread-safe)
│   │   ├── cache.go          # Cache implementation
│   │   └── cache_test.go
│   ├── config/               # YAML configuration management
│   │   ├── config.go         # Config struct, loading, validation
│   │   └── config_test.go
│   ├── email/                # Email parsing and rendering
│   │   ├── parser.go         # Parse email messages (headers, body)
│   │   ├── renderer.go       # HTML→Markdown→ANSI rendering
│   │   ├── types.go          # Email data structures
│   │   └── *_test.go
│   ├── imap/                 # IMAP client wrapper
│   │   ├── client.go         # Connection, auth, mailbox ops, search
│   │   ├── errors.go         # Custom error types
│   │   └── client_test.go
│   └── tui/                  # Terminal UI (Bubble Tea)
│       ├── model.go          # Root model (Elm architecture)
│       ├── commands.go       # tea.Cmd functions (async IMAP operations)
│       ├── messages.go       # tea.Msg types (events)
│       ├── keys.go           # Keyboard bindings
│       ├── styles.go         # Lip Gloss styles
│       ├── mailboxes.go      # Mailbox list view
│       ├── emails.go         # Email list view (sorting, filtering)
│       ├── reader.go         # Email reader view
│       ├── search.go         # Search input view
│       ├── statusbar.go      # Status bar component
│       ├── flags.go          # Email flag management
│       └── *_test.go
```

## Code Patterns & Conventions

### Bubble Tea Architecture (Elm Pattern)

budge follows the **Elm architecture** via Bubble Tea:

1. **Model** - Application state (`Model` struct in `internal/tui/model.go`)
2. **Update** - State transitions based on messages (`Update(msg tea.Msg)`)
3. **View** - Render current state to terminal (`View() string`)
4. **Commands** - Side effects that return messages (`tea.Cmd`)

**Key Components:**
- `Model.state` - Current view (mailboxListView, emailListView, emailReaderView, searchView)
- `Model.focus` - Which pane has focus (focusLeft, focusRight)
- Sub-models: `mailboxList`, `emailList`, `emailReader`, `search`, `statusBar`

### Command Pattern (Async Operations)

All IMAP operations are **async commands** returning `tea.Cmd`:

```go
// Example from commands.go
func loadEmailsCmd(client *imapClient.Client, mailbox string, pageSize uint32) tea.Cmd {
    return func() tea.Msg {
        // Perform IMAP operation (blocking)
        emails, err := client.FetchEmails(...)
        
        // Return message with results
        if err != nil {
            return ErrorMsg{Err: err}
        }
        return EmailsLoadedMsg{Emails: emails}
    }
}
```

**Pattern:**
1. Commands in `commands.go` wrap IMAP operations
2. Commands return messages (`*Msg` types in `messages.go`)
3. `Update()` handles messages and updates state
4. Never block in `Update()` or `View()` - use commands for I/O

### Message Types

All messages are defined in `messages.go`:
```go
type ConnectCompleteMsg struct{}
type ConnectErrorMsg struct{ Err error }
type MailboxesLoadedMsg struct{ Mailboxes []string }
type EmailsLoadedMsg struct{ Emails []EmailEnvelope }
type EmailBodyLoadedMsg struct{ UID uint32; Body string }
// ... etc
```

**Convention:** 
- Success messages: `*LoadedMsg`, `*CompleteMsg`
- Error messages: `*ErrorMsg` with `Err error` field

### State Management

View state transitions:
```go
type viewState uint

const (
    mailboxListView viewState = iota
    emailListView
    emailReaderView
    searchView
)
```

Focus management (dual-pane layout):
```go
type paneFocus uint

const (
    focusLeft paneFocus = iota  // Mailbox list
    focusRight                  // Email list or reader
)
```

### Styling Conventions

Styles are centralized in `styles.go`:
```go
var (
    primaryColor = lipgloss.Color("#7D56F4")
    TitleStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
    // ... etc
)
```

**Current deprecation warnings:** 
- `Style.Copy()` is deprecated - use direct assignment (`a := b`) instead
- Affects `HeaderStyle.Copy()` and focus styles in `model.go:358, 412, 413, 416, 417`

### Email Rendering Pipeline

1. **Parse** (`email/parser.go`): Extract headers, text, HTML from MIME message
2. **Render** (`email/renderer.go`): HTML → Markdown → ANSI-styled terminal output
3. **Cache** (`cache/cache.go`): LRU cache for rendered email bodies
4. **Display** (`tui/reader.go`): Viewport with scrolling

### IMAP Client Wrapper

`internal/imap/client.go` wraps `go-imap/v2`:
- Connection states: `StateDisconnected`, `StateConnecting`, `StateConnected`, `StateAuthenticated`
- Thread-safe with `sync.RWMutex`
- Reconnection with exponential backoff
- Push notifications via `UpdateHandler`

### Testing Conventions

Tests follow Go conventions:
- `*_test.go` files alongside source
- Test functions: `TestFunctionName(t *testing.T)`
- Table-driven tests for multiple cases
- Mock IMAP server in `client_test.go`

Example from `cache_test.go`:
```go
func TestCache_LRUEviction(t *testing.T) {
    cache := New(3)
    // ... test logic
}
```

### Configuration

Config at `~/.config/budge/config.yaml`:
```yaml
server:
  host: imap.gmail.com
  port: 993
  tls: true
  starttls: false

credentials:
  username: user@example.com
  password: app-password

behavior:
  default_folder: INBOX
  page_size: 50
  poll_interval: 30  # seconds

display:
  date_format: "Jan 02 15:04"
  theme: auto
```

**Security:** Config file should be `chmod 600` (checked at startup)

## Key Features Implementation

### Vim-Style Keybindings

Defined in `keys.go`:
- Navigation: `j/k` (down/up), `gg` (top), `G` (bottom)
- Search: `/` (opens search view)
- Actions: `m` (mark read/unread), `d` (delete)
- View switching: `1` (mailboxes), `2` (emails), `3` (reader)
- Global: `q` (quit), `r` (refresh), `Tab` (switch focus)

### Email Sorting & Filtering

In `emails.go`:
```go
type SortMode int
const (
    SortDateNewest SortMode = iota
    SortDateOldest
    SortSenderAZ
    // ... etc (8 modes total)
)

type FilterMode int
const (
    FilterNone FilterMode = iota
    FilterUnread
    FilterRead
    FilterAttachments
)
```

Cycling: `s` key cycles through sort modes, `f` cycles through filters

### Search Implementation

IMAP SEARCH query:
1. User enters query in search view
2. `searchEmailsCmd` calls `client.Search()`
3. Results stored in `EmailsLoadedMsg`
4. `inSearchResults` flag tracks search mode
5. `Esc` exits search and restores previous email list

### Push Notifications

Polling mechanism in `commands.go`:
```go
func startMonitoringCmd(client *Client, mailbox string, interval time.Duration) tea.Cmd {
    return func() tea.Msg {
        ticker := time.NewTicker(interval)
        // Polls mailbox for new messages
        // Returns NewMailMsg when count increases
    }
}
```

Configured via `poll_interval` in config (default 30s)

## Common Tasks

### Adding a New View

1. Add view state constant to `model.go`
2. Create component struct with `Update()`, `View()`, `Init()` methods
3. Add to `Model` struct as sub-model
4. Handle state transitions in `Model.Update()`
5. Add keybinding in `keys.go`

### Adding a New IMAP Operation

1. Add method to `Client` in `internal/imap/client.go`
2. Create command function in `internal/tui/commands.go`
3. Define message types in `internal/tui/messages.go`
4. Handle messages in `Model.Update()` in `model.go`
5. Add tests in `internal/imap/client_test.go`

### Adding a New Keyboard Shortcut

1. Add binding to `KeyMap` struct in `keys.go`
2. Initialize in `NewKeyMap()`
3. Handle key in appropriate `Update()` method
4. Update help text

### Fixing Styling Issues

1. Check `styles.go` for existing styles
2. Use `lipgloss.NewStyle()` for new styles
3. Apply with `.Render(text)` or `.Width(w).Height(h)`
4. Test with different terminal sizes (resize events)

### Working with Email Rendering

1. HTML emails: `html-to-markdown` → `glamour` → ANSI
2. Plain text: `renderPlainText()` preserves formatting
3. Sanitization in `sanitizeMarkdown()` removes broken links/images
4. Caching in `cache.Cache` keyed by UID

## Known Issues & Gotchas

### Deprecated Style.Copy()

**Issue:** Lip Gloss deprecated `.Copy()` method
**Files affected:** `internal/tui/model.go` lines 358, 412-417
**Fix:** Replace `.Copy()` with direct assignment:
```go
// Old
newStyle := oldStyle.Copy().Foreground(color)

// New
newStyle := oldStyle
newStyle = newStyle.Foreground(color)
```

### go.mod Warning

**Issue:** `github.com/muesli/reflow should be direct`
**Fix:** Run `go mod tidy` (already indirect, likely unused)

### IMAP Connection Timeouts

**Context:** 30-second timeout in `connectCmd()`
**Pattern:** All IMAP commands use `context.WithTimeout()`
**Reason:** Prevents hanging on slow/unresponsive servers

### Bubble Tea Message Handling

**Gotcha:** Messages can arrive in any order (async)
**Pattern:** Check state before processing (e.g., `if !client.IsConnected()`)
**Example:** `loadMailboxesCmd` checks connection before fetching

### Email UID vs Sequence Number

**Important:** Always use UID (unique identifier), not sequence number
**Reason:** Sequence numbers change when emails are deleted
**IMAP methods:** Use `FetchMessagesByUID`, not `FetchMessagesBySeqNum`

### Terminal Resize Handling

**Pattern:** `tea.WindowSizeMsg` is sent on resize
**Action:** Update `Model.width` and `Model.height`
**Propagate:** Call `SetSize()` on all sub-components

### Testing with Mock IMAP

**Approach:** `client_test.go` uses mock server
**Note:** Real IMAP testing requires credentials
**Manual testing:** Use config.yaml with test account

## Development Workflow

### Standard Feature Development

1. Read relevant code (use `grep`, `go doc`, or LSP tools)
2. Write tests first (TDD) in `*_test.go`
3. Implement feature in source files
4. Run tests: `go test ./...`
5. Build and manually test: `go build -o budge . && ./budge`
6. Update README.md if user-facing
7. Run `gofmt -w .` before committing

### Debugging TUI Issues

1. **Use logs:** Write to file (don't print to stdout - breaks TUI)
2. **State inspection:** Add temporary fields to track state
3. **Message tracing:** Log all messages in `Update()`
4. **Resize testing:** Shrink/expand terminal while running

### CI/CD

GitHub Actions workflow (`.github/workflows/go.yml`):
- Triggers on push/PR to `main`
- Go 1.25.6
- Runs `go build -v ./...`
- Runs `go test -v ./...`

### Release Process

Not automated yet - see CONTRIBUTING.md for manual process

## External Resources

- [Bubble Tea Docs](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [go-imap/v2 Docs](https://pkg.go.dev/github.com/emersion/go-imap/v2)
- [IMAP RFC 3501](https://datatracker.ietf.org/doc/html/rfc3501)

## Questions & Support

- Issues: [GitHub Issues](https://github.com/chhlga/budge/issues)
- Contributing: See `CONTRIBUTING.md`
- Security: See `SECURITY.md`

---

**Last Updated:** 2025-02-14
**Go Version:** 1.25.6
**Main Dependencies:** Bubble Tea 1.3.10, go-imap v2.0.0-beta.8
