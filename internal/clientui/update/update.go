package update

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gustavo-silva98/adnotes/internal/clientui/model"
	"github.com/gustavo-silva98/adnotes/internal/repository/file"
	"github.com/muesli/reflow/wordwrap"
)

// var termWidth, termHeight, _ = term.GetSize(os.Stdout.Fd())
var ctx = context.Background()

// Mensagem para timeout do resultado da edição
type resultEditTimeoutMsg struct{}

const PageSize = 50

func Update(msg tea.Msg, m *model.Model) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.TermHeight = msg.Height - msg.Height/10
		m.TermWidth = msg.Width

	case resultEditTimeoutMsg:
		m.State = model.ReadNotesState
		m.HelpKeys = helpMaker(m)

		return *m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Read):
			m.State = model.ReadNotesState
		}
	}
	m.HelpKeys = helpMaker(m)

	switch m.State {
	case model.InsertNoteState:
		return updateInsertNoteState(msg, m)
	case model.ReadNotesState:
		return updateReadNoteState(msg, m)
	case model.EditNoteSate:
		return updateEditNoteFunc(msg, m)
	case model.ConfirmEditSate:
		return updateConfirmEditNote(msg, m)
	case model.DeleteNoteState:
		return updateConfirmDeleteNote(msg, m)
	case model.ResultEditState:
		m.State = model.ReadNotesState
	}
	return *m, nil
}

func updateInsertNoteState(msg tea.Msg, m *model.Model) (model.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Save):
			noteExample := file.Note{
				Hour:         (time.Now().Unix() - int64(time.Now().Second())),
				NoteText:     m.Textarea.Value(),
				Reminder:     0,
				PlusReminder: 0,
			}

			sql, _ := file.InitDB("banco.db", ctx)

			_, err := sql.InsertNote(&noteExample, ctx)
			if err != nil {
				file.WriteTxt(err.Error())
			}
			time.Sleep(500 * time.Millisecond)
			m.Quitting = true
			return *m, tea.Quit

		case key.Matches(msg, m.Keys.Esc):
			if m.Textarea.Focused() {
				m.Textarea.Blur()
			}
		case key.Matches(msg, m.Keys.Quit):
			m.Quitting = true
			return *m, tea.Quit
		case key.Matches(msg, m.Keys.Read):
			m.ItemList = queryMapNotes(m)
			m.State = model.ReadNotesState
			return *m, nil
		default:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}
	m.Textarea, cmd = m.Textarea.Update(msg)
	cmds = append(cmds, cmd)

	return *m, tea.Batch(cmds...)

}

type noteItem struct {
	title, desc  string
	NoteText     string
	Id           int
	Reminder     int
	PlusReminder int
}

func (i noteItem) Title() string       { return i.title }
func (i noteItem) Description() string { return i.desc }
func (i noteItem) FilterValue() string { return i.title }
func (i noteItem) IdValue() int        { return i.Id }

func queryMapNotes(m *model.Model) []list.Item {
	mapQuery, err := m.DB.QueryNote(PageSize, (m.CurrentPage-1)*PageSize, m.Context)
	if err != nil {
		file.WriteTxt(err.Error())
	}
	m.MapNotes = mapQuery

	var ids []int
	for id := range m.MapNotes {
		ids = append(ids, id)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	items := make([]list.Item, 0, len(mapQuery))
	for _, id := range ids {
		note := m.MapNotes[id]
		noteTimestamp := time.Unix(note.Hour, 0)
		items = append(items, noteItem{
			title:        titleFormatter(note.NoteText),
			desc:         fmt.Sprintf("%v/%d/%v %v:%02d", noteTimestamp.Day(), noteTimestamp.Month(), noteTimestamp.Year(), noteTimestamp.Hour(), noteTimestamp.Minute()),
			NoteText:     note.NoteText,
			Id:           id,
			Reminder:     0,
			PlusReminder: 0,
		})
	}
	return items
}

func updateReadNoteState(msg tea.Msg, m *model.Model) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Só recarregue a lista se ItemList estiver vazia (primeira vez) ou se mudar de página
	if len(m.ItemList) == 0 {
		m.ItemList = queryMapNotes(m)
		d := list.NewDefaultDelegate()
		c := lipgloss.Color("#FE02FF")
		c1 := lipgloss.Color("#7e40fa")
		d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c).Bold(true)
		d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(lipgloss.Color("#9a6bf8ff")).Faint(true)
		d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(c1).BorderLeftForeground(c)
		d.Styles.NormalDesc = d.Styles.NormalDesc.Foreground(lipgloss.Color("#f2c9faff")).Faint(true)

		l := list.New(m.ItemList, d, m.TermWidth/2, (m.TermHeight/10)*7)
		l.Styles.Title = l.Styles.Title.Background(lipgloss.Color("#9D2EB0")).Foreground(lipgloss.Color("#E0D9F6"))
		l.Title = "Notas"
		l.SetShowHelp(false)

		m.ListModel = l
	}

	var cmd tea.Cmd
	m.ListModel, cmd = m.ListModel.Update(msg)
	cmds = append(cmds, cmd)

	selected := m.ListModel.SelectedItem()
	if selected != nil {
		if note, ok := selected.(noteItem); ok {
			wrapped := wordwrap.String(fmt.Sprintf("%v", note.NoteText), m.TextareaEdit.Width())
			// Só atualize o valor se for diferente do atual
			if m.TextareaEdit.Value() != wrapped {
				m.TextareaEdit.SetValue(wrapped)
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			m.Quitting = true
			return *m, tea.Quit
		case key.Matches(msg, m.Keys.PageBack):
			m.State = model.InsertNoteState
		case key.Matches(msg, m.Keys.Delete):
			m.State = model.DeleteNoteState
		case key.Matches(msg, m.Keys.Enter):
			// Ao entrar no modo de edição, inicialize e foque o TextareaEdit
			m.State = model.EditNoteSate
			if !m.TextareaEdit.Focused() {
				cmd = m.TextareaEdit.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}
	return *m, tea.Batch(cmds...)
}

func updateEditNoteFunc(msg tea.Msg, m *model.Model) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Atualize o TextareaEdit com o evento recebido
	var cmd tea.Cmd

	m.TextareaEdit, cmd = m.TextareaEdit.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			m.State = model.ReadNotesState
			if m.TextareaEdit.Focused() {
				m.TextareaEdit.Blur()
			}

		case key.Matches(msg, m.Keys.Save):
			m.State = model.ConfirmEditSate
		}

	}
	return *m, tea.Batch(cmds...)
}

