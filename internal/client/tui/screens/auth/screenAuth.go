// Package auth содержит компоненты и логику для экрана входа и регистрации в TUI приложении.
package auth

import (
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
	"github.com/romanp1989/gophkeeper/internal/client/storage"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/components"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
)

const (
	posLogin = iota
	posPassword
)

// Mode режимы работы экрана: вход или регистрация.
type Mode int

const (
	modeLogin Mode = iota
	modeRegister
)

// AuthenticateScreen структура для экрана входа и регистрации.
type AuthenticateScreen struct {
	client     grpc.ClientGRPCInterface
	inputGroup components.InputGroup
}

type inputOpts struct {
	charLimit   int
	focus       bool
	placeholder string
	secret      bool
}

// Make создаёт новый экран AuthenticateScreen на основе переданного клиента.
func (s *AuthenticateScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewLoginScreen(msg.Client), nil
}

// NewLoginScreen инициализирует и возвращает новый экран входа/регистрации.
func NewLoginScreen(client grpc.ClientGRPCInterface) *AuthenticateScreen {
	m := AuthenticateScreen{
		client: client,
	}

	inputs := make([]textinput.Model, 2)
	inputs[posLogin] = newInput(inputOpts{placeholder: "Login", charLimit: 64})
	inputs[posPassword] = newInput(inputOpts{placeholder: "Password", charLimit: 64, secret: true})

	var buttons []components.Button
	buttons = append(buttons, components.Button{Title: "[ Login ]", Cmd: func() tea.Cmd {
		return m.Submit(modeLogin)
	}})

	buttons = append(buttons, components.Button{Title: "[ Register ]", Cmd: func() tea.Cmd {
		return m.Submit(modeRegister)
	}})

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

// Init инициализирует компоненты экрана.
func (s *AuthenticateScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

// Update обрабатывает пользовательский ввод и обновляет состояние экрана.
func (s *AuthenticateScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// Submit обрабатывает отправку данных для входа или регистрации.
func (s *AuthenticateScreen) Submit(mode Mode) tea.Cmd {
	var (
		token    string
		err      error
		commands []tea.Cmd
	)

	login := s.inputGroup.Inputs[posLogin].Value()
	password := s.inputGroup.Inputs[posPassword].Value()

	if len(login) == 0 {
		return tui.ReportError(errors.New("please enter login"))
	}
	if len(password) == 0 {
		return tui.ReportError(errors.New("please enter password"))
	}

	switch mode {
	case modeLogin:
		token, err = s.client.Login(context.Background(), login, password)
	case modeRegister:
		token, err = s.client.Register(context.Background(), login, password)
	}

	if err != nil {
		commands = append(commands, tui.ReportError(err))
	} else {
		s.client.SetToken(token)
		s.client.SetPassword(password)

		var store storage.Storage
		store, err = storage.NewRemoteStorage(s.client)
		if err != nil {
			commands = append(commands, tui.ReportError(err))
		} else {
			commands = append(commands, tui.ReportInfo("success!"))
			commands = append(commands, tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(store)))
		}
	}

	return tea.Batch(commands...)
}

// View отображает текущее состояние экрана в виде строки.
func (s *AuthenticateScreen) View() string {
	return screens.RenderContent("Fill in credentials:", s.inputGroup.View())
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

	if opts.secret {
		t.EchoMode = textinput.EchoPassword
		t.EchoCharacter = '*'
	}

	return t
}
