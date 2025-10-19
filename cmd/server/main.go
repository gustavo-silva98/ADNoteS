package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"

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
	fmt.Println("Server iniciando...")
	executeTerminal("InitServer")
	// Captura sinais do sistema (como Ctrl+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Canal para coordenar o encerramento de todas as goroutines
	done := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Goroutine para capturar sinais e fechar o canal done
	go func() {
		<-sigs
		log.Println("Encerrando...")
		close(done) // Fecha o canal, sinalizando para todas as goroutines encerrarem
	}()

	go func() {
		defer wg.Done()
		for {
			// Registra a hotkey
			hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyH)
			err := hk.Register()
			if err != nil {
				log.Println("Erro ao registrar hotkey H:", err)
				continue
			}

			// Usa select para poder cancelar
			select {
			case <-done:
				hk.Unregister()
				return
			case <-hk.Keydown():
				fmt.Println("Foi pressionado o H")
				executeTerminal("InsertNote")
				<-hk.Keyup()
			}
			hk.Unregister()
		}
	}()

	go func() {
		defer wg.Done()
		for {

			// Registra a hotkey
			hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyR)
			err := hk.Register()
			if err != nil {
				log.Println("Erro ao registrar hotkey R:", err)
				continue
			}

			// Usa select para poder cancelar
			select {
			case <-done:
				hk.Unregister()
				return
			case <-hk.Keydown():
				fmt.Println("Foi pressionado o R")
				executeTerminal("ReadNote")
				<-hk.Keyup() // Espera soltar a tecla
			}
			hk.Unregister()
		}
	}()

	go func() {
		defer wg.Done()
		for {
			// Registra a hotkey
			hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyK)
			err := hk.Register()
			if err != nil {
				log.Println("Erro ao registrar hotkey K:", err)
				continue
			}

			// Usa select para poder cancelar
			select {
			case <-done:
				hk.Unregister()
				return
			case <-hk.Keydown():
				fmt.Println("Foi pressionado o K")
				executeTerminal("ExecuteServer")
				<-hk.Keyup()
			}
			hk.Unregister()
		}
	}()

	wg.Wait()
}

var clientCmd *exec.Cmd

func runClient(command string, arg string) {
	log.Println("Processo não está rodando. Iniciando...")

	clientCmd = exec.Command(command, arg)
	//file.WriteTxt(fmt.Sprint(clientCmd.Args))
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

func executeTerminal(arg string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("erro ao achar executavel. Erro: %v", err)
	}
	osName := runtime.GOOS
	clientBinaryName := "client"
	if osName == "windows" {
		clientBinaryName = clientBinaryName + ".exe"
	}
	clientBinaryPath := filepath.Join(filepath.Dir(exePath), clientBinaryName)
	log.Printf("Tentando executar o binário em %v", clientBinaryPath)
	fmt.Println("ClientBinary = ", clientBinaryPath)
	if _, err := os.Stat(clientBinaryPath); err == nil {
		runClient(clientBinaryPath, arg)
	} else {
		log.Printf("Binário do cliente não encontrado em %s: %v", clientBinaryPath, err)
	}
	return nil
}
