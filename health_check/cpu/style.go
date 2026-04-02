package cpu

import "github.com/charmbracelet/lipgloss"

// Styles is a struct to hold our "CSS" rules
type UIStyles struct {
    Header    lipgloss.Style
    Critical  lipgloss.Style
    Warning   lipgloss.Style
    Normal    lipgloss.Style
    Container lipgloss.Style
    Error lipgloss.Style
}

// DefaultStyles acts like your "External Stylesheet"
func DefaultStyles() UIStyles {
    return UIStyles{
        Header: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FFFFFF")).
            Background(lipgloss.Color("#6200EE")).
            Padding(1, 0).
            Bold(true),

        Critical: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FF0000")).
            Blink(true).
            Bold(true),

        Warning: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#FFA500")),

        Normal: lipgloss.NewStyle().
            Foreground(lipgloss.Color("#00FF00")),
            
        Container: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            Padding(1, 0).
            BorderForeground(lipgloss.Color("62")),
        
        Error: lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            Padding(1, 0).
            BorderForeground(lipgloss.Color("#FF0000")).
            Foreground(lipgloss.Color("#F76F25")),
    }
}