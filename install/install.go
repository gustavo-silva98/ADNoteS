package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Define nome da pasta de output e nome de binários.
var buildFolderName string = filepath.Join("..", "bin")

var binServerName string = "server"
var binClientName string = "client"

func main() {

	if MakeFolder(buildFolderName) {
		fmt.Printf("Pasta de compilação criada em %v\n", buildFolderName)
		if os := runtime.GOOS; os == "windows" {
			binClientName += ".exe"
			binServerName += ".exe"
		}
		clientBool := BinCompiler(buildFolderName, binClientName, "client")
		serverBool := BinCompiler(buildFolderName, binServerName, "server")

		if clientBool && serverBool {
			fmt.Println("Binários compilados com sucesso.")
		}
	}
}

// Função que cria pasta de output para compilação de binário
func MakeFolder(nameFolder string) bool {
	maked := true
	if err := os.MkdirAll(nameFolder, os.ModePerm); err != nil {
		fmt.Printf("Erro ao criar diretório de build: %v", err)
		maked = false
	}
	return maked
}

// Função que compila binários no destino de compilação definido pela função
func BinCompiler(installFolder string, compiledName string, codeFolder string) bool {
	compiled := true
	clientPath := filepath.Join(installFolder, compiledName)
	err := exec.Command("go", "build", "-o", clientPath, fmt.Sprintf("./cmd/%v", codeFolder)).Run()
	fmt.Printf(".ADNoteS/cmd/%v\n", codeFolder)
	if err != nil {
		fmt.Printf("Erro ao compilar binario: %v\n", err)
		compiled = false
	} else {
		fmt.Printf("Binário %v compilado com sucesso\n", compiledName)
	}

	return compiled
}
