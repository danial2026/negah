package main

import (
	"fmt"
	"strings"
	"time"

	"the_watchman/scanner"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define color scheme
var (
	ColorPrimary    = lipgloss.Color("#00CED1")
	ColorSecondary  = lipgloss.Color("#F6D2A2")
	ColorAccent     = lipgloss.Color("#89E051")
	ColorError      = lipgloss.Color("#FF6B6B")
	ColorSuccess    = lipgloss.Color("#4ECDC4")
	ColorHighlight  = lipgloss.Color("#44475A")
	ColorForeground = lipgloss.Color("#F8F8F2")
)

// View states
type viewState int

const (
	menuView viewState = iota
	inputView
	runningView
	resultView
	helpView
)

// tickMsg is sent on every tick
type tickMsg time.Time

// scanCompleteMsg is sent when a scan completes
type scanCompleteMsg struct {
	output string
	err    error
}

// Model represents the TUI state
type model struct {
	// UI components
	table     table.Model
	textInput textinput.Model
	width     int
	height    int

	// State
	currentView  viewState
	tools        []scanner.ScanFeature
	selected     scanner.ScanFeature
	targetInput  string
	portsInput   string
	needsTarget  bool
	needsPorts   bool
	scanOutput   string
	scanError    error
	showHelp     bool
	message      string
	messageColor lipgloss.Color

	// Styles
	baseStyle    lipgloss.Style
	headerStyle  lipgloss.Style
	footerStyle  lipgloss.Style
	helpStyle    lipgloss.Style
	tableStyle   table.Styles
	inputStyle   lipgloss.Style
	messageStyle lipgloss.Style
	borderStyle  lipgloss.Style
}

// Init initializes the model
func initialModel() model {
	// Create table columns
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 25},
		{Title: "Description", Width: 50},
	}

	// Get all scan features
	tools := scanner.GetFeatures()

	// Create table rows
	rows := []table.Row{}
	for _, t := range tools {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", t.ID),
			t.Name,
			t.Description,
		})
	}

	// Create table styles
	tableStyles := table.DefaultStyles()
	tableStyles.Header = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ColorPrimary).
		BorderBottom(true).
		Bold(true).
		Foreground(ColorPrimary)
	tableStyles.Selected = lipgloss.NewStyle().
		Foreground(ColorForeground).
		Background(ColorHighlight).
		Bold(true)

	// Create table
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithStyles(tableStyles),
	)

	// Create text input
	ti := textinput.New()
	ti.Placeholder = "Enter target (IP, Domain, or Range)"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	ti.TextStyle = lipgloss.NewStyle().Foreground(ColorForeground)

	return model{
		table:        t,
		textInput:    ti,
		tools:        tools,
		currentView:  menuView,
		baseStyle:    lipgloss.NewStyle(),
		headerStyle:  lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Align(lipgloss.Center),
		footerStyle:  lipgloss.NewStyle().Foreground(ColorSecondary),
		helpStyle:    lipgloss.NewStyle().Foreground(ColorAccent),
		inputStyle:   lipgloss.NewStyle().Foreground(ColorForeground),
		messageStyle: lipgloss.NewStyle().Foreground(ColorSuccess),
		borderStyle:  lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(ColorPrimary),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate table height based on available space
		// Account for: header (21 or 1 lines), spacing, footer, messages
		headerLines := 21 // Full ASCII art
		if m.height < 30 {
			headerLines = 1 // Compact header
		}

		// Reserve space for: header + spacing (2) + footer (1) + messages (3) + padding (3)
		reservedLines := headerLines + 9
		tableHeight := m.height - reservedLines

		// Ensure minimum table height
		if tableHeight < 5 {
			tableHeight = 5
		}

		m.table.SetHeight(tableHeight)
		return m, nil

	case tea.KeyMsg:
		// Global keys
		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentView == menuView {
				return m, tea.Quit
			}
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		}

		// View-specific keys
		switch m.currentView {
		case menuView:
			return m.handleMenuKeys(msg)
		case inputView:
			return m.handleInputKeys(msg)
		case resultView:
			return m.handleResultKeys(msg)
		case runningView:
			if msg.String() == "esc" {
				m.currentView = menuView
				m.message = ""
			}
		}

	case scanCompleteMsg:
		m.scanOutput = msg.output
		m.scanError = msg.err
		m.currentView = resultView
		if msg.err != nil {
			m.message = "Scan completed with errors"
			m.messageColor = ColorError
		} else {
			m.message = "Scan completed! Press 'c' to copy to clipboard"
			m.messageColor = ColorSuccess
		}
		return m, nil

	case tickMsg:
		// Keep ticking while running
		if m.currentView == runningView {
			return m, tickEvery()
		}
		return m, nil
	}

	// Update table in menu view
	if m.currentView == menuView {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleMenuKeys handles keyboard input in menu view
func (m model) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Get selected tool
		selectedIdx := m.table.Cursor()
		m.selected = m.tools[selectedIdx]

		// Check if this is an internal tool that doesn't need target
		if m.selected.Command == "INTERNAL_PUBLIC_IP" || m.selected.Command == "INTERNAL_LOCAL_INFO" {
			m.needsTarget = false
			m.needsPorts = false
			return m, m.runScan()
		}

		// Check if we need ports input
		if m.selected.ID == 4 {
			m.needsPorts = true
		}

		// Move to input view
		m.needsTarget = true
		m.currentView = inputView
		m.textInput.SetValue("")
		m.textInput.Placeholder = "Enter target (IP, Domain, or Range)"
		m.textInput.Focus()
		m.message = ""
		return m, nil

	default:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}
}

