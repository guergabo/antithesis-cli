package cli

import "github.com/charmbracelet/lipgloss"

var (
	// Core colors.

	defaultTextColor = lipgloss.AdaptiveColor{
		Light: "#000000",
		Dark:  "#ffffff",
	}
	primaryColor = lipgloss.ANSIColor(39)  // Bright blue
	accentColor  = lipgloss.ANSIColor(141) // Light purple
	mutedColor   = lipgloss.ANSIColor(242) // Gray

	// Status colors.

	successColor = lipgloss.ANSIColor(78)  // Bright green
	warningColor = lipgloss.ANSIColor(208) // Bright orange
	errorColor   = lipgloss.ANSIColor(196) // Bright red

	// Pre-defined core styles.

	BaseStyle = lipgloss.NewStyle().
			Foreground(defaultTextColor)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			PaddingBottom(1)

	ValueStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	// Pre-defined status styles.

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)
)
