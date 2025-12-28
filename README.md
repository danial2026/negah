# THE WATCHMAN - Network Guardian

**The Watchman** is a beautiful terminal-based network scanner written in Go.
It wraps **Nmap** into an interactive TUI (Terminal User Interface) so you don't have to remember flags, scripts, or arcane syntax.

> ✨ For latest updates, check out the [dev](https://github.com/danial2026/the_watchman/tree/dev) branch ✨

![Screenshot 1](./screenshots/screenshot-1.png)

## ✨ Features

### Network Scanning Capabilities
* 🔍 Discover hosts on a local network
* 💓 Check if machines are alive
* 🌐 Get public IP and ISP info
* 🔎 Find subdomains
* 🛡️ Run common vulnerability checks
* 🔐 Audit SSL/TLS and SSH
* 🚧 Test firewall behavior
* 🥷 Perform stealth scans
* 🗄️ Detect databases and services
* 📊 Identify service versions and web titles

### TUI Features
* ⌨️ Beautiful, modern terminal interface powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
* 📋 **Copy scan results to clipboard** with a single keypress
* 🎨 Syntax highlighting and styled output
* 📱 Responsive design that adapts to terminal size
* 🎯 Interactive table navigation
* ⚡ Real-time scan progress indicators
* ❓ Built-in help menu

## Requirements

* **Go 1.21+**
* **Nmap**
  * macOS: `brew install nmap`
  * Arch Linux: `sudo pacman -S nmap`
  * Ubuntu/Debian: `sudo apt install nmap`

---

## 🚀 Build and Run

```bash
bash run.sh
```

Or build manually:
```bash
go build -o nscanner
./nscanner
```

---

## ⌨️ Keyboard Shortcuts

### Menu View
* `↑/↓` or `j/k` - Navigate through scan options
* `Enter` - Select a scan
* `q` - Quit application
* `?` - Toggle help menu

### Input View
* `Enter` - Confirm and continue
* `Esc` - Back to menu

### Result View
* `c` - **Copy results to clipboard**
* `Enter` or `Esc` - Back to menu

### Global
* `?` - Toggle help menu
* `Ctrl+C` - Force quit

---

## 📖 Usage

1. Launch the application
2. Use arrow keys to navigate through available scans
3. Press `Enter` to select a scan
4. Enter target information (IP, domain, or range) when prompted
5. Wait for the scan to complete
6. Press `c` to copy results to clipboard
7. Press `Enter` to return to the menu

If a scan needs root access, the tool handles `sudo` automatically on macOS and Linux.

---

## ⚠️ Important

I don't condone doing bad things.
Don't hack people. Don't scan random networks.

Use this **only** on systems you own or have explicit permission to test.

---

## 🛠️ Tech Stack

* [Go](https://golang.org/) - Programming language
* [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
* [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
* [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
* [clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard access
* [Nmap](https://nmap.org/) - Network scanning engine
