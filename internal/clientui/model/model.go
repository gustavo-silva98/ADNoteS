package model

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gustavo-silva98/adnotes/internal/clientui/keys"
)


type Model struct {
	Textarea textarea.Model
	Help help.Model
	Keys keys.KeyMap
	InputStyle lipgloss.Style
	Err error
	Quitting bool
}

func New() Model {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time ..."
	ti.Focus()

	return Model{
		Textarea: ti,
		Help: New().Help,
		Keys: keys.Default,
		InputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}
