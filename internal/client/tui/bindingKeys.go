// Package tui предоставляет инструменты и утилиты для работы с текстовым пользовательским интерфейсом (TUI) в приложении.
package tui

import "github.com/charmbracelet/bubbles/key"

// GlobalKeys содержит глобальные горячие клавиши, которые применяются во всём текстовом интерфейсе пользователя.
var GlobalKeys = struct {
	// Quit Горячая клавиша для выхода из приложения
	Quit key.Binding
	// Help Горячая клавиша для вызова помощи
	Help key.Binding
}{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "help"),
	),
}
