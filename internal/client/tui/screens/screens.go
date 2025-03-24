// Package screens предоставляет функции для рендеринга содержимого экранов в TUI приложении GophKeeper.
package screens

import (
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"strings"
)

// RenderContent форматирует заголовок и содержимое экрана с применением стилей.
func RenderContent(header, content string) string {
	var b strings.Builder

	b.WriteString(styles.HeaderStyle.Render(header))
	b.WriteString("\n\n")
	b.WriteString(content)

	return styles.ContentPaddedStyle.Render(b.String())
}
