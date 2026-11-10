package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"term2/internal/icons"
	"term2/internal/preview"
	"term2/internal/theme"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)
var pointer string
func (m Model) View() string {
	
	if m.fuzzyactive{
		return m.renderFuzzyModal()
	}
	if m.fullscreen{
		return lipgloss.NewStyle().
        Width(m.width).
        Height(m.height).
        Padding(1).
        Render(m.previewX.View())
	
		
	}
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}
	if m.screen == SplashScreen {
    return renderSplash(
        m.width,
        m.height,
        m.menuCursor,
    )
}

	colWidth:=m.width/3
	if colWidth<10 {
		colWidth=10
	}

	paneHeight:=m.height-4
	if paneHeight<5{
		paneHeight= 5
	}

	parentView := renderDirPane(m.parent, colWidth, paneHeight, false)
	currentView := renderDirPane(m.current, colWidth, paneHeight, true)

	var previewView string
	if len(m.current.Entries) > 0 && m.current.Cursor < len(m.current.Entries) {
		sel := m.current.Entries[m.current.Cursor]
		if sel.IsDir {
			previewView = renderDirPane(m.preview, colWidth, paneHeight, false)
		} else {
			// path:=filepath.Join(m.cwd,sel.Name)
			// content:=m.renderpreview(path)
			// previewView = renderTextPreview(content, colWidth, paneHeight)
			previewView=theme.PaneStyle(
				colWidth,
				paneHeight,
				false,
			).Render(m.previewX.View())
		}
	} else {
		previewView = renderTextPreview("", colWidth, paneHeight)
	}

	columns := lipgloss.JoinHorizontal(lipgloss.Top, parentView, currentView, previewView)

	permissions:="N/A"
	selectedname:="None"

	if len(m.current.Entries)>0&&m.current.Cursor<len(m.current.Entries){
		sel:=m.current.Entries[m.current.Cursor]
		selectedname=sel.Name
		path:=filepath.Join(m.cwd,sel.Name)
		permissions=getpermissions(path)
	}
	statusBar := lipgloss.NewStyle().
		Foreground(theme.DimColor).
		Padding(0, 1).
		Render("cwd:"+m.cwd,"•",selectedname,"   •",permissions,"         •[q]Quit")
if m.renaming {
    title := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#00D7AF")).
        Render("Rename File")

    input := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Background(lipgloss.Color("#262626")).
        Padding(0, 1).
        Width(40).
        Render(m.renameInput + "█")

    help := lipgloss.NewStyle().
        Foreground(theme.DimColor).
        Render("ENTER: Confirm    ESC: Cancel")

    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        "",
        input,
        "",
        help,
    )

    modal := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(theme.ActiveBorderColor).
        Padding(1, 2).
        Width(46).
        Render(content)

    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center,
        modal,
    )
}
	mainView := lipgloss.JoinVertical(
    lipgloss.Left,
    columns,
    statusBar,
)

if m.statusmessage != "" {
    toast := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Background(lipgloss.Color("#168B72")).
        Bold(true).
        Padding(0, 2).
       Border(lipgloss.RoundedBorder()).
        Render(m.statusmessage)

    return lipgloss.JoinVertical(
        lipgloss.Left,
        mainView,
        lipgloss.PlaceHorizontal(
            m.width,
            lipgloss.Center,
            toast,
        ),
    )
}

return mainView

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
indicatorwidth:=ansi.StringWidth(theme.CursorIndicator)
	var b strings.Builder
	for i := offset; i < len(p.Entries) && i < offset+visible; i++ {
		e := p.Entries[i]

		
		if i==p.Cursor&&active{
			pointer=theme.CursorIndicatorStyle.Render(theme.CursorIndicator)
		}else{
			pointer=theme.CursorIndicatorEmpty
		}
		fullpath:=filepath.Join(p.Path,e.Name)
		iconprefix:=icons.ForPath(fullpath)
		iconwidth:=ansi.StringWidth(iconprefix)

		namewidth:=innerWidth-indicatorwidth-iconwidth
		if namewidth<1{
			namewidth=1
		}
		
		name := e.Name
		if e.IsDir {
			name += "/"
		}

		truncatedName:=ansi.Truncate(name,namewidth,"...")
		renderedname:=theme.EntryStyle(e.IsDir).Render(truncatedName)
		line:=iconprefix+renderedname
		if i == p.Cursor&& active {
			pointer =theme.CursorIndicatorStyle.Render(theme.CursorIndicator)
			linewidth:=ansi.StringWidth(line)
			rem:=(innerWidth-indicatorwidth)-linewidth
			if rem>0{
				line+=strings.Repeat(" ",rem)
			}
			
			
			line=theme.CursorIndicatorStyle.Render(line)
		} else{
			pointer=strings.Repeat(" ",indicatorwidth)
		}

		b.WriteString(pointer+line + "\n")
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
		b.WriteString(ansi.Truncate(l, innerWidth,"") + "\n")
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
func(m Model)renderpreview(path string)string{

	if preview.IsBinary(path){
		return "binary file"
	}
e:= strings.ToLower(filepath.Ext(path))

switch e{


	case ".go",".c",".cpp",".h",".rs",".ts",".js",
	".py",".html",".css",".json",".yml",".toml",".sh",
	".sql",".java":
	content,err:=preview.Highlightfile(path)
	if err!=nil{
		return "preview error:"+err.Error()
	}
	return content

default:
	content ,err:=os.ReadFile(path)

	if err!=nil{
		return "preview error"+err.Error()
	}
	return string(content)
}
}
func getpermissions(path string)string{
	info,err:=os.Lstat(path)
	if err!=nil{
		return "N/A"
	}
	mode:=info.Mode().Perm()
	return fmt.Sprintf("%s(%04o)",mode.String(),mode)
}
func (m Model)renderFuzzyModal()string{
	mw:=max(m.width-10,40)
	mh:=max(m.height-6,10)
	var b strings.Builder
	b.WriteString(m.fuzzyinput.View()+"\n\n")
	mv:=mh-5
	s:=0
	if m.fuzzycursor>=mv{
		s=m.fuzzycursor-mv+1
	}
	for i:=s;i<len(m.fuzzyfiltered)&&i<s+mv;i++{
		item:=m.fuzzyfiltered[i]
		l:=item.Relpath
		if item.IsDir{
			l+="/"
		}
		l=ansi.Truncate(l,mw-8,"...")

		if i==m.fuzzycursor{
			l=theme.CursorIndicatorStyle.Render("> ")+theme.CursorIndicatorStyle.Render(l)
		}else{
			l=" "+l
		}
		b.WriteString(l+"\n")
	}
	mb:=lipgloss.NewStyle().Width(mw).Height(mh).Border(lipgloss.RoundedBorder()).BorderForeground(theme.ActiveBorderColor).Padding(1,2).Render(b.String())

	return lipgloss.Place(m.width,m.height,lipgloss.Center,lipgloss.Center,mb)
}