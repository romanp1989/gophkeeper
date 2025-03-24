// Package credentials содержит компоненты и логику для создания и редактирования учетных данных в TUI приложении.
package credentials

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
	credTitle = iota
	credMetadata
	credLogin
	credPassword
)

// CredentialEditScreen структура для экрана редактирования учетных данных.
type CredentialEditScreen struct {
	inputGroup components.InputGroup
	secret     *domain.Secret
	storage    storage.Storage
}

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
}

// Make создаёт новый экран CredentialEditScreen на основе переданных параметров.
func (s *CredentialEditScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewCredentialEditScreen(msg.Secret, msg.Storage), nil
}

// NewCredentialEditScreen создаёт и инициализирует новый экземпляр CredentialEditScreen.
func NewCredentialEditScreen(secret *domain.Secret, store storage.Storage) *CredentialEditScreen {
	m := CredentialEditScreen{
		secret:  secret,
		storage: store,
	}

	inputs := make([]textinput.Model, 4)
	inputs[credTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[credMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
	inputs[credLogin] = newInput(inputOpts{placeholder: "Login", charLimit: 64})
	inputs[credPassword] = newInput(inputOpts{placeholder: "Password", charLimit: 64})

	var buttons []components.Button
	buttons = append(buttons, components.Button{Title: "[ Submit ]", Cmd: func() tea.Cmd {
		if err := m.Submit(); err != nil {
			return tui.ReportError(err)
		}
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	if secret.ID > 0 {
		inputs[credTitle].SetValue(secret.Title)
		inputs[credMetadata].SetValue(secret.Metadata)
		inputs[credLogin].SetValue(secret.Credentials.Login)
		inputs[credPassword].SetValue(secret.Credentials.Password)
	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

// Init инициализирует компоненты экрана.
func (s *CredentialEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

// Update обрабатывает пользовательский ввод и обновляет состояние экрана.
func (s *CredentialEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// Submit обрабатывает отправку данных учетной записи в хранилище.
func (s *CredentialEditScreen) Submit() error {
	var (
		err error
	)

	err = s.Validate()
	if err != nil {
		return err
	}

	title := s.inputGroup.Inputs[credTitle].Value()
	metadata := s.inputGroup.Inputs[credMetadata].Value()
	login := s.inputGroup.Inputs[credLogin].Value()
	password := s.inputGroup.Inputs[credPassword].Value()

	s.secret.Title = title
	s.secret.Metadata = metadata
	s.secret.Credentials = &domain.Credentials{Login: login, Password: password}
	s.secret.UpdatedAt = time.Now()

	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), s.secret)
	} else {
		err = s.storage.Update(context.Background(), s.secret)
	}

	return err
}

// Validate Валидация данных, введенных пользователем
func (s *CredentialEditScreen) Validate() error {
	if len(s.inputGroup.Inputs[credMetadata].Value()) == 0 {
		return errors.New("please enter metadata")
	}

	if len(s.inputGroup.Inputs[credTitle].Value()) == 0 {
		return errors.New("please enter title")
	}

	if len(s.inputGroup.Inputs[credLogin].Value()) == 0 {
		return errors.New("please enter login")
	}

	if len(s.inputGroup.Inputs[credPassword].Value()) == 0 {
		return errors.New("please enter password")
	}

	return nil
}

// View отображает текущее состояние экрана в виде строки.
func (s *CredentialEditScreen) View() string {
	return screens.RenderContent("Fill in credential details:", s.inputGroup.View())
}

// newInput создаёт новую модель ввода текста с заданными параметрами.
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
