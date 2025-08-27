package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

// runState mantém o estado do processo do cliente.
var runState = struct {
	sync.Mutex
	running bool
}{}

func main() { mainthread.Init(fn) }
func fn() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyC)
	err := hk.Register()
	if err != nil {
		log.Fatalf("hotkey: Failed to register hotkey: %v", err)
		return
	}
	log.Printf("hotkey: %v is registered\n", hk)
	log.Printf("hotkey: %v is testeeee\n", hk)
	<-hk.Keydown()
	log.Printf("hotkey: %v is down\n", hk)
	executeTerminal()

	<-hk.Keyup()
	log.Printf("hotkey: %v is up\n", hk)
	hk.Unregister()
	log.Printf("hotkey: %v is unregistered\n", hk)

}

var clientCmd *exec.Cmd

func runClient(command string) {
	log.Println("Processo não está rodando. Iniciando...")

	// O binário do cliente deve abrir a sua própria janela de terminal.1
	clientCmd = exec.Command(command)

	if err := clientCmd.Start(); err != nil {
		log.Printf("Erro ao iniciar o comando: %v", err)
		return
	}
	runState.Lock()
	runState.running = true
	runState.Unlock()

	log.Println("Comando executado com sucesso.")

	// Espera pelo processo do cliente terminar para atualizar o estado.
	go func() {
		clientCmd.Wait()
		runState.Lock()
		runState.running = false
		clientCmd = nil
		runState.Unlock()
		log.Println("Cliente encerrado.")
	}()
}

func executeTerminal() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Erro ao achar executavel. Erro: %v", err)
	}
	clientBinaryName := "client"
	clientBinaryPath := filepath.Join(filepath.Dir(exePath), clientBinaryName)
	log.Printf("Tentando executar o binário em %v", clientBinaryPath)

	if _, err := os.Stat(clientBinaryPath); err == nil {
		runClient(clientBinaryPath)
	} else {
		log.Printf("Binário do cliente não encontrado em %s: %v", clientBinaryPath, err)
	}
	return nil
}
