// Package components предоставляет компоненты для построения текстового пользовательского интерфейса (TUI).
// В этом пакете определены элементы управления вводом данных и кнопки, а также их логика взаимодействия и отображения.
package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"strings"
)

// InputGroup представляет группу текстовых полей ввода и кнопок.
// Содержит логику для управления фокусом ввода и обработки событий ввода.
type InputGroup struct {
	Inputs     []textinput.Model
	Buttons    []Button
	FocusIndex int
	totalPos   int
}

// Button представляет кнопку в пользовательском интерфейсе.
type Button struct {
	Title string
	Cmd   func() tea.Cmd
}

// NewInputGroup создает новый экземпляр InputGroup с заданными полями ввода и кнопками.
func NewInputGroup(inputs []textinput.Model, buttons []Button) InputGroup {
	for i, input := range inputs {
		if i == 0 {
			input.Focus()
			input.PromptStyle = styles.Focused
			input.TextStyle = styles.Focused
		}

		inputs[i] = input
	}

	return InputGroup{
		Inputs:   inputs,
		Buttons:  buttons,
		totalPos: len(inputs) + len(buttons) - 1,
	}
}

// Init инициализирует компонент, возвращая начальную команду для анимации курсора.
func (m InputGroup) Init() tea.Cmd {
	return textinput.Blink
}

// Update обновляет состояние компонента в ответ на события.
func (m InputGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "up", "down":
			s := msg.String()

			if s == "enter" {
				butIdx := m.FocusIndex - len(m.Inputs)
				if butIdx >= 0 {
					return m, m.Buttons[butIdx].Cmd()
				}
			}

			if s == "up" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > m.totalPos {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = m.totalPos
			}

			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					commands = append(commands, m.Inputs[i].Focus())
					m.Inputs[i].PromptStyle = styles.Focused
					m.Inputs[i].TextStyle = styles.Focused
					continue
				}

				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = styles.Regular
				m.Inputs[i].TextStyle = styles.Regular
			}

		}
	}

	cmd := m.updateInputs(msg)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

// updateInputs обновляет поля ввода, обрабатывая сообщения.
func (m InputGroup) updateInputs(msg tea.Msg) tea.Cmd {
	commands := make([]tea.Cmd, len(m.Inputs))

	for i := range m.Inputs {
		m.Inputs[i], commands[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(commands...)
}

// View отображает группу ввода в виде строки.
func (m InputGroup) View() string {
	var (
		b       strings.Builder
		style   lipgloss.Style
		padding int
	)

	maxLabelLength := 0
	for _, input := range m.Inputs {
		if len(input.Placeholder) > maxLabelLength {
			maxLabelLength = len(input.Placeholder)
		}
	}

	for _, input := range m.Inputs {
		label := input.Placeholder
		padding = maxLabelLength - len(label)

		b.WriteString(fmt.Sprintf("%s: %s\n",
			strings.Repeat(" ", padding)+label,
			input.View(),
		))
	}

	b.WriteRune('\n')

	buttonPadding := maxLabelLength + 2
	for i, but := range m.Buttons {
		title := but.Title

		if m.FocusIndex == len(m.Inputs)+i {
			style = styles.Focused
		} else {
			style = styles.Blurred
		}

		b.WriteString(fmt.Sprintf("%s%s\n",
			strings.Repeat(" ", buttonPadding),
			style.Render(title),
		))
	}

	return b.String()
}