// handleInputKeys handles keyboard input in input view
func (m model) handleInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		if m.needsTarget && m.targetInput == "" {
			// Save target and check if we need ports
			m.targetInput = m.textInput.Value()
			if m.needsPorts {
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Enter ports (e.g., 80,443 or 1-1000)"
				m.needsTarget = false
				return m, nil
			}
		} else if m.needsPorts && m.portsInput == "" {
			// Save ports
			m.portsInput = m.textInput.Value()
			m.selected.Command = "-p " + m.portsInput
		}

		// Run the scan
		return m, m.runScan()

	case "esc":
		m.currentView = menuView
		m.targetInput = ""
		m.portsInput = ""
		m.needsTarget = false
		m.needsPorts = false
		m.message = ""
		return m, nil

	default:
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
}

// handleResultKeys handles keyboard input in result view
func (m model) handleResultKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		// Copy to clipboard
		err := clipboard.WriteAll(m.scanOutput)
		if err != nil {
			m.message = "Failed to copy to clipboard"
			m.messageColor = ColorError
		} else {
			m.message = "✓ Copied to clipboard!"
			m.messageColor = ColorSuccess
		}
		return m, nil

	case "esc", "enter", "backspace":
		m.currentView = menuView
		m.targetInput = ""
		m.portsInput = ""
		m.needsTarget = false
		m.needsPorts = false
		m.scanOutput = ""
		m.message = ""
		return m, nil

	default:
		return m, nil
	}
}

// runScan executes the selected scan
func (m model) runScan() tea.Cmd {
	m.currentView = runningView
	m.message = "Running scan..."
	m.messageColor = ColorAccent

	return tea.Batch(
		tickEvery(),
		func() tea.Msg {
			output, err := scanner.RunScanWithOutput(m.selected, m.targetInput)
			return scanCompleteMsg{output: output, err: err}
		},
	)
}

// tickEvery returns a command that sends a tick message
func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// View renders the UI
func (m model) View() string {
	if m.showHelp {
		return m.viewHelp()
	}

	switch m.currentView {
	case menuView:
		return m.viewMenu()
	case inputView:
		return m.viewInput()
	case runningView:
		return m.viewRunning()
	case resultView:
		return m.viewResult()
	default:
		return ""
	}
}

// viewMenu renders the menu view
func (m model) viewMenu() string {
	var b strings.Builder

	// Header (with adaptive spacing)
	header := m.headerStyle.Width(m.width).Render(m.renderHeader())
	b.WriteString(header)
	if m.height < 30 {
		b.WriteString("\n")
	} else {
		b.WriteString("\n\n")
	}

	// Check nmap
	if !scanner.CheckNmap() {
		warning := lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true).
			Render("⚠ Warning: 'nmap' not found in PATH. Install it first!")
		b.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(warning))
		b.WriteString("\n")
	}

	// Table
	b.WriteString(m.table.View())
	b.WriteString("\n")

	// Message
	if m.message != "" {
		msg := m.messageStyle.Foreground(m.messageColor).Render(m.message)
		b.WriteString(msg)
		b.WriteString("\n")
	}

	// Footer
	footer := m.footerStyle.Render(
		"↑/↓: Navigate • Enter: Select • ?: Help • q: Quit",
	)
	b.WriteString(footer)

	return b.String()
}

// viewInput renders the input view
func (m model) viewInput() string {
	var b strings.Builder

	// Header (with adaptive spacing)
	header := m.headerStyle.Width(m.width).Render(m.renderHeader())
	b.WriteString(header)
	if m.height < 30 {
		b.WriteString("\n")
	} else {
		b.WriteString("\n\n")
	}

	// Selected tool info
	toolInfo := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Render(fmt.Sprintf("Selected: %s - %s", m.selected.Name, m.selected.Description))
	b.WriteString(toolInfo)
	b.WriteString("\n")

	// Input
	b.WriteString(m.textInput.View())
	b.WriteString("\n")

	// Footer
	footer := m.footerStyle.Render("Enter: Continue • Esc: Back to menu")
	b.WriteString(footer)

	return b.String()
}

// viewRunning renders the running view
func (m model) viewRunning() string {
	var b strings.Builder

	// Header (with adaptive spacing)
	header := m.headerStyle.Width(m.width).Render(m.renderHeader())
	b.WriteString(header)
	if m.height < 30 {
		b.WriteString("\n")
	} else {
		b.WriteString("\n\n")
	}

	// Spinner and message
	spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"
	frame := int(time.Now().UnixNano()/1e8) % len(spinner)

	msg := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Render(fmt.Sprintf("%c Running: %s", rune(spinner[frame]), m.selected.Name))
	b.WriteString(msg)
	b.WriteString("\n")

	target := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Render(fmt.Sprintf("Target: %s", m.targetInput))
	b.WriteString(target)
	b.WriteString("\n")

	// Footer
	footer := m.footerStyle.Render("Please wait... • Esc: Cancel")
	b.WriteString(footer)

	return b.String()
}

