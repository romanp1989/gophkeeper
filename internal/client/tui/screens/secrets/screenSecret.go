// Package secrets содержит компоненты для выбора типа секрета в TUI.
package secrets

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/storage"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/screens"
	"maps"
	"slices"
	"sort"
)

const (
	selectBack = iota
	selectCredential
	selectText
	selectCard
	selectBlob
)

// SecretTypeScreen представляет экран выбора типа секрета.
type SecretTypeScreen struct {
	choice  string
	list    list.Model
	storage storage.Storage
	tea.Model
}

// Make создает экземпляр SecretTypeScreen на основе переданного сообщения.
func (s *SecretTypeScreen) Make(msg tui.NavigationMsg, _, _ int) (tui.TeaLike, error) {
	return NewSecretTypeScreen(msg.Storage), nil
}

// NewSecretTypeScreen создает и инициализирует новый экран выбора типа секрета.
func NewSecretTypeScreen(store storage.Storage) *SecretTypeScreen {
	m := &SecretTypeScreen{storage: store}
	m.prepareSecretListModel()

	return m
}

func (s *SecretTypeScreen) prepareSecretListModel() {
	choices := map[int]string{
		selectBack:       "Go back",
		selectCredential: "Add credentials",
		selectText:       "Add text",
		selectCard:       "Add card info",
		selectBlob:       "Upload file",
	}

	keys := slices.Collect(maps.Keys(choices))
	sort.Ints(keys)

	var items []list.Item
	for i := range keys {
		items = append(items, secretItem{id: i, name: choices[i]})
	}

	l := list.New(items, secretItemDelegate{}, 0, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.KeyMap.Quit.SetEnabled(false)

	s.list = l
}

// Init инициализирует экран выбора типа секрета.
func (s *SecretTypeScreen) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

// Update обновляет состояние и обрабатывает действия пользователя на экране выбора типа секрета.
func (s *SecretTypeScreen) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd      tea.Cmd
		commands []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.list.SetWidth(msg.Width)
		s.list.SetHeight(msg.Height - 2 - 2)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":

			i, ok := s.list.SelectedItem().(secretItem)
			if ok {
				s.choice = i.name
			}

			switch i.id {
			case selectBack:
				cmd = tui.SetBodyPane(
					tui.StorageBrowseScreen,
					tui.WithStorage(s.storage),
				)

				return cmd

			case selectCredential:
				sec := domain.NewSecret(domain.CredSecret)

				cmd = tui.SetBodyPane(
					tui.CredentialEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			case selectText:
				sec := domain.NewSecret(domain.TextSecret)

				cmd = tui.SetBodyPane(
					tui.TextEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			case selectCard:
				sec := domain.NewSecret(domain.CardSecret)

				cmd = tui.SetBodyPane(
					tui.CardEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			case selectBlob:
				sec := domain.NewSecret(domain.BlobSecret)

				cmd = tui.SetBodyPane(
					tui.BlobEditScreen,
					tui.WithStorage(s.storage),
					tui.WithSecret(sec),
				)
			}

			commands = append(commands, cmd)
		}
	}

	s.list, cmd = s.list.Update(msg)
	commands = append(commands, cmd)

	return tea.Batch(commands...)
}

// View отображает экран выбора типа секрета.
func (s *SecretTypeScreen) View() string {
	return screens.RenderContent("Select type of secret:", s.list.View())
}
