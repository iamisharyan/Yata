package tui

import (
	
	"path/filepath"
	"term2/internal/fs"

	tea "github.com/charmbracelet/bubbletea"
)

// type editorFinishedMsg struct{
// 	err error
// }
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

	// case editorFinishedMsg:
	// 	if msg.err!=nil{
	// 		m.err=msg.err
	// 		return m,nil
	// 	}
	case tea.MouseMsg:
		 
		switch msg.Button{
		case tea.MouseButtonLeft:
			if msg.Action==tea.MouseActionPress{
				crow:=msg.Y -m.fileListY

				if crow>=0 && crow<len(m.current.Entries) {
					m.current.Cursor=crow
					
				}

			}
		case tea.MouseButtonRight:
			if msg.Action==tea.MouseActionPress{
				crow:=msg.Y-m.fileListY

				if crow>=0 && crow<len(m.current.Entries){
					m.current.Cursor=crow
				}
			}
		}
		return m ,nil

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
		case "f":
			if m.fullscreen{
				m.fullscreen=false
				return m,nil
			}
			if len(m.current.Entries)>0&&m.current.Cursor<len(m.current.Entries){
				sel:=m.current.Entries[m.current.Cursor]
				if !sel.IsDir{
					m.fullscreen=true
				}

			}
			return m,nil
		case "esc":
			if m.fullscreen{
				m.fullscreen=false
				return m,nil
			}
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

		case "n":
			selected:=m.current.Entries[m.current.Cursor]
			if selected.IsDir{
				return m,nil
			}
			path:=filepath.Join(m.cwd,selected.Name)
			if imagefile(selected.Name){
				return m,openimage(path)
			}
			if isdocfile(selected.Name){
				return m,docfile(path)
			}
			// if ispdffile(selected.Name){
			// 	return m,pdffile(path)
			// }
			return m,openinNvim(path)

			


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
