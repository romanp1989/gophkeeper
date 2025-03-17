// Package blobs содержит компоненты и логику для редактирования и загрузки файловых данных в TUI приложении.
package blobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/storage"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/components"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"os"
	"time"
)

const (
	blobTitle = iota
	blobMetadata
)

// BlobEditScreen представляет экран для редактирования и загрузки данных файлов (blob).
type BlobEditScreen struct {
	inputGroup components.InputGroup
	secret     *domain.Secret
	storage    storage.Storage
}

type inputOpts struct {
	placeholder string
	charLimit   int
	focus       bool
}

// Make создаёт новый экран BlobEditScreen на основе переданных параметров.
func (s *BlobEditScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewBlobEditScreen(msg.Secret, msg.Storage), nil
}

// NewBlobEditScreen создаёт и инициализирует новый экземпляр BlobEditScreen.
func NewBlobEditScreen(secret *domain.Secret, store storage.Storage) *BlobEditScreen {
	m := BlobEditScreen{
		secret:  secret,
		storage: store,
	}

	inputs := make([]textinput.Model, 2)
	inputs[blobTitle] = newInput(inputOpts{placeholder: "Title", charLimit: 64})
	inputs[blobMetadata] = newInput(inputOpts{placeholder: "Metadata", charLimit: 64})

	var buttons []components.Button
	buttons = append(buttons, components.Button{Title: "[ Pick file ]", Cmd: func() tea.Cmd {

		err := m.validateInputs()
		if err != nil {
			return tui.ReportError(err)
		}

		f := func(args ...any) tea.Cmd {
			str, ok := args[0].(string)
			if !ok {
				return tui.ReportError(fmt.Errorf("error opening file"))
			}

			err := m.Submit(str)
			if err != nil {
				return tui.ReportError(fmt.Errorf("error uploading file: %w", err))
			}

			return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
		}

		return tui.SetBodyPane(tui.FilePickScreen, tui.WithStorage(m.storage), tui.WithCallback(f), tui.WithSecret(secret))
	}})

	buttons = append(buttons, components.Button{Title: "[ Back ]", Cmd: func() tea.Cmd {
		return tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(m.storage))
	}})

	if secret.ID > 0 {
		inputs[blobTitle].SetValue(secret.Title)
		inputs[blobMetadata].SetValue(secret.Metadata)
	}

	m.inputGroup = components.NewInputGroup(inputs, buttons)

	return &m
}

// Init инициализирует компоненты экрана.
func (s *BlobEditScreen) Init() tea.Cmd {
	return s.inputGroup.Init()
}

// Update обрабатывает пользовательский ввод и обновляет состояние экрана.
func (s *BlobEditScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	ig, cmd := s.inputGroup.Update(msg)
	s.inputGroup = ig.(components.InputGroup)

	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// View отображает текущее состояние экрана в виде строки.
func (s *BlobEditScreen) View() string {
	return screens.RenderContent("Fill in file details:", s.inputGroup.View())
}

// Проверяет валидность введенных данных.
func (s *BlobEditScreen) validateInputs() error {
	title := s.inputGroup.Inputs[blobTitle].Value()
	metadata := s.inputGroup.Inputs[blobMetadata].Value()

	if len(title) == 0 {
		return errors.New("please enter title")
	}

	if len(metadata) == 0 {
		return errors.New("please enter metadata")
	}

	return nil
}

// Submit отправляет данные на сервер после проверки валидности.
func (s *BlobEditScreen) Submit(path string) error {
	var (
		err error
	)

	err = s.validateInputs()
	if err != nil {
		return err
	}

	bts, err := readFileToBytes(path)
	if err != nil {
		return err
	}

	s.secret.Title = s.inputGroup.Inputs[blobTitle].Value()
	s.secret.Metadata = s.inputGroup.Inputs[blobMetadata].Value()
	s.secret.Blob = &domain.Blob{
		FileName:  path,
		FileBytes: bts,
	}
	s.secret.UpdatedAt = time.Now()

	if s.secret.ID == 0 {
		s.secret.CreatedAt = time.Now()
		err = s.storage.Create(context.Background(), s.secret)
	} else {
		err = s.storage.Update(context.Background(), s.secret)
	}

	return err
}

// Читает файл по указанному пути и возвращает его содержимое в виде байтового массива.
func readFileToBytes(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
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
