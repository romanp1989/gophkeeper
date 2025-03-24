package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
)

type TeaLike interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
}

type ScreenMaker interface {
	Make(msg NavigationMsg, width, height int) (TeaLike, error)
}

// Page identifies an instance of a model
type Page struct {
	Screen Screen
}

// ModelHelpBindings is implemented by models
// that pass up own help bindings specific to that model.
type ModelHelpBindings interface {
	HelpBindings() []key.Binding
}
