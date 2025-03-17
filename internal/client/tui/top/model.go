// Package top содержит модель представления для главного пользовательского интерфейса в TUI-приложении GophKeeper.
package top

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/romanp1989/gophkeeper/internal/client/config"
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
	"github.com/romanp1989/gophkeeper/internal/client/tui"
	"github.com/romanp1989/gophkeeper/internal/client/tui/styles"
	"reflect"
	"strings"
)

type mode int

const (
	normalMode mode = iota
	promptMode
)

const (
	// PromptHeight определяет высоту всплывающего окна ввода, включая границы.
	PromptHeight = 3

	// FooterHeight определяет высоту нижнего колонтитула в пользовательском интерфейсе.
	FooterHeight = 1

	// HelpWidgetHeight определяет высоту виджета помощи, включая границы.
	HelpWidgetHeight = 12

	// MinContentWidth определяет минимальную ширину содержимого на экране.
	MinContentWidth = 80
)

var (
	helpStyle    = styles.Padded.Background(styles.Grey).Foreground(styles.White)
	versionStyle = styles.Padded.Background(styles.DarkGrey).Foreground(styles.White)
)

// Model структура, представляющая модель данных для главного интерфейса пользователя.
type Model struct {
	*tui.PaneManager
	client                    grpc.ClientGRPCInterface
	config                    *config.Config
	err                       error
	height                    int
	helpWidget, versionWidget string
	info                      string
	makers                    map[tui.Screen]tui.ScreenMaker
	mode                      mode
	prompt                    *tui.Prompt
	showHelp                  bool
	width                     int
}

// NewModel создает и инициализирует новую модель интерфейса пользователя.
func NewModel(config *config.Config, client grpc.ClientGRPCInterface) (*Model, error) {
	makers := prepareMakers(client)

	m := Model{
		config:        config,
		client:        client,
		PaneManager:   tui.NewPaneManager(makers),
		makers:        makers,
		helpWidget:    helpStyle.Render(fmt.Sprintf("%s for help", tui.GlobalKeys.Help.Help().Key)),
		versionWidget: versionStyle.Render(fmt.Sprintf("%s (%s - %s)", config.BuildVersion, config.BuildDate, config.BuildCommit)),
	}

	return &m, nil
}

// Init инициализирует состояние модели и запускает начальные команды для настройки интерфейса.
func (m *Model) Init() tea.Cmd {
	return m.PaneManager.Init()
}

// Update обрабатывает сообщения от пользователя и системы, обновляя состояние модели.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tui.PromptMsg:
		m.mode = promptMode
		var blink tea.Cmd
		m.prompt, blink = tui.NewPrompt(msg)

		cmd := m.PaneManager.Update(tea.WindowSizeMsg{
			Height: m.viewHeight(),
			Width:  m.viewWidth(),
		})

		return m, tea.Batch(cmd, blink)

	case tea.KeyMsg:
		m.info = ""
		m.err = nil

		switch m.mode {
		case promptMode:
			closePrompt, cmd := m.prompt.HandleKey(msg)
			if closePrompt {
				m.mode = normalMode
				m.PaneManager.Update(tea.WindowSizeMsg{
					Height: m.viewHeight(),
					Width:  m.viewWidth(),
				})
			}
			return m, cmd
		}

		switch {
		case key.Matches(msg, tui.GlobalKeys.Quit):
			return m, tui.YesNoPrompt("Are you sure you want to exit?", tea.Quit)
		case key.Matches(msg, tui.GlobalKeys.Help):
			m.showHelp = !m.showHelp

			m.PaneManager.Update(tea.WindowSizeMsg{
				Height: m.viewHeight(),
				Width:  m.viewWidth(),
			})

		default:
			return m, m.PaneManager.Update(msg)
		}

	case tui.ErrorMsg:
		m.err = error(msg)

	case tui.InfoMsg:
		m.info = string(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.PaneManager.Update(tea.WindowSizeMsg{
			Height: m.viewHeight(),
			Width:  m.viewWidth(),
		})

	case cursor.BlinkMsg:
		var cmd tea.Cmd
		if m.mode == promptMode {
			cmd = m.prompt.HandleBlink(msg)
		} else {
			cmd = m.PaneManager.FocusedModel().Update(msg)
		}
		return m, cmd
	default:
		commands = append(commands, m.PaneManager.Update(msg))
	}

	return m, tea.Batch(commands...)
}

