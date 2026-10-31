package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)
type entry struct{
	name string
	isDir bool
}
type model struct{
	cwd string
	entries []entry
	cursor int
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
	case tea.KeyMsg:
		switch msg.String(){
		case "q","ctrl+c":
			return m,tea.Quit
		case "down","j":
			if m.cursor<len(m.entries)-1{
				m.cursor++
			}
		case "up","k":
			if m.cursor>0{
				m.cursor--
			}
		}
	}
	return m,nil
}
func (m model) View() string{
var b strings.Builder
b.WriteString(m.cwd+"\n\n")
for i,e:=range m.entries{
	cursor := " "
	if i==m.cursor{
		cursor=">"
	}
	name :=e.name
	if e.isDir{
		name+="/"
	}
	b.WriteString(fmt.Sprintf("%s%s\n",cursor,name))
}
	return b.String()
}
func main(){
	p:=tea.NewProgram(initialModel())
	if _,err:=p.Run();err!=nil{
		fmt.Println("error running program:",err)
		os.Exit(1)
	}
}
