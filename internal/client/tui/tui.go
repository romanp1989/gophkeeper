// Package tui предоставляет утилиты и компоненты для создания текстового пользовательского интерфейса в приложении.
package tui

import "github.com/charmbracelet/bubbletea"

// CmdHandler создаёт команду для Bubble Tea, которая возвращает предоставленное сообщение.
func CmdHandler(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
