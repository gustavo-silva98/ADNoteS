package update

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gustavo-silva98/adnotes/internal/clientui/file"
	"github.com/gustavo-silva98/adnotes/internal/clientui/model"
)

func Update(msg tea.Msg, m model.Model) (model.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Save):
			noteExample := file.Note{
				Hour:         0,
				Note:         m.Textarea.Value(),
				Reminder:     0,
				PlusReminder: 0,
			}

			ctx := context.Background()
			sql, _ := file.InitDB("banco.db", ctx)

			id, err := sql.InsertNote(&noteExample, ctx)
			if err != nil {
				file.WriteTxt(err.Error())
			}
			fmt.Println("ID É ", id)
			time.Sleep(5 * time.Second)
			m.Quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.Keys.Esc):
			if m.Textarea.Focused() {
				m.Textarea.Blur()
			}
		case key.Matches(msg, m.Keys.Quit):
			m.Quitting = true
			return m, tea.Quit
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
