package view

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/gustavo-silva98/adnotes/internal/clientui/model"
)

var termWid, termHeight, _ = term.GetSize(os.Stdout.Fd())

var logoLines = []string{
	" _____     _         _____     _       ",
	"|  _  |_ _| |___ ___|   | |___| |_ ___ ",
	"|   __| | | |_ -| -_| | | | _ |  _| -_|",
	"|  |  | | | |   |   | |   |   | | |   |",
	"|__|  |___|_|___|___|_|___|___|_| |___|",
}

var gradientColors = []string{
	"#6d40f3ff", // roxo profundo
	"#7e40faff", // lilás mais vivo
	"#8b3bfcff", // lavanda
	"#BC78FE",   // rosa arroxeado (nova cor)
	"#B262FD",   // violeta claro
}

func renderLogo() string {
	var rendered string
	for i, line := range logoLines {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(gradientColors[i])).
			Align(lipgloss.Center)
		rendered += style.Render(line) + "\n"
	}
	return rendered
}

func View(m model.Model) string {
	var output string
	if m.Quitting {
		return "Bye!\n"
	}
	helpview := m.Help.ShortHelpView(m.HelpKeys)

	var helpStyle = lipgloss.NewStyle().
		AlignVertical(lipgloss.Bottom).
		AlignHorizontal(lipgloss.Center).
		MarginBottom(1)

	switch m.State {
	case model.InsertNoteState:
		output = InsertNoteView(m)
	case model.ReadNotesState:

		horizontal := lipgloss.JoinHorizontal(lipgloss.Top, ListModelView(m), lipgloss.PlaceHorizontal(termWid/2, lipgloss.Center, EditNoteView(m)))
		output = lipgloss.JoinVertical(lipgloss.Center, horizontal, helpStyle.Render(helpview))
	case model.EditNoteSate:
		horizontal := lipgloss.JoinHorizontal(lipgloss.Top, ListModelView(m), lipgloss.PlaceHorizontal(termWid/2, lipgloss.Center, EditNoteView(m)))
		output = lipgloss.JoinVertical(lipgloss.Center, horizontal, helpStyle.Render(helpview))
	}

	return output
}

func InsertNoteView(m model.Model) string {
	termWid, termHeight, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return ""
	}

	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6d40f3ff")).
		Align(lipgloss.Center).
		Width(termWid).
		PaddingBottom(3)

	var textStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7e40fa")).
		Padding(1, 1).
		PaddingBottom(3).
		Width(termWid * 9 / 10)

	var helpStyle = lipgloss.NewStyle().
		AlignVertical(lipgloss.Bottom).
		AlignHorizontal(lipgloss.Center).
		PaddingTop(3).
		MarginBottom(1)

	content := fmt.Sprintf(
		"Digite sua anotação abaixo. \n\n%s",
		m.Textarea.View(),
	)
	helpView := m.Help.View(m.Keys)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		logoStyle.Render(renderLogo()),
		textStyle.Render(content),
		helpStyle.Render(helpView),
	)

	output := lipgloss.Place(
		termWid,
		termHeight,
		lipgloss.Center, lipgloss.Center,
		mainContent,
	)

	return output
}

func EditNoteView(m model.Model) string {
	var textStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7e40faff")).
		MarginLeft(8).
		MarginTop(1).
		Width((termWid / 2)).
		Height(int(float64(termHeight) / 10 * 7.5)).
		AlignVertical(lipgloss.Center)

	//output := lipgloss.PlaceHorizontal((termWid / 2), lipgloss.Right, textStyle.Render(m.TextareaEdit.View()))
	return textStyle.Render(m.TextareaEdit.View())

}

func ListModelView(m model.Model) string {
	var listModelStyle = lipgloss.NewStyle().
		PaddingTop(1).
		Width(int(float64(termWid) / 2.5))

	return listModelStyle.Render(m.ListModel.View())
}
