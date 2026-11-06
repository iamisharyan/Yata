package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"term2/internal/fs"
	"term2/internal/preview"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
)
type Style struct{
	Icon string
	Color string
}

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
	fileListY   int 
	fullscreen  bool
	previewX viewport.Model
	fuzzyactive  bool
	fuzzyinput   textinput.Model
	fuzzyall     []fs.FuzzyItem
	fuzzyfiltered []fs.FuzzyItem
	fuzzycursor   int
	
}


func InitialModel() Model {
	cwd, _ := os.Getwd()
	vp:=viewport.New(0,0)
	ti:=textinput.New()
	ti.Placeholder="search files and dirs"
	ti.Prompt=" \ued23 "
	m := Model{
		cwd:        cwd,
		showHidden: false,
		previewX: vp,
		fuzzyinput: ti,
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
	m.fileListY=3

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

func openinNvim(path string)tea.Cmd{
	editor:=os.Getenv("EDITOR")

	if editor==""{
		editor="nvim"
	}
	c:=exec.Command(editor,path)
	return tea.ExecProcess(c,func(err error)tea.Msg{
		return editorFinishingmsg{err:err}
	})
}
func imagefile(name string)bool{
	ext:=strings.ToLower(filepath.Ext(name))
	switch ext{
		case ".png",".jpg",".jpeg",".gif",".bmp",
		".webp",".tiff",".ico":
		return true
	}
	return false
}
func openimage(path string) tea.Cmd{
	cmd:=exec.Command(
		"cmd",
		"/c",
		"start",
		"",
		path,

	)
	return tea.ExecProcess(cmd,func(err error)tea.Msg{
		return editorFinishingmsg{err:err}
	})
}
func isdocfile(name string)bool{
	ext:=strings.ToLower(filepath.Ext(name))
	switch ext{
		case ".doc",".docx",".docm",".dot",".dotx",
		".rtf",".odt",".dotm",".pdf":
		return true
	}
	return false
}
func  docfile(path string) tea.Cmd{
	cmd:=exec.Command(
		"cmd",
		"/c",
		"start",
		"",
		path,

	)
	return tea.ExecProcess(cmd,func(err error)tea.Msg{
		return editorFinishingmsg{err:err}
	})
}

// func ispdffile(name string)bool{
// 	ext:=strings.ToLower(filepath.Ext(name))
// 	switch ext{
// 		case ".pdf":
// 		return true
// 	}
// 	return false
// }
// func pdffile(path string) tea.Cmd{
	
		
// c:=os.Getenv("PDFOPENER")
// if c==""{
// 	c="brave"
// }
// cmd:=exec.Command(c,path)
// 	return tea.ExecProcess(cmd,func(err error)tea.Msg{
// 		return editorFinishingmsg{err:err}
// 	})
// }


type editorFinishingmsg struct{
	err error
}
func (m *Model) resizepreview() {
    var width, height int

    if m.fullscreen {
        width = m.width - 2
        height = m.height - 2
    } else {
        colWidth := m.width / 3
        paneHeight := m.height - 4

        width = colWidth - 4
        height = paneHeight - 2
    }

    m.previewX.Width=max(width, 1)
    m.previewX.Height=max(height, 1)
}