// viewResult renders the result view
func (m model) viewResult() string {
	var b strings.Builder

	// Header (with adaptive spacing)
	header := m.headerStyle.Width(m.width).Render(m.renderHeader())
	b.WriteString(header)
	if m.height < 30 {
		b.WriteString("\n")
	} else {
		b.WriteString("\n\n")
	}

	// Tool info
	toolInfo := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Render(fmt.Sprintf("Scan: %s | Target: %s", m.selected.Name, m.targetInput))
	b.WriteString(toolInfo)
	b.WriteString("\n")

	// Results box - calculate height based on available space
	headerLines := 21
	if m.height < 30 {
		headerLines = 1
	}
	resultBoxHeight := m.height - headerLines - 8 // header + spacing + toolInfo + message + footer + padding
	if resultBoxHeight < 5 {
		resultBoxHeight = 5
	}

	resultBox := m.borderStyle.
		Width(m.width - 4).
		Height(resultBoxHeight).
		Padding(1).
		Render(m.scanOutput)
	b.WriteString(resultBox)
	b.WriteString("\n")

	// Message
	if m.message != "" {
		msg := m.messageStyle.Foreground(m.messageColor).Render(m.message)
		b.WriteString(msg)
		b.WriteString("\n")
	}

	// Footer
	footer := m.footerStyle.Render("c: Copy to Clipboard • Enter/Esc: Back to menu")
	b.WriteString(footer)

	return b.String()
}

// viewHelp renders the help view
func (m model) viewHelp() string {
	var b strings.Builder

	// Header (with adaptive spacing)
	header := m.headerStyle.Width(m.width).Render(m.renderHeader())
	b.WriteString(header)
	if m.height < 30 {
		b.WriteString("\n")
	} else {
		b.WriteString("\n\n")
	}

	// Help title
	helpTitle := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Render("⌨ Keyboard Shortcuts")
	b.WriteString(helpTitle)
	b.WriteString("\n")

	// Help content
	helpContent := [][]string{
		{"Menu View", ""},
		{"  ↑/↓, j/k", "Navigate through scan options"},
		{"  Enter", "Select a scan"},
		{"  q", "Quit application"},
		{"", ""},
		{"Input View", ""},
		{"  Enter", "Confirm and continue"},
		{"  Esc", "Back to menu"},
		{"", ""},
		{"Result View", ""},
		{"  c", "Copy results to clipboard"},
		{"  Enter/Esc", "Back to menu"},
		{"", ""},
		{"Global", ""},
		{"  ?", "Toggle this help"},
		{"  Ctrl+C", "Force quit"},
	}

	for _, line := range helpContent {
		if line[0] == "" {
			b.WriteString("\n")
			continue
		}

		if line[1] == "" {
			// Section header
			section := lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				Render(line[0])
			b.WriteString(section)
		} else {
			// Key binding
			key := lipgloss.NewStyle().
				Foreground(ColorAccent).
				Render(line[0])
			desc := lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Render(line[1])
			b.WriteString(fmt.Sprintf("  %s - %s", key, desc))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	footer := m.footerStyle.Render("Press ? to close help")
	b.WriteString(footer)

	return b.String()
}

// renderHeader returns the ASCII art header (responsive to terminal height)
func (m model) renderHeader() string {
	// If terminal is too small, show compact header
	if m.height < 30 {
		return lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Render("THE WATCHMAN - Network Security Scanner")
	}

	// Full ASCII art header
	art := `
╔════════════════════════════════════════════════════════════════╗
║                                                                ║
║            ████████╗██╗  ██╗███████╗                           ║
║            ╚══██╔══╝██║  ██║██╔════╝                           ║
║               ██║   ███████║█████╗                             ║
║               ██║   ██╔══██║██╔══╝                             ║
║               ██║   ██║  ██║███████╗                           ║
║               ╚═╝   ╚═╝  ╚═╝╚══════╝                           ║
║                                                                ║
║        ██╗    ██╗ █████╗ ████████╗ ██████╗██╗  ██╗             ║
║        ██║    ██║██╔══██╗╚══██╔══╝██╔════╝██║  ██║             ║
║        ██║ █╗ ██║███████║   ██║   ██║     ███████║             ║
║        ██║███╗██║██╔══██║   ██║   ██║     ██╔══██║             ║
║        ╚███╔███╔╝██║  ██║   ██║   ╚██████╗██║  ██║             ║
║         ╚══╝╚══╝ ╚═╝  ╚═╝   ╚═╝    ╚═════╝╚═╝  ╚═╝             ║
║                 M A N                                          ║
║                                                                ║
║               Network Security Scanner                         ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝`
	return lipgloss.NewStyle().Foreground(ColorPrimary).Render(art)
}
