package model

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gustavo-silva98/adnotes/internal/clientui/keys"
	"github.com/gustavo-silva98/adnotes/internal/repository/file"
)

type SessionState uint

const (
	InsertNoteState SessionState = iota
	ReadNotesState
)

type Model struct {
	State          SessionState
	Textarea       textarea.Model
	Help           help.Model
	Keys           keys.KeyMap
	InputStyle     lipgloss.Style
	Err            error
	Quitting       bool
	MapNotes       map[int]file.Note
	IndexQuery     int
	Context        context.Context
	DB             file.Writer
	TotalItemsNote int
	ListModel      list.Model
	ItemList       []list.Item
}

func New() Model {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time ..."
	ti.Focus()
	ctx := context.Background()
	sql, _ := file.InitDB("banco.db", ctx)

	firstIndex, err := sql.GetFirsIndexPage(ctx)
	if err != nil {
		file.WriteTxt("GET INDEX ERROR: " + err.Error())
	}
	return Model{
		State:          InsertNoteState,
		Textarea:       ti,
		Help:           help.New(),
		Keys:           keys.Default,
		InputStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		IndexQuery:     firstIndex,
		TotalItemsNote: firstIndex,
		Context:        ctx,
		DB:             sql,
	}
}
