// client.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

func main() {

	// A lógica para abrir o terminal é diferente por sistema operacional.
	// O nosso "client" vai executar um novo terminal e passar a si mesmo como argumento.
	// Este é o método "double-exec".
	if len(os.Args) > 1 && os.Args[1] == "in-terminal" {
		notes := `


:::'###::::'########:::::'##::: ##::'#######::'########:'########::'######::
::'## ##::: ##.... ##:::: ###:: ##:'##.... ##:... ##..:: ##.....::'##... ##:
:'##:. ##:: ##:::: ##:::: ####: ##: ##:::: ##:::: ##:::: ##::::::: ##:::..::
'##:::. ##: ##:::: ##:::: ## ## ##: ##:::: ##:::: ##:::: ######:::. ######::
 #########: ##:::: ##:::: ##. ####: ##:::: ##:::: ##:::: ##...:::::..... ##:
 ##.... ##: ##:::: ##:::: ##:. ###: ##:::: ##:::: ##:::: ##:::::::'##::: ##:
 ##:::: ##: ########::::: ##::. ##:. #######::::: ##:::: ########:. ######::
..:::::..::........::::::..::::..:::.......::::::..:::::........:::......:::
		

		
		`

		/*
			fmt.Println("\n\n\nEscreva sua anotação abaixo!")
			annotation := scanner()
			writeTxt(annotation)
			time.Sleep(1 * time.Second)
			return*/

		termWid, termHeight, err := term.GetSize(os.Stdout.Fd())
		if err != nil {
			fmt.Println(nil)
			return
		}

		var style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Height(termHeight / 2)
		// Centraliza o conteúdo
		output := lipgloss.Place(
			termWid,
			termHeight/2,
			lipgloss.Center,
			lipgloss.Center,
			style.Render(notes),
		)

		fmt.Println(output)
		p := tea.NewProgram(initialModel())

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

/*
	func scanner() string {
		for {
			reader := bufio.NewReader(os.Stdin)
			annotation_str, _ := reader.ReadString('\n')
			annotation_str = strings.TrimSpace(annotation_str)

			return annotation_str
		}
	}
*/
func writeTxt(msg string) {
	// Abre o arquivo para anexar, cria se não existir
	f, err := os.OpenFile("annotations.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(msg + "\n")
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Println("Arquivo escrito com sucesso!")
}

type errMsg error

type model struct {
	textarea textarea.Model
	err      error
}

func initialModel() model {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.Focus()

	return model{
		textarea: ti,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlS:
			//fmt.Println(m.textarea.Value())
			writeTxt(m.textarea.Value())
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	termWid, termHeight, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		fmt.Println(nil)
		return ""
	}
	var textStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 1).
		Width((termWid / 2) * 9 / 10)

	content := fmt.Sprintf(
		"Tellme a story. \n\n%s \n\n%s",
		m.textarea.View(),
		"(Ctrl+s to Save and Quit)",
	) + "\n\n"
	output := lipgloss.Place(
		termWid/2,
		termHeight/2,
		lipgloss.Center,
		lipgloss.Center,
		textStyle.Render(content))

	return output
}
