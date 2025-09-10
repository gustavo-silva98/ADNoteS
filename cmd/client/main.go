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
	// O nosso "client" vai executar um novo terminal e passar a si mesmo como argumento.
	// Este é o método "double-exec".
	if len(os.Args) > 1 && os.Args[1] == "in-terminal" {
		p := tea.NewProgram(&app{Model: model.New()})
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
	case "darwin":
		// Abre o terminal e executa o comando.
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell app "Terminal" to do script "%s in-terminal"`, exePath))
	case "windows":
		log.Println("o exePath é ", exePath)
		cmd = exec.Command(
			"cmd.exe", "/C", "start", "cmd.exe", "/C", exePath, "in-terminal",
		)
		log.Println("O comando inserido é : ", cmd)
	default:
		fmt.Println("Sistema operacional não suportado.")
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Erro ao iniciar o novo terminal:", err)
	}
}