func updateConfirmDeleteNote(msg tea.Msg, m *model.Model) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Yes):
			if selected := m.ListModel.SelectedItem(); selected != nil {
				if note, ok := selected.(noteItem); ok {
					rowsUpdated, err := m.DB.DeleteNoteRepository(ctx, note.Id)
					if err != nil {
						m.ResultMessage = fmt.Sprintf("Erro: %v\nErro ao deletar a nota.", err.Error())
						m.State = model.ReadNotesState
					}
					if rowsUpdated == 1 {
						m.ResultMessage = fmt.Sprintf("Nota %v deletada com sucesso.", note.title)
						m.ItemList = queryMapNotes(m)
						m.ListModel.SetItems(m.ItemList)
						m.State = model.ResultEditState
						return updateResultEditState(msg, m)
					}
				}
			}
		case key.Matches(msg, m.Keys.No):
			m.State = model.ReadNotesState
		}

	}
	return *m, tea.Batch(cmds...)
}

func updateConfirmEditNote(msg tea.Msg, m *model.Model) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Yes):
			if selected := m.ListModel.SelectedItem(); selected != nil {
				if note, ok := selected.(noteItem); ok {
					noteInput := file.Note{
						ID:           note.Id,
						Hour:         time.Now().Unix(),
						NoteText:     m.TextareaEdit.Value(),
						Reminder:     note.Reminder,
						PlusReminder: note.PlusReminder,
					}
					rowsUpdated, err := m.DB.UpdateEditNoteRepository(ctx, noteInput)
					if err != nil {
						m.ResultMessage = fmt.Sprintf("Erro: %v\nErro ao salvar a nota. Necessário averiguar.", err.Error())
						m.State = model.ResultEditState
					}
					if rowsUpdated == 1 {
						m.ResultMessage = fmt.Sprintf("Nota %v editada com sucesso.", note.title)
						m.ItemList = queryMapNotes(m)
						m.ListModel.SetItems(m.ItemList)
						m.State = model.ResultEditState
						return updateResultEditState(msg, m)

					}
				}
			}
		case key.Matches(msg, m.Keys.No):
			m.State = model.EditNoteSate
		}
	}
	return *m, tea.Batch(cmds...)
}

func updateResultEditState(_ tea.Msg, m *model.Model) (model.Model, tea.Cmd) {
	// retorna o cmd que vai enviar resultEditTimeoutMsg após 500ms
	return *m, tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return resultEditTimeoutMsg{}
	})
}

func helpMaker(m *model.Model) []key.Binding {
	// helper pra formatar tecla+descrição
	b := func(keys, helpText string) key.Binding {
		return key.NewBinding(
			key.WithKeys(keys),
			key.WithHelp(keys, helpText),
		)
	}

	switch m.State {
	case model.InsertNoteState:
		return []key.Binding{
			b("Ctrl+s", "Save and Quit"),
			b("Ctrl+r", "Read Notes"),
			b("q", "Quit"),
		}
	case model.ReadNotesState:
		return []key.Binding{
			b("Alt + ←", "Return"),
			b("Enter", "Edit Note"),
			b("Ctrl+d", "Delete Note"),
			b("q", "Quit"),
		}
	case model.EditNoteSate:
		return []key.Binding{
			b("Ctrl+s", "Save Note"),
			b("q", "Quit Editing"),
		}
	}
	return []key.Binding{}
}

func titleFormatter(title string) string {
	maxLineLenght := 40
	splitStr := strings.Split(title, ",")[0]
	if len(splitStr) > maxLineLenght {
		return splitStr[0:maxLineLenght] + "..."
	}
	return splitStr
}
