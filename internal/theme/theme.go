package theme

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ActiveBorderColor = lipgloss.Color("212")
	InactiveBorderColor = lipgloss.Color("240")
	DirColor = lipgloss.Color("33")
	FileColor = lipgloss.Color("252")
	DimColor = lipgloss.Color("243")
	CursorStyle = lipgloss.NewStyle().
			Reverse(true).
			Bold(true)
)


func EntryStyle(isDir bool) lipgloss.Style {
	if isDir {
		return lipgloss.NewStyle().Foreground(DirColor).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(FileColor)
}


func PaneStyle(width, height int, active bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(0, 1)

	if active {
		return style.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ActiveBorderColor)
	}

	return style.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(InactiveBorderColor)
}
