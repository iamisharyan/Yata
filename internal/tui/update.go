package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"term2/internal/fs"
)


func (m Model) Init() tea.Cmd {
	return nil
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m, cmd = m.RefreshPanes()
		return m, cmd

	case previewLoadedMsg:
		if msg.path == m.SelectedPath() {
			if msg.err != nil {
				m.previewText = "[error reading file]"
			} else {
				m.previewText = msg.content
			}
			m.previewPath = msg.path
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		
		if m.pendingG {
			m.pendingG = false
			if key == "g" {
				m.current.Cursor = 0
				m, cmd = m.RefreshPanes()
				return m, cmd
			}
			
		}

		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "down", "j":
			if m.current.Cursor < len(m.current.Entries)-1 {
				m.current.Cursor++
				m, cmd = m.RefreshPanes()
			}

		case "up", "k":
			if m.current.Cursor > 0 {
				m.current.Cursor--
				m, cmd = m.RefreshPanes()
			}

		case "g":
			m.pendingG = true
			return m, nil

		case "G":
			if len(m.current.Entries) > 0 {
				m.current.Cursor = len(m.current.Entries) - 1
				m, cmd = m.RefreshPanes()
			}

		case "h", "left":
			m, cmd = m.GoToParent()

		case "enter", "l", "right":
			m, cmd = m.EnterSelected()

		case "ctrl+d", "w":
			if len(m.current.Entries) > 0 {
				m.current.Cursor = min(m.current.Cursor+10, len(m.current.Entries)-1)
				m, cmd = m.RefreshPanes()
			}

		case "ctrl+u", "o":
			m.current.Cursor = max(m.current.Cursor-10, 0)
			m, cmd = m.RefreshPanes()

		case ".":
			m.showHidden = !m.showHidden
			m.current.Entries = m.filterEntries(fs.ReadDir(m.cwd))
			if m.current.Cursor >= len(m.current.Entries) {
				m.current.Cursor = max(len(m.current.Entries)-1, 0)
			}
			m, cmd = m.RefreshPanes()
		}
	}

	return m, cmd
}
