// Package texts содержит компоненты для редактирования текстовых секретов в TUI приложении GophKeeper.
package texts

import (
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/storage"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/components"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"time"
)

const (
	textTitle = iota
	textMetadata
	textContent
)

// TextEditScreen структура для экрана редактирования текстовых секретов
type TextEditScreen struct {
	inputGroup components.InputGroup
	secret     *domain.Secret
	storage    storage.Storage
}

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
}

// Make создает экран редактирования текстового секрета
func (s TextEditScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewTextEditScreen(msg.Secret, msg.Storage), nil
}

// NewTextEditScreen создает и инициализирует новый экран для редактирования текстового секрета
func NewTextEditScreen(secret *domain.Secret, store storage.Storage) *TextEditScreen {
	m := TextEditScreen{
		secret:  secret,
		storage: store,
	}

	inputs := make([]textinput.Model, 3)
	inputs[textTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[textMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
	inputs[textContent] = newInput(inputOpts{placeholder: "Content", charLimit: 164})

	var buttons []components.Button
	buttons = append(buttons, components.Button{Title: "[ Submit ]", Cmd: func() tea.Cmd {
		err := m.Submit()
		if err != nil {
			return tui.ReportError(err)
		} else {
			return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
		}
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	if secret.ID > 0 {
		inputs[textTitle].SetValue(secret.Title)
		inputs[textMetadata].SetValue(secret.Metadata)
		inputs[textContent].SetValue(secret.Text.Content)

	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

// Init инициализирует экран при его создании
func (s TextEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

// Update обновляет состояние экрана в ответ на пользовательские действия
func (s *TextEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// Submit обрабатывает отправку формы редактирования текстового секрета
func (s *TextEditScreen) Submit() error {
	var (
		err error
	)

	title := s.inputGroup.Inputs[textTitle].Value()
	metadata := s.inputGroup.Inputs[textMetadata].Value()
	content := s.inputGroup.Inputs[textContent].Value()

	if len(metadata) == 0 {
		return errors.New("please enter metadata")
	}

	if len(title) == 0 {
		return errors.New("please enter title")
	}

	if len(content) == 0 {
		return errors.New("please enter content")
	}

	s.secret.Title = title
	s.secret.Metadata = metadata
	s.secret.Text = &domain.Text{Content: content}
	s.secret.UpdatedAt = time.Now()

	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), s.secret)
	} else {
		err = s.storage.Update(context.Background(), s.secret)
	}

	return err
}

// View отображает представление экрана
func (s TextEditScreen) View() string {
	return screens.RenderContent("Fill in text details:", s.inputGroup.View())
}

func newInput(opts inputOpts) textinput.Model {
	t := textinput.New()
	t.CharLimit = opts.charLimit
	t.Placeholder = opts.placeholder

	if opts.focus {
		t.Focus()
		t.PromptStyle = styles.Focused
		t.TextStyle = styles.Focused
	}

	return t
}
