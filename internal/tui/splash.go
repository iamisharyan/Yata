package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	logoColor   = lipgloss.Color("#C99AFF") // Primary lavender accent
	mutedColor  = lipgloss.Color("#70518D") // Dim secondary accent
	subtleColor = lipgloss.Color("#3D2D54") // Dim border tone for inactive pills
	textColor   = lipgloss.Color("#E2D9F3") // Clean foreground text
)

var yataFont = []string{
	"10001 01110 11111 01110",
	"01010 10001 00100 10001",
	"00100 10001 00100 10001",
	"00100 11111 00100 11111",
	"00100 10001 00100 10001",
	"00100 10001 00100 10001",
	"00100 10001 00100 10001",
}

func renderYataLogo(width int) string {
	pixelWidth := width * 6 / 10 / 23
	if pixelWidth < 1 {
		pixelWidth = 1
	} else if pixelWidth > 4 {
		pixelWidth = 4
	}

	const rows = 7
	const cols = 23

	var foreground [rows][cols]bool

	for y, line := range yataFont {
		parts := strings.Fields(line)
		x := 0

		for _, part := range parts {
			for _, pixel := range part {
				if pixel == '1' {
					foreground[y][x] = true
				}
				x++
			}
			x++
		}
	}

	var b strings.Builder
	fgStyle := lipgloss.NewStyle().Foreground(logoColor)
	solidBlock := strings.Repeat("█", pixelWidth)
	emptyBlock := strings.Repeat(" ", pixelWidth)

	for y := 0; y < rows; y++ {
		for scaleY := 0; scaleY < 2; scaleY++ {
			for x := 0; x < cols; x++ {
				if foreground[y][x] {
					b.WriteString(fgStyle.Render(solidBlock))
				} else {
					b.WriteString(emptyBlock)
				}
			}
			b.WriteString("\n")
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func renderSplash(width, height, cursor int) string {
	if width < 1 || height < 1 {
		return ""
	}

	logo := renderYataLogo(width)

	kanjiBadge := lipgloss.NewStyle().
		Foreground(mutedColor).
		Bold(true).
		Render("―  八  咫  ―")

	tagline := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Render("High-Performance Terminal File Manager")

	// Base pill button layout using Lipgloss RoundedBorder with no background fills
	baseBtn := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(22).
		Align(lipgloss.Center)

	activeBtn := baseBtn.
		BorderForeground(logoColor).
		Foreground(logoColor).
		Bold(true)

	inactiveBtn := baseBtn.
		BorderForeground(subtleColor).
		Foreground(mutedColor)

	var openBtn, exitBtn string
	if cursor == 0 {
		openBtn = activeBtn.Render("▶  OPEN")
		exitBtn = inactiveBtn.Render("EXIT")
	} else {
		openBtn = inactiveBtn.Render("OPEN")
		exitBtn = activeBtn.Render("▶  EXIT")
	}

	footerHelp := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("↑/k • ↓/j navigate   •   enter select   •   q quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		kanjiBadge,
		"",
		logo,
		"",
		tagline,
		"",
		openBtn,
		exitBtn,
		"",
		footerHelp,
	)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}