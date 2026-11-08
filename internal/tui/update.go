package tui

import (
	"path/filepath"
	"term2/internal/fs"
	"term2/internal/ops"

	"github.com/charmbracelet/bubbles/textinput"
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
		m.resizepreview()
		
		

		m, cmd = m.RefreshPanes()
		return m, cmd

	// case editorFinishedMsg:
	// 	if msg.err!=nil{
	// 		m.err=msg.err
	// 		return m,nil
	// 	}
	case tea.MouseMsg:
		if m.fullscreen{
				var cmd tea.Cmd
				m.previewX,cmd=m.previewX.Update(msg)
				return m,cmd
			}
		 
		switch msg.Button{
		case tea.MouseButtonWheelDown:
			m.previewX.ScrollDown(3)
			
		case tea.MouseButtonWheelUp:
			m.previewX.ScrollUp(3)
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
				m.previewText = m.renderpreview(msg.path)
			}
			m.previewPath = msg.path
			m.previewX.SetContent(m.previewText)
			m.previewX.GotoTop()
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()
if m.renaming{
	switch msg.String(){
	case "enter":
		selected:=m.current.Entries[m.current.Cursor]
		src:=filepath.Join(m.cwd,selected.Name)

		_,err:=ops.Rename(src,m.renameInput)
		if err!=nil{
			return m,nil
		}else{
			m.RefreshPanes()
		}
		m.renaming=false
		m.renameInput=""
		return m,nil

	case "esc":
		m.renaming=false
		m.renameInput=""
		return m,nil
	case "backspace":
		runes:=[]rune(m.renameInput)
		if len(runes)>0{
			m.renameInput=string(runes[:len(runes)-1])
		}
		return m,nil
	default:
		key:=msg.String()
		if len([]rune(key))==1{
			m.renameInput+=key
		}
		return m,nil

	}
	
}

		if m.fuzzyactive{
			switch key{
			case "esc":
				m.fuzzyactive=false
				m.fuzzyinput.Blur()
				return m,nil
			case "down","ctrl+l":
				if m.fuzzycursor<len(m.fuzzyfiltered)-1{
					m.fuzzycursor++
				}
				return m,nil
			case "up","ctrl+k":
				if m.fuzzycursor>0{
					m.fuzzycursor--
				}
				return m,nil
			case "enter","=":
				if len(m.fuzzyfiltered)==0||m.fuzzycursor>=len(m.fuzzyfiltered){
					
				return m,nil
			}
					target:=m.fuzzyfiltered[m.fuzzycursor]
					fullpath:=filepath.Join(m.cwd,target.Relpath)
					m.fuzzyactive=false
					m.fuzzyinput.Blur()

					if target.IsDir{
						m.cwd=fullpath
						m.current=Pane{
							Path:  fullpath,
							Entries: m.filterEntries(fs.ReadDir(fullpath)),
							Cursor: 0,
						}
						return m.RefreshPanes()
					}
					td:=filepath.Dir(fullpath)
					tf:=filepath.Base(fullpath)
					m.cwd=td
					m.current=Pane{
						Path:td,
						Entries: m.filterEntries(fs.ReadDir(td)),
						Cursor: 0,
					}

						for i,e:=range m.current.Entries{
							if e.Name==tf{
								m.current.Cursor=i
								break
							}
						}
						var cmd tea.Cmd
						m,cmd=m.RefreshPanes()
						m.fullscreen=true
						m.resizepreview()
						return m,tea.Batch(cmd,loadPreviewCmd(fullpath))
						
					}
					
				
			
			var tiCmd tea.Cmd
			m.fuzzyinput,tiCmd=m.fuzzyinput.Update(msg)
			m.fuzzyfiltered=fs.FuzzySearch(m.fuzzyinput.Value(),m.fuzzyall)
			m.fuzzycursor=0
			return m,tiCmd
		}
		switch key{
		case "/":
			m.fuzzyactive=true
			m.fuzzyinput.Reset()
			m.fuzzyinput.Focus()
			m.fuzzyall=fs.CollectFiles(m.cwd,m.showHidden)
			m.fuzzyfiltered=m.fuzzyall
			m.fuzzycursor=0
			return m,textinput.Blink
		}
		

		if m.fullscreen {
    switch key {
    case "alt+1":
        m.previewX.ScrollDown(1)
        return m, nil

    case "alt+2":
        m.previewX.ScrollUp(1)
        return m, nil

    case "0":
        m.previewX.HalfPageDown()
        return m, nil

    case "9":
        m.previewX.HalfPageUp()
        return m, nil
    }
}
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
				
			}
			if len(m.current.Entries)>0&&m.current.Cursor<len(m.current.Entries){
				sel:=m.current.Entries[m.current.Cursor]
				if !sel.IsDir{
					m.fullscreen=true
				}

			}
			m.resizepreview()
			return m,nil
		case "esc":
			if m.fullscreen{
				m.fullscreen=false
				m.resizepreview()
				return m,nil
			}
			//file operations
		case "c":
			if len(m.current.Entries)==0{
				return m,nil
			}
			selected:=m.current.Entries[m.current.Cursor]
			path:=filepath.Join(m.cwd,selected.Name)
			m.clipboard.Set(path,ops.Copy)
		case "x":
			if len(m.current.Entries)==0{
				return m,nil
			}
			selected:=m.current.Entries[m.current.Cursor]
			path:=filepath.Join(m.cwd,selected.Name)
			m.clipboard.Set(path,ops.Cut)
		case "p":
			if m.clipboard.Source==""{
				return m,nil
			}
			_,err:=ops.Paste(
				m.clipboard.Source,
				m.cwd,
				m.clipboard.Mode,
			)
			if err!=nil{
				return m,nil
			}
			if m.clipboard.Mode==ops.Cut{
				m.clipboard.Clear()
			}
			m.RefreshPanes()
		case "r":
			if len(m.current.Entries)==0{
				return m,nil
			}
			selected:=m.current.Entries[m.current.Cursor]
			m.renaming=true
			m.renameInput=selected.Name


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
