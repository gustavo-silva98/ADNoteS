package view_test

import (
	"fmt"
	"testing"

	"github.com/gustavo-silva98/adnotes/internal/clientui/view"
)

func TestKeysForInitState(t *testing.T) {
	totalLenghtTest := 40
	sliceTest := []string{"Ctrl + Shift + H -> Save Note", "Ctrl + Shift + R -> Read Note"}
	sliceTest = view.KeysForInitState(sliceTest, totalLenghtTest)
	//fmt.Println(sliceTest)
	for _, value := range sliceTest {
		if len(value) != totalLenghtTest {
			t.Error("Len diferente do esperado de 10")
		}
	}
}

func TestFormatCenterString(t *testing.T) {
	text := "Texto teste"
	lenght := 20
	result := view.FormatCenterString(text, lenght)
	//fmt.Println(result)
	if len(result) != lenght {
		t.Errorf("Len diferente do esperado de %v", lenght)
	}
}
func TestSliceFormatter(t *testing.T) {
	sliceIn := []string{"Função 1", "Função 2", "Função3"}
	sliceOut := view.SliceFormatter(sliceIn)
	fmt.Println(sliceOut)
	if len(sliceOut) != 2 {
		t.Error("Len do slice de resultado diferente do esperado")
	}
}
