package view

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/gustavo-silva98/adnotes/internal/clientui/model"
)

var logoLines = []string{
	" _____     _         _____     _       ",
	"|  _  |_ _| |___ ___|   | |___| |_ ___ ",
	"|   __| | | |_ -| -_| | | | _ |  _| -_|",
	"|  |  | | | |   |   | |   |   | | |   |",
	"|__|  |___|_|___|___|_|___|___|_| |___|",
}

var gradientColors = []string{
	"#6d40f3ff", // roxo profundo
	"#8E5AF7",   // lilás mais vivo
	"#A05EFA",   // lavanda
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
	if m.Quitting {
		return "Bye!\n"
	}

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
		BorderForeground(lipgloss.Color("#6d40f3ff")).
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
