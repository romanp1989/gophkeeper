// Package tui предоставляет функции для работы с текстовым интерфейсом пользователя в приложении GophKeeper.
package tui

import (
	"errors"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"golang.org/x/exp/maps"
	"slices"
)

const (
	// StorageBrowseScreen Экран просмотра хранилища
	StorageBrowseScreen Screen = iota

	// SecretTypeScreen Экран выбора типа секрета
	SecretTypeScreen

	// FilePickScreen Экран выбора файла
	FilePickScreen

	// LoginScreen Экран входа
	LoginScreen

	// RemoteOpenScreen Экран удаленного доступа
	RemoteOpenScreen

	// CredentialEditScreen Экран редактирования учетных данных
	CredentialEditScreen

	// TextEditScreen Экран редактирования текста
	TextEditScreen

	// CardEditScreen Экран редактирования карт
	CardEditScreen

	// BlobEditScreen Экран редактирования файлов
	BlobEditScreen
)

const (
	// BodyPane Основная панель
	BodyPane Position = iota
)

type (
	// Screen представляет собой перечисление доступных экранов в приложении.
	Screen int

	// Position представляет позицию элемента в текстовом интерфейсе пользователя.
	Position int
)

const borderSize = 2

// PaneManager управляет панелями в текстовом интерфейсе пользователя.
type PaneManager struct {
	makers        map[Screen]ScreenMaker
	cache         *Cache
	focused       Position
	panes         map[Position]pane
	width, height int
}

type pane struct {
	model TeaLike
	page  Page
}

// NewPaneManager создает новый менеджер панелей с предоставленными создателями экранов.
func NewPaneManager(makers map[Screen]ScreenMaker) *PaneManager {
	p := &PaneManager{
		makers: makers,
		cache:  NewCache(),
		panes:  make(map[Position]pane),
	}
	return p
}

// Init инициализирует менеджер панелей.
func (pm *PaneManager) Init() tea.Cmd {
	return tea.Batch(
		SetBodyPane(RemoteOpenScreen),
	)
}

// Update обрабатывает сообщения и обновляет состояние панелей.
func (pm *PaneManager) Update(msg tea.Msg) tea.Cmd {
	var (
		commands []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		commands = append(commands, pm.updateModel(pm.focused, msg))

	case tea.WindowSizeMsg:
		pm.width = msg.Width
		pm.height = msg.Height

		pm.updateChildSizes()
	case NavigationMsg:
		commands = append(commands, pm.setPane(msg))
	default:
		commands = pm.cache.UpdateAll(msg)
	}

	return tea.Batch(commands...)
}

// FocusedModel возвращает модель фокусированной панели.
func (pm *PaneManager) FocusedModel() TeaLike {
	return pm.panes[pm.focused].model
}

func (pm *PaneManager) cycleFocusedPane() {
	positions := maps.Keys(pm.panes)
	slices.Sort(positions)

	focusedIndex := int(pm.focused)
	totalPanes := len(pm.panes)

	if focusedIndex >= totalPanes-1 {
		focusedIndex = 0
	} else {
		focusedIndex++
	}

	pm.focusPane(positions[focusedIndex])
}

func (pm *PaneManager) updateChildSizes() {
	for position := range pm.panes {
		pm.updateModel(position, tea.WindowSizeMsg{
			Width:  pm.paneWidth(position) - borderSize,
			Height: pm.paneHeight(position) - borderSize,
		})
	}
}

func (pm *PaneManager) updateModel(position Position, msg tea.Msg) tea.Cmd {

	if p, ok := pm.panes[position]; ok {
		return p.model.Update(msg)
	}

	return nil
}

func (pm *PaneManager) setPane(msg NavigationMsg) tea.Cmd {
	var (
		cmd tea.Cmd
	)

	if p, ok := pm.panes[msg.Position]; ok && p.page == msg.Page {
		if !msg.DisableFocus {
			pm.focusPane(msg.Position)
		}

		return nil
	}

	model := pm.cache.Get(msg.Page)

	maker, ok := pm.makers[msg.Page.Screen]

	if !ok {
		return ReportError(errors.New("no maker for requested screen"))
	}

	var err error
	model, err = maker.Make(msg, 0, 0)
	if err != nil {
		return ReportError(err)
	}

	pm.cache.Put(msg.Page, model)
	cmd = model.Init()

	pm.panes[msg.Position] = pane{model: model, page: msg.Page}
	pm.updateChildSizes()

	if !msg.DisableFocus {
		pm.focusPane(msg.Position)
	}

	return cmd
}

func (pm *PaneManager) focusPane(position Position) {
	if _, ok := pm.panes[position]; ok {
		pm.focused = position
	}
}

func (pm *PaneManager) paneWidth(_ Position) int {
	return pm.width
}

func (pm *PaneManager) paneHeight(_ Position) int {
	return pm.height
}

// View отображает текущий вид менеджера панелей.
func (pm *PaneManager) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top,
			pm.renderPane(BodyPane),
		),
	)
}

func (pm *PaneManager) renderPane(position Position) string {
	if _, ok := pm.panes[position]; !ok {
		return ""
	}

	paneStyle := styles.InactiveBorder.
		Width(pm.paneWidth(position) - borderSize).
		Height(pm.paneHeight(position) - borderSize)

	if position == pm.focused {
		paneStyle = styles.ActiveBorder.Inherit(paneStyle)
	}

	model := pm.panes[position].model
	return paneStyle.Render(model.View())

}

// HelpBindings возвращает привязки клавиш для справки.
func (pm *PaneManager) HelpBindings() (bindings []key.Binding) {
	if model, ok := pm.FocusedModel().(ModelHelpBindings); ok {
		bindings = append(bindings, model.HelpBindings()...)
	}
	return bindings
}
