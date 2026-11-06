package icons

import (
	

"github.com/charmbracelet/lipgloss"
	"github.com/epilande/go-devicons"
)
func ForPath(path string)string{
	style:=devicons.IconForPath(path)
	iconstyle:=lipgloss.NewStyle()
	if style.Color!=""{
		iconstyle=iconstyle.Foreground(lipgloss.Color(style.Color))
	}
	return iconstyle.Render(style.Icon)+" "

}