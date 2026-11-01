package tui

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"term2/internal/fs"
	"term2/internal/preview"
)


type Pane struct {
	Path    string
	Entries []fs.Entry
	Cursor  int
	Offset  int
}


type Model struct {
	cwd         string
	parent      Pane
	current     Pane
	preview     Pane
	previewText string
	previewPath string
	pendingG    bool
	showHidden  bool
	height      int
	width       int
}


func InitialModel() Model {
	cwd, _ := os.Getwd()
	m := Model{
		cwd:        cwd,
		showHidden: false,
	}
	m.current = Pane{
		Path:    cwd,
		Entries: m.filterEntries(fs.ReadDir(cwd)),
	}
	m, _ = m.RefreshPanes()
	return m
}


func (m Model) SelectedPath() string {
	if len(m.current.Entries) == 0 || m.current.Cursor >= len(m.current.Entries) {
		return ""
	}
	sel := m.current.Entries[m.current.Cursor]
	return filepath.Join(m.cwd, sel.Name)
}


func (m Model) filterEntries(entries []fs.Entry) []fs.Entry {
	if m.showHidden {
		return entries
	}
	var out []fs.Entry
	for _, e := range entries {
		if len(e.Name) > 0 && e.Name[0] != '.' {
			out = append(out, e)
		}
	}
	return out
}


func (m Model) RefreshPanes() (Model, tea.Cmd) {
	parentDir := filepath.Dir(m.cwd)
	m.parent = Pane{
		Path:    parentDir,
		Entries: m.filterEntries(fs.ReadDirsOnly(parentDir)),
	}
	for i, e := range m.parent.Entries {
		if e.Name == filepath.Base(m.cwd) {
			m.parent.Cursor = i
			break
		}
	}

	m.preview = Pane{}
	m.previewText = ""
	m.previewPath = ""

	if len(m.current.Entries) == 0 || m.current.Cursor >= len(m.current.Entries) {
		return m, nil
	}

	sel := m.current.Entries[m.current.Cursor]
	selPath := filepath.Join(m.cwd, sel.Name)

	if sel.IsDir {
		m.preview = Pane{
			Path:    selPath,
			Entries: m.filterEntries(fs.ReadDir(selPath)),
		}
		return m, nil
	}

	
	return m, loadPreviewCmd(selPath)
}


func (m Model) GoToParent() (Model, tea.Cmd) {
	parent := filepath.Dir(m.cwd)
	if parent == m.cwd {
		return m, nil
	}
	oldBase := filepath.Base(m.cwd)
	m.cwd = parent
	m.current = Pane{
		Path:    parent,
		Entries: m.filterEntries(fs.ReadDir(parent)),
	}

	
	m.current.Cursor = 0
	for i, e := range m.current.Entries {
		if e.Name == oldBase {
			m.current.Cursor = i
			break
		}
	}

	return m.RefreshPanes()
}


func (m Model) EnterSelected() (Model, tea.Cmd) {
	if len(m.current.Entries) == 0 {
		return m, nil
	}
	sel := m.current.Entries[m.current.Cursor]
	if sel.IsDir {
		newPath := filepath.Join(m.cwd, sel.Name)
		m.cwd = newPath
		m.current = Pane{
			Path:    newPath,
			Entries: m.filterEntries(fs.ReadDir(newPath)),
			Cursor:  0,
		}
		return m.RefreshPanes()
	}
	return m, nil
}

type previewLoadedMsg struct {
	path    string
	content string
	err     error
}

func loadPreviewCmd(path string) tea.Cmd {
	return func() tea.Msg {
		content, err := preview.TextPreview(path)
		return previewLoadedMsg{
			path:    path,
			content: content,
			err:     err,
		}
	}
}
