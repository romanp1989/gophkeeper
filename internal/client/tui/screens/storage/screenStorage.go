// Package storage предоставляет интерфейс и функциональность для просмотра и управления хранением секретов.
package storage

import (
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
	"github.com/romanp1989/gophkeeper/internal/client/storage"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	tableBorderSize = 4
)

type savePathMsg = struct {
	path   string
	secret *domain.Secret
}

// BrowseStorageScreen предоставляет модель экрана для просмотра хранилища секретов.
type BrowseStorageScreen struct {
	storage storage.Storage
	table   table.Model
}

// Make создает экран для просмотра хранилища.
func (s *BrowseStorageScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewStorageBrowseScreenScreen(msg.Storage), nil
}

// NewStorageBrowseScreenScreen создает новый экран для просмотра хранилища.
func NewStorageBrowseScreenScreen(storage storage.Storage) *BrowseStorageScreen {
	scr := &BrowseStorageScreen{
		storage: storage,
		table:   prepareTable(),
	}

	scr.updateRows()

	return scr
}

// Init инициализирует экран и обновляет строки таблицы.
func (s *BrowseStorageScreen) Init() tea.Cmd {
	s.updateRows()
	return nil
}

// Update обновляет состояние экрана в ответ на сообщения.
func (s *BrowseStorageScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	switch msg := msg.(type) {
	case grpc.ReloadSecretList:
		s.updateRows()
	case savePathMsg:
		err := os.WriteFile(msg.path, msg.secret.Blob.FileBytes, 0644)
		if err != nil {
			commands = append(commands, tui.ReportError(err))
		} else {
			commands = append(commands, infoCmd("file saved successfully"))
		}
	case tea.WindowSizeMsg:
		s.table.SetWidth(min(msg.Width, s.colsWidth()))
		s.table.SetHeight(msg.Height - tableBorderSize)
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			cmd = tui.SetBodyPane(tui.SecretTypeScreen, tui.WithStorage(s.storage))
			commands = append(commands, cmd)
		case "e", "enter":
			commands = append(commands, s.handleEdit())
		case "c":
			commands = append(commands, s.handleCopy())
		case "d":
			commands = append(commands, s.handleDelete())

			s.updateRows()
			commands = append(commands, tui.SetBodyPane(tui.StorageBrowseScreen, tui.WithStorage(s.storage)))
		}
	}

	s.table.Focus()
	s.table, cmd = s.table.Update(msg)
	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// View отображает текущий экран.
func (s *BrowseStorageScreen) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Operating storage %s\n", styles.Highlighted.Render(s.storage.String())))
	b.WriteString("Use ↑↓ to navigate, add[a], edit[e], delete[d], copy[c]\n")
	b.WriteString(styles.TableStyle.Render(s.table.View()))

	return styles.StorageScreenStyle.Render(b.String())
}

// HelpBindings возвращает набор горячих клавиш для экрана.
func (s *BrowseStorageScreen) HelpBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add secret")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit secret")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete secret")),
		key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy/save secret")),
	}
}

func (s *BrowseStorageScreen) updateRows() {
	secrets, _ := s.storage.GetAll(context.Background())

	sortSecrets(secrets)

	var rows []table.Row
	for _, sec := range secrets {
		rows = append(rows, table.Row{
			strconv.Itoa(int(sec.ID)),
			sec.Title,
			sec.SecretType,
			sec.CreatedAt.Format("02 Jan 06 15:04"),
			sec.UpdatedAt.Format("02 Jan 06 15:04"),
		})
	}

	s.table.SetRows(rows)
}

func (s *BrowseStorageScreen) handleEdit() tea.Cmd {
	secret, err := s.getSelectedSecret()
	if err != nil {
		return errCmd("failed to load secret: %w", err)
	}

	screen, err := s.getScreenForSecret(secret)
	if err != nil {
		return errCmd("failed to get screen: %w", err)
	}

	return tui.SetBodyPane(screen, tui.WithSecret(secret), tui.WithStorage(s.storage))
}

func (s *BrowseStorageScreen) handleCopy() tea.Cmd {
	secret, err := s.getSelectedSecret()
	if err != nil {
		return errCmd("failed to load secret: %w", err)
	}

	if secret.SecretType == string(domain.BlobSecret) {
		return tui.StringPrompt("choose path to save", func(str string) tea.Cmd { return func() tea.Msg { return savePathMsg{path: str, secret: secret} } })
	}

	if err := clipboard.WriteAll(secret.ToClipboard()); err != nil {
		return errCmd("failed to copy to clipboard: %w", err)
	}

	return infoCmd("secret copied successfully")
}

func (s *BrowseStorageScreen) handleDelete() tea.Cmd {
	secret, err := s.getSelectedSecret()
	if err != nil {
		return errCmd("failed to load secret", err)
	}

	err = s.storage.Delete(context.Background(), secret.ID)
	if err != nil {
		return errCmd("failed to delete secret", err)
	}

	return infoCmd("secret deleted")
}

func errCmd(msg string, err error) tea.Cmd {
	return tui.ReportError(fmt.Errorf("%s: %w", msg, err))
}

func infoCmd(msg string) tea.Cmd {
	return tui.ReportInfo("%s", msg)
}

func (s *BrowseStorageScreen) getSelectedSecret() (secret *domain.Secret, err error) {
	row := s.table.SelectedRow()

	secret, err = s.loadSecret(row[0])
	if err != nil {
		return nil, err
	}

	return secret, err
}

func (s *BrowseStorageScreen) loadSecret(rawID string) (*domain.Secret, error) {
	var err error

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return nil, err
	}

	sec, err := s.storage.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return sec, err
}

func (s *BrowseStorageScreen) getScreenForSecret(secret *domain.Secret) (tui.Screen, error) {
	switch secret.SecretType {
	case string(domain.CredSecret):
		return tui.CredentialEditScreen, nil
	case string(domain.TextSecret):
		return tui.TextEditScreen, nil
	case string(domain.BlobSecret):
		return tui.BlobEditScreen, nil
	case string(domain.CardSecret):
		return tui.CardEditScreen, nil
	default:
		return -1, fmt.Errorf("unknown secret type")
	}
}

func (s *BrowseStorageScreen) colsWidth() int {
	cols := s.table.Columns()
	total := tableBorderSize
	for _, c := range cols {
		total += c.Width
	}

	return total
}

func sortSecrets(secrets []*domain.Secret) {
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].UpdatedAt.After(secrets[j].UpdatedAt)
	})
}

func prepareTable() table.Model {
	columns := []table.Column{
		{Title: "id", Width: 5},
		{Title: "Title", Width: 20},
		{Title: "Secret Type", Width: 20},
		{Title: "Created", Width: 20},
		{Title: "Updated", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	st := table.DefaultStyles()
	st.Header = styles.TableHeaderStyle
	st.Selected = styles.TableSelectedStyle
	t.SetStyles(st)

	return t
}
