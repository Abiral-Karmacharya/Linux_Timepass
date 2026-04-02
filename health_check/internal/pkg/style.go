package tui

import "github.com/charmbracelet/lipgloss"

// Styles is a struct to hold our "CSS" rules
type UIStyles struct {
    Title    lipgloss.Style
    Critical  lipgloss.Style
    Warning   lipgloss.Style
    Normal    lipgloss.Style
    Info lipgloss.Style
    Error lipgloss.Style
}

// DefaultStyles acts like your "External Stylesheet"
func DefaultStyles() UIStyles {
    return UIStyles{
        Title: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FFFFFF")).
            Background(lipgloss.Color("#6200EE")).
            Padding(0, 2).
            Bold(true),

        Critical: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FF0000")).
            Blink(true).
            Bold(true),

        Warning: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FFA500")),

        Normal: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#00FF00")),
            
        Info: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            Padding(0, 2).
            BorderForeground(lipgloss.Color("62")),
        
        Error: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            Padding(0, 2).
            BorderForeground(lipgloss.Color("#FF0000")).
            Foreground(lipgloss.Color("#F76F25")),
    }
}