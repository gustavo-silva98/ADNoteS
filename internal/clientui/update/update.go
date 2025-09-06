package update

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gustavo-silva98/adnotes/internal/clientui/model"
	"github.com/gustavo-silva98/adnotes/internal/repository/file"
)

func Update(msg tea.Msg, m model.Model) (model.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Read):
			m.State = model.ReadNotesState
			fmt.Println("STATE É ", m.State)
		}
	}

	switch m.State {
	case model.InsertNoteState:
		return updateInsertNoteState(msg, m)
	case model.ReadNotesState:
		return updateReadNoteState(msg, m)
	}
	return m, nil
}

func updateInsertNoteState(msg tea.Msg, m model.Model) (model.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Save):
			noteExample := file.Note{
				Hour:         0,
				NoteText:     m.Textarea.Value(),
				Reminder:     0,
				PlusReminder: 0,
			}

			ctx := context.Background()
			sql, _ := file.InitDB("banco.db", ctx)

			_, err := sql.InsertNote(&noteExample, ctx)
			if err != nil {
				file.WriteTxt(err.Error())
			}
			time.Sleep(500 * time.Millisecond)
			m.Quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.Keys.Esc):
			if m.Textarea.Focused() {
				m.Textarea.Blur()
			}
		case key.Matches(msg, m.Keys.Quit):
			m.Quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.Keys.Read):
			m.ItemList = queryMapNotes(m)
			m.State = model.ReadNotesState
			return m, nil
		default:
			if !m.Textarea.Focused() {
				cmd = m.Textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}
	m.Textarea, cmd = m.Textarea.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)

}

type noteItem struct {
	title, desc string
}

func (i noteItem) Title() string       { return i.title }
func (i noteItem) Description() string { return i.desc }
func (i noteItem) FilterValue() string { return i.title }

func queryMapNotes(m model.Model) []list.Item {
	mapQuery, err := m.DB.QueryNote(m.IndexQuery-9, m.IndexQuery, m.Context)
	if err != nil {
		file.WriteTxt(err.Error())
	}
	m.MapNotes = mapQuery
	var items []list.Item
	for i := m.IndexQuery; i >= m.IndexQuery-9; i-- {
		if note, ok := m.MapNotes[i]; ok {
			items = append(items, noteItem{
				title: fmt.Sprintf("Note %d", i),
				desc:  note.NoteText,
			})
		}
	}
	return items
}

func updateReadNoteState(msg tea.Msg, m model.Model) (model.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Só recarregue a lista se ItemList estiver vazia (primeira vez) ou se mudar de página
	if len(m.ItemList) == 0 {
		m.ItemList = queryMapNotes(m)
		l := list.New(m.ItemList, list.NewDefaultDelegate(), 50, 30)
		l.Title = "Notas"
		m.ListModel = l
	}

	// Sempre atualize a navegação da lista
	var cmd tea.Cmd
	m.ListModel, cmd = m.ListModel.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Back):
			m.IndexQuery = m.IndexQuery + 10
			if m.IndexQuery > m.TotalItemsNote {
				m.IndexQuery = m.TotalItemsNote
			}
		case key.Matches(msg, m.Keys.Quit):
			m.Quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.Keys.PageBack):
			m.State = model.InsertNoteState
		}
	}

	return m, tea.Batch(cmds...)
}
