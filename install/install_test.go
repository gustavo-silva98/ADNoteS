package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMakeFolder(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "folder_test")
	defer os.Remove(tmpDir)
	maker := MakeFolder(tmpDir)
	if maker == false {
		t.Error()
	}
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Errorf("Pasta não foi criada")
	}
}

func TestBinCompiler(t *testing.T) {
	cmdDir := "./cmd/client"
	os.MkdirAll(cmdDir, os.ModePerm)
	defer os.RemoveAll("./cmd")

	binName := "test_bin"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	mainContent := `package main
	func main() {}`

	err := os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte(mainContent), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo de teste: %v", err)
	}
	tmpDir := filepath.Join(os.TempDir(), "folder_test")
	defer os.RemoveAll(tmpDir)

	compiler := BinCompiler(tmpDir, binName, "client")
	if !compiler {
		t.Errorf("Teste BinCompiler falhou")
	}

}
