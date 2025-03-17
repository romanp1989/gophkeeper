// Package secrets содержит компоненты для отображения и взаимодействия со списком секретов в TUI.
package secrets

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"io"
	"strings"
)

type secretItem struct {
	id   int
	name string
}

// FilterValue возвращает значение для фильтрации списка.
func (i secretItem) FilterValue() string { return "" }

type secretItemDelegate struct{}

// Height возвращает высоту элемента списка.
func (d secretItemDelegate) Height() int { return 1 }

// Spacing возвращает пробел между элементами списка.
func (d secretItemDelegate) Spacing() int { return 0 }

// Update обрабатывает сообщения для обновления модели списка.
func (d secretItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render отображает элемент списка.
func (d secretItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		itemStyle         = styles.Regular.PaddingLeft(4)
		itemSelectedStyle = styles.Regular.PaddingLeft(2)
	)

	i, ok := listItem.(secretItem)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return itemSelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(i.name))
}
