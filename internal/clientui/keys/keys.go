package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Save key.Binding
	Up   key.Binding
	Quit key.Binding
	Esc  key.Binding
	Help key.Binding
	Read key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Save, k.Up},
		{k.Help, k.Quit},
	}
}

var Default = KeyMap{
	Save: key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "Save and Quit")),
	Esc:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "unfocus textarea")),
	Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("^k", "Move Up")),
	Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit: key.NewBinding(key.WithKeys("q", "esc", "ctrl+q"), key.WithHelp("q", "quit")),
	Read: key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "Read notes")),
}
