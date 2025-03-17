// Package tui предоставляет инструменты для создания текстового пользовательского интерфейса (TUI) для приложения GophKeeper.
package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
)

// PromptMsg структура сообщения для создания диалогового окна с запросом ввода.
type PromptMsg struct {
	// Действие, выполняемое после подтверждения ввода
	Action PromptAction

	// Флаг, позволяющий отмену при любом вводе
	AnyCancel bool

	// Клавиша для отмены ввода
	Cancel key.Binding

	// Начальное значение поля
	InitialValue string

	// Клавиша для подтверждения ввода
	Key key.Binding

	// Подсказка для поля ввода
	Placeholder string

	// Текст запроса
	Prompt string
}

// Prompt представляет диалоговое окно с пользовательским вводом.
type Prompt struct {
	action    PromptAction
	anyCancel bool
	cancel    key.Binding
	model     textinput.Model
	trigger   key.Binding
}

// PromptAction определяет тип функции, вызываемой после ввода пользователя.
type PromptAction func(text string) tea.Cmd

// StringPrompt создает команду для отображения диалога с однострочным текстовым полем.
func StringPrompt(prompt string, action PromptAction) tea.Cmd {
	return CmdHandler(PromptMsg{
		Prompt: fmt.Sprintf("%s: ", prompt),
		Action: action,
		Key: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AnyCancel: false,
	})
}

// YesNoPrompt создает команду для отображения диалога с вопросом Да/Нет.
func YesNoPrompt(prompt string, action tea.Cmd) tea.Cmd {
	return CmdHandler(PromptMsg{
		Prompt: fmt.Sprintf("%s (y|n): ", prompt),
		Action: func(_ string) tea.Cmd {
			return action
		},
		Key: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "confirm"),
		),
		AnyCancel: true,
	})
}

// NewPrompt создает новое диалоговое окно с пользовательским вводом.
func NewPrompt(msg PromptMsg) (*Prompt, tea.Cmd) {
	model := textinput.New()
	model.Prompt = msg.Prompt
	model.SetValue(msg.InitialValue)
	model.Placeholder = msg.Placeholder
	model.PlaceholderStyle = styles.Regular.Faint(true)
	blink := model.Focus()

	prompt := Prompt{
		model:     model,
		action:    msg.Action,
		trigger:   msg.Key,
		cancel:    msg.Cancel,
		anyCancel: msg.AnyCancel,
	}
	return &prompt, blink
}

// HandleKey обрабатывает ввод с клавиатуры в диалоговом окне.
func (p *Prompt) HandleKey(msg tea.KeyMsg) (closePrompt bool, cmd tea.Cmd) {
	switch {
	case key.Matches(msg, p.trigger):
		cmd = p.action(p.model.Value())
		closePrompt = true
	case key.Matches(msg, p.cancel), p.anyCancel:
		cmd = ReportInfo("canceled operation")
		closePrompt = true
	default:
		p.model, cmd = p.model.Update(msg)
	}
	return
}

// HandleBlink обрабатывает мигание курсора в текстовом поле.
func (p *Prompt) HandleBlink(msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	default:
		p.model, cmd = p.model.Update(msg)
	}
	return
}

// View настраивает вид диалогового окна.
func (p *Prompt) View(width int) string {
	paddedBorder := styles.ThickBorder.BorderForeground(styles.Red).Padding(0, 1)
	paddedBorderWidth := paddedBorder.GetHorizontalBorderSize() + paddedBorder.GetHorizontalPadding()
	p.model.Width = max(0, width-lipgloss.Width(p.model.Prompt)-paddedBorderWidth)
	content := styles.Regular.Inline(true).MaxWidth(width - paddedBorderWidth).Render(p.model.View())
	return paddedBorder.Width(width - paddedBorder.GetHorizontalBorderSize()).Render(content)
}

// HelpBindings возвращает привязки клавиш для действий в диалоговом окне.
func (p *Prompt) HelpBindings() []key.Binding {
	bindings := []key.Binding{
		p.trigger,
	}
	if p.anyCancel {
		bindings = append(bindings, key.NewBinding(key.WithHelp("n", "cancel")))
	} else {
		bindings = append(bindings, p.cancel)
	}
	return bindings
}
