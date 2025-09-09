package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Save       key.Binding
	Up         key.Binding
	Quit       key.Binding
	Esc        key.Binding
	Help       key.Binding
	Read       key.Binding
	Back       key.Binding
	PageBack   key.Binding
	PageFoward key.Binding
	Enter      key.Binding
}

var Default = KeyMap{
	Save:       key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "Save and Quit")),
	Esc:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "unfocus textarea")),
	Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("^k", "Move Up")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
	Quit:       key.NewBinding(key.WithKeys("q", "esc", "ctrl+q"), key.WithHelp("q", "quit")),
	Read:       key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "Read notes")),
	Back:       key.NewBinding(key.WithKeys("left"), key.WithHelp("left", "Get back")),
	PageBack:   key.NewBinding(key.WithKeys("alt+left"), key.WithHelp("alt+left", "Page Back")),
	PageFoward: key.NewBinding(key.WithKeys("alt+right"), key.WithHelp("alt+right", "Page Foward")),
	Enter:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Enter Note")),
}

var ReadNoteKeys = KeyMap{
	PageBack: key.NewBinding(key.WithKeys("alt+left"), key.WithHelp("alt+left", "Return Home")),
	Quit:     key.NewBinding(key.WithKeys("q", "esc", "ctrl+q"), key.WithHelp("q", "quit")),
	Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Enter Note")),
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
