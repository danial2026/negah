# Refactoring Notes - TUI Implementation

## Overview
Negah has been completely refactored to use [Bubble Tea](https://github.com/charmbracelet/bubbletea) for a modern, interactive Terminal User Interface (TUI). The application now features a beautiful interface with clipboard support, styled output, and improved user experience.

## What Changed

### 1. New Dependencies Added
- `github.com/charmbracelet/bubbletea` - Terminal UI framework (Elm architecture for Go)
- `github.com/charmbracelet/lipgloss` - Style definitions and color management
- `github.com/charmbracelet/bubbles` - Pre-built TUI components (table, text input)
- `github.com/atotto/clipboard` - Cross-platform clipboard access

### 2. New Files Created

#### `tui.go` (490+ lines)
Complete TUI implementation with:
- **Model struct**: Manages application state
- **Multiple views**: Menu, Input, Running, Result, and Help views
- **Styled components**: Custom color scheme and ASCII art header
- **Keyboard shortcuts**: Comprehensive keybindings for navigation
- **Real-time updates**: Animated spinner during scans
- **Clipboard integration**: One-key copy of scan results

### 3. Modified Files

#### `main.go`
- **Before**: Simple text-based menu with `bufio.Reader`
- **After**: Launches Bubble Tea program with alternate screen mode
- Reduced from ~70 lines to ~20 lines
- Much cleaner and more maintainable

#### `scanner/scanner.go`
Added new functions for output capture:
- `ExecuteCommandWithOutput()` - Captures command output to string
- `GetPublicIPWithOutput()` - Returns IP info as string
- `GetLocalInfoWithOutput()` - Returns interface info as string
- `RunScanWithOutput()` - Main scan function with output capture

Original functions preserved for backward compatibility.

#### `README.md`
- Updated with TUI features and keyboard shortcuts
- Added emoji icons for better visual appeal
- Documented clipboard functionality
- Added tech stack section
- Improved formatting and structure

## Key Features Implemented

### 🎨 Beautiful UI
- Custom color scheme (cyan, gold, green accents)
- Styled ASCII art header
- Professional table layout
- Bordered result display
- Responsive to terminal size

### ⌨️ Keyboard Shortcuts
- **Menu Navigation**: `↑/↓` or `j/k`
- **Selection**: `Enter`
- **Copy to Clipboard**: `c` (in result view)
- **Help Toggle**: `?`
- **Back**: `Esc`
- **Quit**: `q` or `Ctrl+C`

### 📋 Clipboard Support
Results can be copied with a single keypress (`c`):
- Full scan output
- Command information
- Error messages
- IP/network info

### 🎯 User Experience Improvements
- No more remembering tool numbers
- Visual feedback during scans
- Clear navigation hints
- Contextual help system
- Error handling with styled messages
- Smooth transitions between views

## Architecture

### Bubble Tea Pattern (Elm Architecture)
```
Model → Update → View
  ↑        ↓
  └────────┘
     Msg
```

1. **Model**: Application state
2. **Update**: Handles messages and updates state
3. **View**: Renders current state to string
4. **Messages**: Events (keypresses, tick, scan complete)

### View States
1. **menuView**: Browse available scans
2. **inputView**: Enter target/ports
3. **runningView**: Animated progress indicator
4. **resultView**: Display results with copy option
5. **helpView**: Keyboard shortcuts reference

## Technical Highlights

### Message Passing
- `tea.KeyMsg` - Keyboard events
- `tea.WindowSizeMsg` - Terminal resize
- `tickMsg` - Animation frames
- `scanCompleteMsg` - Async scan completion

### Concurrent Execution
Scans run in goroutines via `tea.Cmd`:
```go
func() tea.Msg {
    output, err := scanner.RunScanWithOutput(m.selected, m.targetInput)
    return scanCompleteMsg{output: output, err: err}
}
```

### Styling with Lipgloss
```go
lipgloss.NewStyle().
    Foreground(ColorPrimary).
    Bold(true).
    Align(lipgloss.Center)
```

## Testing

Build and run:
```bash
go build -o nscanner
./nscanner
```

Or use the run script:
```bash
bash run.sh
```

## Backward Compatibility

All original scanner functions remain intact:
- `RunScan()` - Still works for direct execution
- `GetPublicIP()` - Still prints to stdout
- `GetLocalInfo()` - Still prints to stdout

New `*WithOutput()` variants don't break existing code.

## Performance

- Lazy rendering (only when state changes)
- Efficient string building
- Minimal allocations
- Fast table rendering via `bubbles/table`

## Future Enhancements (Optional)

- [ ] History of recent scans
- [ ] Save results to file
- [ ] Custom scan profiles
- [ ] Multi-target scanning
- [ ] Progress bars for long scans
- [ ] Syntax highlighting for results
- [ ] Export to JSON/CSV
- [ ] Dark/Light theme toggle

## Credits

Inspired by:
- [Build a System Monitor TUI in Go](https://penchev.com/posts/create-tui-with-go/) by Ivan Penchev
- [Charm](https://charm.sh/) - Bubble Tea framework creators

## Developer Notes

### Adding New Scans
1. Add to `scanner.GetFeatures()`
2. No TUI changes needed - automatically appears in table

### Modifying Colors
Edit color constants in `tui.go`:
```go
var (
    ColorPrimary   = lipgloss.Color("#00CED1")
    ColorSecondary = lipgloss.Color("#F6D2A2")
    // ... etc
)
```

### Adding New Views
1. Add view state to `viewState` enum
2. Implement `view*()` method
3. Add case to `View()` switch
4. Handle transitions in `Update()`

---

**Refactored by**: AI Assistant  
**Date**: December 27, 2025  
**Framework**: Bubble Tea v0.25.0  
**Go Version**: 1.23.1

