package file

import (
	"log"
	"os"
)

func WriteTxt(msg string) {
	f, err := os.OpenFile("notes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %v", err)
	}

	defer f.Close()

	_, err = f.WriteString(msg + "\n")
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Println("Arquivo escrito com sucesso")
}