// View генерирует текстовое представление интерфейса пользователя для отображения в терминале.
func (m *Model) View() string {
	var components []string

	if m.mode == promptMode {
		components = append(components, m.prompt.View(m.width))
	}

	components = append(components, styles.Regular.
		Height(m.viewHeight()).
		Width(m.viewWidth()).
		Render(m.PaneManager.View()),
	)

	if m.showHelp {
		components = append(components, m.help())
	}

	footer := m.helpWidget
	if m.err != nil {
		footer += styles.Regular.Padding(0, 1).
			Background(styles.Red).
			Foreground(styles.White).
			Width(m.availableFooterMsgWidth()).
			Render(m.err.Error())
	} else if m.info != "" {
		footer += styles.Padded.
			Foreground(styles.Black).
			Background(styles.LightGreen).
			Width(m.availableFooterMsgWidth()).
			Render(m.info)
	} else {
		footer += styles.Padded.
			Foreground(styles.Black).
			Background(styles.EvenLighterGrey).
			Width(m.availableFooterMsgWidth()).
			Render(m.info)
	}
	footer += m.versionWidget

	components = append(components, styles.Regular.
		Inline(true).
		MaxWidth(m.width).
		Width(m.width).
		Render(footer),
	)
	return strings.Join(components, "\n")
}

func (m *Model) availableFooterMsgWidth() int {
	return max(0, m.width-lipgloss.Width(m.helpWidget)-lipgloss.Width(m.versionWidget))
}

func (m *Model) viewHeight() int {
	vh := m.height - FooterHeight
	if m.mode == promptMode {
		vh -= PromptHeight
	}
	if m.showHelp {
		vh -= HelpWidgetHeight
	}

	return vh
}

func (m *Model) viewWidth() int {
	return max(MinContentWidth, m.width)
}

func (m *Model) help() string {
	bindings := []key.Binding{tui.GlobalKeys.Help, tui.GlobalKeys.Quit}

	switch m.mode {
	case promptMode:
		bindings = append(bindings, m.prompt.HelpBindings()...)
	default:
		bindings = append(bindings, m.HelpBindings()...)
	}

	bindings = append(bindings, keyMapToSlice(tui.GlobalKeys)...)
	bindings = removeDuplicateBindings(bindings)

	var (
		pairs []string
		width int
		rows  = HelpWidgetHeight - 2
	)
	for i := 0; i < len(bindings); i += rows {
		var (
			keys         []string
			descriptions []string
		)
		for j := i; j < min(i+rows, len(bindings)); j++ {
			keys = append(keys, styles.HelpKeyStyle.Render(bindings[j].Help().Key))
			descriptions = append(descriptions, styles.HelpDescStyle.Render(bindings[j].Help().Desc))
		}

		var cols []string
		if len(pairs) > 0 {
			cols = []string{"   "}
		}
		cols = append(cols,
			strings.Join(keys, "\n"),
			strings.Join(descriptions, "\n"),
		)

		pair := lipgloss.JoinHorizontal(lipgloss.Top, cols...)
		width += lipgloss.Width(pair)
		if width > m.width-2 {
			break
		}
		pairs = append(pairs, pair)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, pairs...)
	return styles.Border.Height(rows).Width(m.width - 2).Render(content)
}

func removeDuplicateBindings(bindings []key.Binding) []key.Binding {
	seen := make(map[string]struct{})
	var i int
	for _, b := range bindings {
		bindKey := strings.Join(b.Keys(), " ")
		if _, ok := seen[bindKey]; ok {
			continue
		}
		seen[bindKey] = struct{}{}
		bindings[i] = b
		i++
	}
	return bindings[:i]
}

func keyMapToSlice(t any) (bindings []key.Binding) {
	typ := reflect.TypeOf(t)
	if typ.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < typ.NumField(); i++ {
		v := reflect.ValueOf(t).Field(i)
		bindings = append(bindings, v.Interface().(key.Binding))
	}
	return
}
