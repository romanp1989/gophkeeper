// Package tui предоставляет утилиты для управления текстовым интерфейсом пользователя и навигацией внутри приложения.
package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	// ErrorMsg представляет сообщение об ошибке.
	ErrorMsg error

	// InfoMsg представляет информационное сообщение.
	InfoMsg string
)

// ReportInfo создает команду, которая отправляет информационное сообщение в систему Bubble Tea.
func ReportInfo(msg string, args ...any) tea.Cmd {
	return CmdHandler(InfoMsg(fmt.Sprintf(msg, args...)))
}

// ReportError создает команду, которая отправляет сообщение об ошибке в систему Bubble Tea.
func ReportError(err error) tea.Cmd {
	return CmdHandler(ErrorMsg(err))
}

// NavigateTo создает команду для навигации на указанный экран.
func NavigateTo(screen Screen, opts ...NavigateOption) tea.Cmd {
	return CmdHandler(NewNavigationMsg(screen, opts...))
}

// SetBodyPane создает команду для установки основного содержимого экрана.
func SetBodyPane(screen Screen, opts ...NavigateOption) tea.Cmd {
	opts = append(opts, WithPosition(BodyPane))
	return NavigateTo(screen, opts...)
}
