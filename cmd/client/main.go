// client.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	// A lógica para abrir o terminal é diferente por sistema operacional.
	// O nosso "client" vai executar um novo terminal e passar a si mesmo como argumento.
	// Este é o método "double-exec".
	if len(os.Args) > 1 && os.Args[1] == "in-terminal" {
		fmt.Println(`


:::'###::::'########:::::'##::: ##::'#######::'########:'########::'######::
::'## ##::: ##.... ##:::: ###:: ##:'##.... ##:... ##..:: ##.....::'##... ##:
:'##:. ##:: ##:::: ##:::: ####: ##: ##:::: ##:::: ##:::: ##::::::: ##:::..::
'##:::. ##: ##:::: ##:::: ## ## ##: ##:::: ##:::: ##:::: ######:::. ######::
 #########: ##:::: ##:::: ##. ####: ##:::: ##:::: ##:::: ##...:::::..... ##:
 ##.... ##: ##:::: ##:::: ##:. ###: ##:::: ##:::: ##:::: ##:::::::'##::: ##:
 ##:::: ##: ########::::: ##::. ##:. #######::::: ##:::: ########:. ######::
..:::::..::........::::::..::::..:::.......::::::..:::::........:::......:::
		

		
		`)

		fmt.Println("\n\n\nEscreva sua anotação abaixo!")
		annotation := scanner()
		writeTxt(annotation)
		time.Sleep(1 * time.Second)
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

func scanner() string {
	for {
		reader := bufio.NewReader(os.Stdin)
		annotation_str, _ := reader.ReadString('\n')
		annotation_str = strings.TrimSpace(annotation_str)

		return annotation_str
	}
}

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
