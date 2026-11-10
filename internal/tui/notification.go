
package tui

import (
    "time"

    tea "github.com/charmbracelet/bubbletea"
)

type notificationMsg struct {
    id uint64
}

func (m *Model) Notify(message string) tea.Cmd {
    m.notifID++
    id := m.notifID
    m.statusmessage = message

    return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
        return notificationMsg{id: id}
    })
}
