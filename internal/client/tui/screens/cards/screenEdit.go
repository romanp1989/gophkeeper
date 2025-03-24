// Package cards содержит компоненты и логику для создания и редактирования информации о кредитных картах в TUI приложении.
package cards

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/storage"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/components"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"strconv"
	"time"
)

var errFieldEmpty = func(label string) error { return fmt.Errorf("please enter %s", label) }

const (
	cardTitle = iota
	cardMetadata
	cardNumber
	cardExpYear
	cardExpMonth
	cardCVV
)

// CardEditScreen представляет экран для редактирования и создания информации о картах.
type CardEditScreen struct {
	inputGroup components.InputGroup
	secret     *domain.Secret
	storage    storage.Storage
}

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
}

// Make создаёт новый экран CardEditScreen на основе переданных параметров.
func (s *CardEditScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewCardEditScreen(msg.Secret, msg.Storage), nil
}

// NewCardEditScreen создаёт и инициализирует новый экземпляр CardEditScreen.
func NewCardEditScreen(secret *domain.Secret, store storage.Storage) *CardEditScreen {
	m := CardEditScreen{
		secret:  secret,
		storage: store,
	}

	inputs := make([]textinput.Model, 6)
	inputs[cardTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[cardMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})
	inputs[cardNumber] = newInput(inputOpts{placeholder: "Card number", charLimit: 64})
	inputs[cardExpYear] = newInput(inputOpts{placeholder: "Exp Year", charLimit: 2})
	inputs[cardExpMonth] = newInput(inputOpts{placeholder: "Exp Month", charLimit: 2})
	inputs[cardCVV] = newInput(inputOpts{placeholder: "CVV", charLimit: 6})

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
		inputs[cardTitle].SetValue(secret.Title)
		inputs[cardMetadata].SetValue(secret.Metadata)
		inputs[cardNumber].SetValue(secret.Card.Number)
		inputs[cardExpMonth].SetValue(strconv.FormatUint(uint64(secret.Card.ExpMonth), 10))
		inputs[cardExpYear].SetValue(strconv.FormatUint(uint64(secret.Card.ExpYear), 10))
		inputs[cardCVV].SetValue(strconv.FormatUint(uint64(secret.Card.CVV), 10))
	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

// Init инициализирует компоненты экрана.
func (s *CardEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

// Update обрабатывает пользовательский ввод и обновляет состояние экрана.
func (s *CardEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// Submit обрабатывает отправку данных о карте в хранилище.
func (s *CardEditScreen) Submit() error {
	var (
		err error
	)

	title := s.inputGroup.Inputs[cardTitle].Value()
	metadata := s.inputGroup.Inputs[cardMetadata].Value()
	cardNumberValue := s.inputGroup.Inputs[cardNumber].Value()
	cardExpMonthValue := s.inputGroup.Inputs[cardExpMonth].Value()
	cardExpYearValue := s.inputGroup.Inputs[cardExpYear].Value()
	cardCVVValue := s.inputGroup.Inputs[cardCVV].Value()

	if len(metadata) == 0 {
		return errFieldEmpty("metadata")
	}

	if len(title) == 0 {
		return errFieldEmpty("title")
	}

	if len(cardNumberValue) == 0 {
		return errFieldEmpty("card number")
	}

	if len(cardNumberValue) == 0 {
		return errFieldEmpty("card number")
	}

	if len(cardExpYearValue) == 0 {
		return errFieldEmpty("exp year")
	}

	if len(cardExpMonthValue) == 0 {
		return errFieldEmpty("exp month")
	}

	if len(cardCVVValue) == 0 {
		return errFieldEmpty("CVV")
	}

	s.secret.Title = title
	s.secret.Metadata = metadata
	card := &domain.Card{Number: cardNumberValue}

	card.ExpYear = strToUint32(cardExpYearValue)
	card.ExpMonth = strToUint32(cardExpMonthValue)
	card.CVV = strToUint32(cardCVVValue)

	s.secret.Card = card
	s.secret.UpdatedAt = time.Now()

	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), s.secret)
	} else {
		err = s.storage.Update(context.Background(), s.secret)
	}

	return err
}

// View отображает текущее состояние экрана в виде строки.
func (s *CardEditScreen) View() string {
	return screens.RenderContent("Fill in card details:", s.inputGroup.View())
}

// Создаёт новую модель ввода текста с заданными параметрами.
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

// strToUint32 преобразует строку в uint32.
func strToUint32(str string) uint32 {
	i64, _ := strconv.ParseUint(str, 10, 32)
	return uint32(i64)
}
