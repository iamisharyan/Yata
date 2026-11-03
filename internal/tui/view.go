package tui

import (
	
	"strings"

	"github.com/charmbracelet/lipgloss"
	"term2/internal/theme"
)

func (m Model) View() string {
	
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	colWidth := m.width / 3
	if colWidth < 10 {
		colWidth = 10
	}

	paneHeight := m.height - 4
	if paneHeight < 5 {
		paneHeight = 5
	}

	parentView := renderDirPane(m.parent, colWidth, paneHeight, false)
	currentView := renderDirPane(m.current, colWidth, paneHeight, true)

	var previewView string
	if len(m.current.Entries) > 0 && m.current.Cursor < len(m.current.Entries) {
		sel := m.current.Entries[m.current.Cursor]
		if sel.IsDir {
			previewView = renderDirPane(m.preview, colWidth, paneHeight, false)
		} else {
			previewView = renderTextPreview(m.previewText, colWidth, paneHeight)
		}
	} else {
		previewView = renderTextPreview("", colWidth, paneHeight)
	}

	columns := lipgloss.JoinHorizontal(lipgloss.Top, parentView, currentView, previewView)

	
	statusBar := lipgloss.NewStyle().
		Foreground(theme.DimColor).
		Padding(0, 1).
		Render("CWD: " + m.cwd + "  •  [q] Quit  [h/j/k/l] Navigate  [.] Toggle Hidden")
		

	return lipgloss.JoinVertical(lipgloss.Left, columns, statusBar)

}

func renderDirPane(p Pane, width, height int, active bool) string {
	visible := height - 2
	if visible < 1 {
		visible = 1
	}

	innerWidth := width - 4
	if innerWidth < 1 {
		innerWidth = 1
	}

	
	offset := p.Offset
	if p.Cursor < offset {
		offset = p.Cursor
	}
	if p.Cursor >= offset+visible {
		offset = p.Cursor - visible + 1
	}

	var b strings.Builder
	for i := offset; i < len(p.Entries) && i < offset+visible; i++ {
		e := p.Entries[i]
		name := e.Name
		if e.IsDir {
			name += "/"
		}

		line := truncate(name, innerWidth)
		if i == p.Cursor {
			line = theme.CursorStyle.Render(line)
		} else {
			line = theme.EntryStyle(e.IsDir).Render(line)
		}

		b.WriteString(line + "\n")
	}

	return theme.PaneStyle(width, height, active).Render(b.String())
	
}

func renderTextPreview(content string, width, height int) string {
	visible := height - 2
	if visible < 1 {
		visible = 1
	}

	innerWidth := width - 4
	if innerWidth < 1 {
		innerWidth = 1
	}

	lines := strings.Split(content, "\n")
	var b strings.Builder
	count := 0
	for _, l := range lines {
		if count >= visible {
			break
		}
		b.WriteString(truncate(l, innerWidth) + "\n")
		count++
	}

	return theme.PaneStyle(width, height, false).Render(b.String())
}


func truncate(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) > maxRunes {
		return string(r[:maxRunes])
	}
	return s
}
