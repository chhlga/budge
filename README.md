# budge

A terminal-based email client built with Go and Bubble Tea.

## Features

- 📬 List mailboxes and navigate folders
- 📧 Browse emails with pagination
- 📖 Read emails with HTML rendering (via Markdown)
- 🔍 Search emails (UI implemented, IMAP backend pending)
- ⌨️ Vim-style keyboard navigation
- 🎨 Clean TUI interface with Charm libraries

## Installation

### Prerequisites

- Go 1.22 or higher
- IMAP email account

### Build from source

```bash
git clone https://github.com/chhlga/budge.git
cd budge
go build -o budge .
```

## Configuration

1. Create the config directory:
```bash
mkdir -p ~/.config/budge
```

2. Copy the example config:
```bash
cp config.example.yaml ~/.config/budge/config.yaml
```

3. Edit the config file with your IMAP credentials:
```bash
nano ~/.config/budge/config.yaml
```

4. Secure the config file (contains password):
```bash
chmod 600 ~/.config/budge/config.yaml
```

### Example Configuration

```yaml
server:
  host: imap.gmail.com
  port: 993                    # 993 for TLS, 143 for STARTTLS/plain
  tls: true                    # Implicit TLS (direct secure connection)
  starttls: false              # STARTTLS (upgrade from plain to secure)

credentials:
  username: your.email@gmail.com
  password: your-app-specific-password

behavior:
  default_folder: INBOX
  page_size: 50

display:
  date_format: "Jan 02 15:04"
  theme: auto
```

**Security Options**:
- **TLS (port 993)**: Implicit TLS - connection is encrypted from the start (recommended for Gmail)
- **STARTTLS (port 143)**: Starts as plain connection, then upgrades to TLS
- **Plain (port 143)**: Insecure - only use for testing/debugging

**Note**: Cannot enable both `tls` and `starttls` - choose one based on your server's configuration.

**Note for Gmail users**: You need to generate an [App Password](https://myaccount.google.com/apppasswords) instead of using your regular password.

## Usage

Run budge:
```bash
./budge
```

### Keyboard Shortcuts

#### Global
- `q` or `Ctrl+C` - Quit
- `r` or `Ctrl+R` - Refresh
- `?` - Help

#### Navigation
- `1` - View mailboxes
- `2` - View emails
- `3` - View reader
- `/` - Search
- `↑`/`k` - Move up
- `↓`/`j` - Move down
- `Enter` - Select
- `Esc` - Back

## Architecture

```
budge/
├── main.go                    # Entry point
├── config.example.yaml        # Configuration template
├── internal/
│   ├── config/               # YAML config loading & validation
│   ├── email/                # Email types, parser, renderer
│   │   ├── types.go
│   │   ├── parser.go
│   │   ├── renderer.go
│   │   └── testdata/         # .eml test fixtures
│   ├── imap/                 # IMAP client wrapper
│   │   ├── client.go         # Connection management & reconnection
│   │   └── errors.go
│   ├── cache/                # LRU cache for email bodies
│   └── tui/                  # Bubble Tea TUI components
│       ├── model.go          # Root model with state machine
│       ├── mailboxes.go      # Mailbox list view
│       ├── emails.go         # Email list view
│       ├── reader.go         # Email reader with viewport
│       ├── search.go         # Search interface
│       ├── statusbar.go      # Status bar component
│       ├── keys.go           # Key bindings
│       ├── styles.go         # Lip Gloss styles
│       └── messages.go       # Custom message types
```

## Technology Stack

- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Elm architecture for terminals
- **TUI Components**: [Bubbles](https://github.com/charmbracelet/bubbles) - Pre-built UI components (list, viewport, textinput)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling and layout
- **Markdown Rendering**: [Glamour](https://github.com/charmbracelet/glamour) - Styled markdown output
- **IMAP**: [go-imap v2](https://github.com/emersion/go-imap) - IMAP4rev2 client
- **Email Parsing**: [go-message](https://github.com/emersion/go-message) - MIME message parsing
- **HTML Conversion**: [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) - HTML to Markdown conversion

## Development Status

### ✅ Phase 1: Foundation Layer (COMPLETE)
- [x] Configuration management with YAML
- [x] IMAP client wrapper with reconnection logic
- [x] Email parsing (plain text, HTML, multipart)
- [x] HTML → Markdown → ANSI rendering pipeline
- [x] LRU cache for email bodies
- [x] 37 tests passing

### ✅ Phase 2: TUI Layer (COMPLETE)
- [x] TUI root model with state machine
- [x] Mailbox list view
- [x] Email list view with pagination support
- [x] Email reader with scrolling viewport
- [x] Search UI
- [x] Status bar with connection health
- [x] Vim-style keyboard bindings
- [x] Binary compiles successfully (21MB)

### ✅ Phase 3: IMAP Integration (COMPLETE)
- [x] Async IMAP operations via tea.Cmd (non-blocking)
- [x] Connect and authenticate on startup
- [x] Load mailboxes on connection
- [x] Fetch emails from selected mailbox with pagination
- [x] Render email bodies with HTML support
- [x] IMAP SEARCH integration
- [x] Mark emails as read/unread (m key)
- [x] Delete emails (d key)
- [x] Search functionality (/ key)
- [x] All keyboard actions wired to IMAP operations

### 🎯 Ready for Testing
The application is **functionally complete** and ready for end-to-end testing with a real IMAP server. All core features are implemented and the TUI is fully interactive.

**Next Step**: Test with Gmail using an app-specific password to verify real-world functionality.

### 🔮 Future Enhancements
- [ ] Compose new emails
- [ ] Reply to emails
- [ ] Forward emails
- [ ] Attachments support
- [ ] Multiple account support
- [ ] Offline mode with local cache
- [ ] IMAP IDLE (real-time notifications)

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Current test coverage:
- **37 tests** across foundation packages
- Config loading and validation
- Email parsing (plain, HTML, multipart)
- HTML rendering pipeline
- IMAP client state management
- LRU cache functionality

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Format code: `gofmt -w .`
6. Submit a pull request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Charm](https://charm.sh/) for the excellent TUI libraries
- [emersion](https://github.com/emersion) for go-imap and go-message
- The Go community for awesome tooling

## Contact

Issues and feature requests: [GitHub Issues](https://github.com/chhlga/budge/issues)
