package model

import (
	"context"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	EditNoteSate
	ConfirmEditSate
	ResultEditState
	DeleteNoteState
	KillServerState
)

type Model struct {
	State           SessionState
	Textarea        textarea.Model
	Help            help.Model
	Keys            keys.KeyMap
	InputStyle      lipgloss.Style
	Err             error
	Quitting        bool
	MapNotes        map[int]file.Note
	IndexQuery      int
	Context         context.Context
	DB              file.Writer
	TotalItemsNote  int
	ListModel       list.Model
	ItemList        []list.Item
	CurrentPage     int
	ViewpoerContent string
	Ready           bool
	TextareaEdit    textarea.Model
	HelpKeys        []key.Binding
	SelectedNote    list.Item
	ResultMessage   string
	TermHeight      int
	TermWidth       int
}

func NewTextAreaEdit() textarea.Model {
	var (
		cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DF21FF"))

		cursorLineStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#5F03BD")).
				Foreground(lipgloss.Color("#84F5D5"))

		placeholderStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("238"))

		endOfBufferStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("235"))

		focusedPlaceholderStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("99"))

		focusedBorderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("238"))

		blurredBorderStyle = lipgloss.NewStyle().
					Border(lipgloss.HiddenBorder())
	)

	t := textarea.New()
	t.BlurredStyle.Base.MarginBottom(5)
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()

	return t
}

func New() Model {
	ti := textarea.New()
	ti.Placeholder = "Digite sua nota..."
	ti.Focus()
	ctx := context.Background()
	os.Mkdir(filepath.Join("..", "data"), os.ModePerm)
	dbPath := filepath.Join("..", "/data/banco.db")

	sql, _ := file.InitDB(dbPath, ctx)

	textEdit := NewTextAreaEdit()
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
		CurrentPage:    1,
		TextareaEdit:   textEdit,
	}
}
