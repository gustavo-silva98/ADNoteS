// client.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gustavo-silva98/adnotes/internal/clientui/model"
	"github.com/gustavo-silva98/adnotes/internal/clientui/update"
	"github.com/gustavo-silva98/adnotes/internal/clientui/view"
)

type app struct {
	model.Model
}

func (a *app) Init() tea.Cmd { return nil }

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := update.Update(msg, &a.Model)
	a.Model = m
	return a, cmd
}

func (a *app) View() string {
	return view.View(a.Model)
}

func main() {

	// A lógica para abrir o terminal é diferente por sistema operacional.
	// O "client" vai executar um novo terminal e passar a si mesmo como argumento.
	// Este é o método "double-exec".
	if len(os.Args) > 1 && os.Args[len(os.Args)-1] == "in-terminal" {
		m := model.New()

		switch os.Args[1] {
		case "InsertNote":
			m.State = model.InsertNoteState
		case "ReadNote":
			m.State = model.ReadNotesState
		case "ExecuteServer":
			m.State = model.ConfirmKillServerState
		case "InitServer":
			m.State = model.InitServerState
		case "AdvancedSearch":
			m.State = model.FullSearchNoteState
		}
		p := tea.NewProgram(&app{Model: m})
		if _, err := p.Run(); err != nil {
			log.Fatal((err))
		}
		return
	}

	// Obtém o caminho do executável atual.
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Erro ao obter o caminho do executável:", err)
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("%s in-terminal", exePath))
	case "windows":
		log.Println("o exePath é ", exePath)
		cmd = exec.Command(
			"cmd.exe", "/C", "start", "cmd.exe", "/C", exePath, os.Args[1], "in-terminal",
		)
	default:
		fmt.Println("Sistema operacional não suportado.")
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Erro ao iniciar o novo terminal:", err)
	}
}
