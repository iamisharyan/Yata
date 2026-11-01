package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)
type entry struct{
	name string
	isDir bool
}
type model struct{
	cwd string
	entries []entry
	parent pane
	current pane
	preview pane
	cursor int
	pendingG bool
	height int
	width int
}
type pane struct{
	path string
	entries []entry
	cursor int
}
func (m model) refreshpanes()model{
	m.parent=pane{
		path:filepath.Dir(m.cwd),
		entries:readDir(filepath.Dir(m.cwd)),

	}
	for i,e:=range m.parent.entries{
		if e.name==filepath.Base(m.cwd){
			m.parent.cursor=i
		}
	}
	if(len(m.current.entries)>0){
		sel:=m.current.entries[m.current.cursor]
		selPath:=filepath.Join(m.cwd,sel.name)
		if sel.isDir{
			m.preview=pane{path:selPath,entries: readDir(selPath)}
		}
	}
	return m

}
func readDir(path string)[]entry{
files,err:=os.ReadDir(path)
if err!=nil{
	return nil
}
var out []entry
for _,f:=range files{
	out =append(out,entry{name:f.Name(),isDir: f.IsDir()})
}
sort.Slice(out,func(i,j int)bool{
	if out[i].isDir!=out[j].isDir{
		return out[i].isDir
	}
	return strings.ToLower(out[i].name)<strings.ToLower(out[j].name)
})
return out
}
func initialModel() model{
	cwd,_:=os.Getwd()
	return model{
		cwd: cwd,
		entries:readDir(cwd),
	}
}

func (m model)Init() tea.Cmd{
return nil
}
func (m model) Update(msg tea.Msg)(tea.Model,tea.Cmd){
	switch msg:=msg.(type){
	case tea.WindowSizeMsg:
		m.width=msg.Width
		m.height=msg.Height
	case tea.KeyMsg:
		key:= msg.String()
		if m.pendingG{
			m.pendingG=false
			if key=="g"{
				m.cursor=0
				return m,nil
			}
		}
		switch key{
		case "q","ctrl+c":
			return m,tea.Quit
		case "j","down":
			if m.cursor<len(m.current.entries)-1{
				m.cursor++
			}
		case "k","up":
			if m.cursor>0{
				m.cursor--
			}
		case "g":
			m.pendingG=true
		case "G":
			m.cursor=len(m.current.entries)-1
		case "h","left":
			m=m.goToParent()
		case "l","right":
			m=m.enterSelected()
		case "ctrl+d","w":
			m.cursor=min(m.cursor+10,len(m.current.entries)-1)
		case "ctrl+u","o":
			m.cursor=max(m.cursor-10,0)
		}

	}
	return m.refreshpanes(),nil
}
func(m model)goToParent()model{
	parent:=filepath.Dir(m.cwd)
	if parent==m.cwd{
		return m
	}
	m.cwd=parent
	m.current.entries=readDir(parent)
	m.cursor=0
	return m
}
func (m model)enterSelected()model{
	if len(m.entries)==0{
		return m
	}
	sel:=m.entries[m.cursor]
	if sel.isDir{
		newpath:=filepath.Join(m.cwd,sel.name)
		m.cwd =newpath
		m.current.entries=readDir(newpath)
		m.cursor=0
	}
	return m

}
func (m model) View() string{
colwidth:=m.width/3
parentView:=renderPane(m.parent,colwidth,false)
currentView:=renderPane(m.current,colwidth,true)
previewView:=renderPane(m.preview,colwidth,false)

row:=lipgloss.JoinHorizontal(lipgloss.Top,parentView,currentView,previewView)
return row
}
func renderPane(p pane ,width int , active bool) string{
	style:=lipgloss.NewStyle().Width(width).Height(20).Padding(0,1)
	if active{
		style=style.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("212"))
	}else{
		style=style.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
	}

	var b strings.Builder
	for i,e:= range p.entries{
		line :=e.name
		if e.isDir{
			line +="/"

		}
		if i==p.cursor{
			line=lipgloss.NewStyle().Reverse(true).Render(line)
		}
		b.WriteString(line+"\n")
	}
	return style.Render(b.String())
}
func main(){
	p:=tea.NewProgram(initialModel(),tea.WithAltScreen())
	if _,err:=p.Run();err!=nil{
		fmt.Println("error running program:",err)
		os.Exit(1)
	}
}
