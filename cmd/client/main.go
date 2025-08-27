// client.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	// A lógica para abrir o terminal é diferente por sistema operacional.
	// O nosso "client" vai executar um novo terminal e passar a si mesmo como argumento.
	// Este é o método "double-exec".
	if len(os.Args) > 1 && os.Args[1] == "in-terminal" {
		fmt.Println(`"


:::'###::::'########:::::'##::: ##::'#######::'########:'########::'######::
::'## ##::: ##.... ##:::: ###:: ##:'##.... ##:... ##..:: ##.....::'##... ##:
:'##:. ##:: ##:::: ##:::: ####: ##: ##:::: ##:::: ##:::: ##::::::: ##:::..::
'##:::. ##: ##:::: ##:::: ## ## ##: ##:::: ##:::: ##:::: ######:::. ######::
 #########: ##:::: ##:::: ##. ####: ##:::: ##:::: ##:::: ##...:::::..... ##:
 ##.... ##: ##:::: ##:::: ##:. ###: ##:::: ##:::: ##:::: ##:::::::'##::: ##:
 ##:::: ##: ########::::: ##::. ##:. #######::::: ##:::: ########:. ######::
..:::::..::........::::::..::::..:::.......::::::..:::::........:::......:::
		

		
		"`)

		fmt.Println("\n\n\nHello, World!")
		time.Sleep(5 * time.Second)
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
		// No Windows, usa start para abrir um novo cmd.exe
		cmd = exec.Command("cmd.exe", "/C", fmt.Sprintf("start cmd.exe /C \"%s in-terminal\"", exePath))
	default:
		fmt.Println("Sistema operacional não suportado.")
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Erro ao iniciar o novo terminal:", err)
	}
}
