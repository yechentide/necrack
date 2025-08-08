package styles

import "github.com/charmbracelet/lipgloss"

// color definitions
const (
	colorGreen   = "#04B575"
	colorBlue    = "#00ADD8"
	colorYellow  = "#FFD700"
	colorOrange  = "#FF6B35"
	colorRed     = "#FF4444"
	colorSuccess = "#00FF00"
	colorMuted   = "#888888"
)

// Common styles
var (
	HeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGreen)).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSuccess)).
		Bold(true)

	PathStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBlue))

	URLStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBlue)).
		Underline(true)

	InfoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorYellow))

	KeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorYellow))

	ErrorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorRed)).
		Bold(true)

	MutedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorMuted))
)

// Command-specific header styles
var (
	DecodeHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGreen)).
		Bold(true)

	EncodeHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorOrange)).
		Bold(true)

	ServerHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGreen)).
		Bold(true)
)