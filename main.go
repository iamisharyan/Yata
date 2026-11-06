package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"term2/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen(),tea.WithMouseCellMotion(),)
	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
