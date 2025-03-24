// Package tui предоставляет инструменты для создания и управления текстовым пользовательским интерфейсом (TUI) приложения GophKeeper.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
	"github.com/romanp1989/gophkeeper/internal/client/storage"
)

// NavigationCallback определяет тип функции обратного вызова для навигационных сообщений.
type NavigationCallback func(args ...any) tea.Cmd

// NavigationMsg представляет сообщение, используемое для навигации и конфигурации экранов в TUI.
type NavigationMsg struct {
	Callback     NavigationCallback
	Client       grpc.ClientGRPCInterface
	DisableFocus bool
	Page         Page
	Position     Position
	Screen       Screen
	Secret       *domain.Secret
	Storage      storage.Storage
}

// NewNavigationMsg создаёт новое навигационное сообщение с указанными настройками экрана и опциональными параметрами.
func NewNavigationMsg(screen Screen, opts ...NavigateOption) NavigationMsg {
	msg := NavigationMsg{Page: Page{Screen: screen}}
	for _, fn := range opts {
		fn(&msg)
	}
	return msg
}

// NavigateOption определяет тип функции, которая может модифицировать NavigationMsg для добавления дополнительных параметров.
type NavigateOption func(msg *NavigationMsg)

// WithCallback определяет опцию навигации для установки пользовательского callback.
func WithCallback(c NavigationCallback) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Callback = c
	}
}

// WithClient определяет опцию навигации для установки клиента gRPC.
func WithClient(c grpc.ClientGRPCInterface) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Client = c
	}
}

// WithPosition определяет опцию навигации для установки позиции элемента.
func WithPosition(position Position) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Position = position
	}
}

// WithStorage определяет опцию навигации для установки объекта хранилища.
func WithStorage(store storage.Storage) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Storage = store
	}
}

// WithSecret определяет опцию навигации для установки связанного секрета.
func WithSecret(sec *domain.Secret) NavigateOption {
	return func(msg *NavigationMsg) {
		msg.Secret = sec
	}
}
