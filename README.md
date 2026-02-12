<a id="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/chhlga/budge">
    <img src="assets/blobashik_read.png" alt="Logo" width="254" height="254">
  </a>

<h3 align="center">budge</h3>

  <p align="center">
    Terminal-based email client with Vim keybindings
  </p>
</div>

### Built With
[![Go][golang-shield]][golang-url] [![Bubble Tea][bubbletea-shield]][bubbletea-url] [![Bubbles][bubbles-shield]][bubbles-url] [![Lip Gloss][lipgloss-shield]][lipgloss-url] [![Glamour][glamour-shield]][glamour-url]

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#tldr">TL;DR</a>
      <ul>
        <li><a href="#what-im-looking-at">What I'm Looking At?</a></li>
        <li><a href="#install">Install</a></li>
        <li><a href="#run">Run</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#usage">Usage</a></li>
        <li><a href="#configuration">Configuration</a></li>
      </ul>
    </li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>


<!-- TL;DR -->
## TL;DR
### What I'm Looking At?

**budge** is a terminal-based email client built with Go and Bubble Tea.

**Key features:**
- 📬 IMAP support - Gmail, Outlook, self-hosted servers
- 📧 Rich HTML emails - rendered as styled Markdown in terminal
- ⌨️ Vim-style navigation - hjkl, gg, G, / for search
- 🎨 Beautiful TUI - built with Charm's Bubble Tea framework
- 🔍 Email search - full IMAP SEARCH integration
- ⚡ Fast & lightweight - ~21MB binary, async operations

### Install

```bash
# Via Go install
go install github.com/chhlga/budge@latest
```

### Run

```bash
# First-time setup
mkdir -p ~/.config/budge
cp config.example.yaml ~/.config/budge/config.yaml
nano ~/.config/budge/config.yaml  # Add your IMAP credentials
chmod 600 ~/.config/budge/config.yaml

# Run budge
./budge
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running, follow these steps.

### Prerequisites

* **Go 1.25+** - Only if building from source
* **IMAP email account** - Gmail, Outlook, or any IMAP server
* **App password** - Gmail and most providers require app-specific passwords

### Installation

**Option 1: Go Install**

```sh
go install github.com/chhlga/budge@latest
```

**Option 2: Build from Source**

1. Clone the repo
   ```sh
   git clone https://github.com/chhlga/budge.git
   ```
2. Build
   ```sh
   cd budge
   go build -o budge .
   ```
3. Optionally install to system path
   ```sh
   sudo mv budge /usr/local/bin/
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

### Initial Setup

1. Create config directory:
   ```sh
   mkdir -p ~/.config/budge
   ```

2. Copy example config:
   ```sh
   cp config.example.yaml ~/.config/budge/config.yaml
   ```

3. Edit with your IMAP credentials:
   ```sh
   nano ~/.config/budge/config.yaml
   ```

4. Secure the config file:
   ```sh
   chmod 600 ~/.config/budge/config.yaml
   ```

5. Run budge:
   ```sh
   ./budge
   ```

### Keyboard Shortcuts

**Global**
- `q` or `Ctrl+C` - Quit
- `r` or `Ctrl+R` - Refresh
- `?` - Help
- `Esc` - Back

**Navigation**
- `1` - Mailboxes view
- `2` - Email list view
- `3` - Email reader
- `/` - Search
- `↑`/`k` - Move up
- `↓`/`j` - Move down
- `Enter` - Select

**Email Actions**
- `m` - Toggle read/unread
- `d` - Delete email
- `Space` - Page down

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Configuration

budge uses YAML configuration at `~/.config/budge/config.yaml`:

```yaml
server:
  host: imap.gmail.com
  port: 993                    # 993 for TLS, 143 for STARTTLS
  tls: true                    # Implicit TLS (recommended)
  starttls: false              # STARTTLS (upgrade from plain)

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

### Provider Examples

**Gmail**
```yaml
server:
  host: imap.gmail.com
  port: 993
  tls: true
  starttls: false
```
Generate [App Password](https://myaccount.google.com/apppasswords)

**Outlook / Office 365**
```yaml
server:
  host: outlook.office365.com
  port: 993
  tls: true
  starttls: false
```

**ProtonMail Bridge**
```yaml
server:
  host: 127.0.0.1
  port: 1143
  tls: false
  starttls: true
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Project Link: [https://github.com/chhlga/budge](https://github.com/chhlga/budge)

Issues: [https://github.com/chhlga/budge/issues](https://github.com/chhlga/budge/issues)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

Built with ❤️ using the Charm stack:

* [Charm](https://charm.sh/) - Bubble Tea, Bubbles, Lip Gloss, Glamour
* [emersion](https://github.com/emersion) - go-imap v2 and go-message
* [JohannesKaufmann](https://github.com/JohannesKaufmann) - html-to-markdown
* The Go community

* Approved by [Bashik](https://github.com/chhlga/chhlga/blob/main/bashik.md) 🐱

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & BADGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/chhlga/budge.svg?style=for-the-badge
[contributors-url]: https://github.com/chhlga/budge/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/chhlga/budge.svg?style=for-the-badge
[forks-url]: https://github.com/chhlga/budge/network/members
[stars-shield]: https://img.shields.io/github/stars/chhlga/budge.svg?style=for-the-badge
[stars-url]: https://github.com/chhlga/budge/stargazers
[issues-shield]: https://img.shields.io/github/issues/chhlga/budge.svg?style=for-the-badge
[issues-url]: https://github.com/chhlga/budge/issues
[license-shield]: https://img.shields.io/github/license/chhlga/budge.svg?style=for-the-badge
[license-url]: https://github.com/chhlga/budge/blob/master/LICENSE

[golang-shield]: https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white
[golang-url]: https://go.dev/
[bubbletea-shield]: https://img.shields.io/badge/Bubble%20Tea-TUI-FF69B4?style=for-the-badge
[bubbletea-url]: https://github.com/charmbracelet/bubbletea
[bubbles-shield]: https://img.shields.io/badge/Bubbles-Components-FF69B4?style=for-the-badge
[bubbles-url]: https://github.com/charmbracelet/bubbles
[lipgloss-shield]: https://img.shields.io/badge/Lip%20Gloss-Styling-FF69B4?style=for-the-badge
[lipgloss-url]: https://github.com/charmbracelet/lipgloss
[glamour-shield]: https://img.shields.io/badge/Glamour-Markdown-FF69B4?style=for-the-badge
[glamour-url]: https://github.com/charmbracelet/glamour
